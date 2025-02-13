package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
)

var (
	xrtPath = kingpin.Flag(
		"xrt.path",
		"Path to the XRT installation directory",
	).Default("/opt/xilinx/xrt").Envar("XILINX_XRT").ExistingDir()

	cacheTtl = kingpin.Flag(
		"xrt.cache-ttl",
		"Time to cache XRT device information",
	).Default("5s").Duration()
)

type HostInformation struct {
	Xrt struct {
		Version string `json:"version"`
		Branch  string `json:"branch"`
	} `json:"xrt"`
	Devices []struct {
		BDF     string `json:"bdf"`
		IsReady bool   `json:"is_ready,string"`
	} `json:"devices"`
}

type ThermalInfo struct {
	LocationID  string  `json:"location_id"`
	Description string  `json:"description"`
	TempC       float64 `json:"temp_C,string"`
	IsPresent   bool    `json:"is_present,string"`
}

type PowerRail struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Voltage     struct {
		Volts     float64 `json:"volts,string"`
		IsPresent bool    `json:"is_present,string"`
	} `json:"voltage"`
	Current struct {
		Amps      float64 `json:"amps,string"`
		IsPresent bool    `json:"is_present,string"`
	} `json:"current"`
}

type ElectricalInfo struct {
	PowerRails               []PowerRail `json:"power_rails"`
	PowerConsumptionMaxWatts float64     `json:"power_consumption_max_watts,string"`
	PowerConsumptionWatts    float64     `json:"power_consumption_watts,string"`
	PowerConsumptionWarning  bool        `json:"power_consumption_warning,string"`
}

type PlatformInfo struct {
	StaticRegion struct {
		VBNV        string `json:"vbnv"`
		LogicalUUID string `json:"logic_uuid"`
		JTagIdCode  string `json:"jtag_idcode"`
		FPGAName    string `json:"fpga_name"`
	} `json:"static_region"`
	Controller struct {
		CardMgmtController struct {
			SerialNumber string `json:"serial_number"`
			OEMID        string `json:"oem_id"`
		} `json:"card_mgmt_controller"`
	} `json:"controller"`
}

type DeviceInfo struct {
	InterfaceType string         `json:"interface_type"`
	DeviceID      string         `json:"device_id"`
	Thermals      []ThermalInfo  `json:"thermals"`
	Electrical    ElectricalInfo `json:"electrical"`
	Platforms     []PlatformInfo `json:"platforms"`
}

type DeviceInfoRepository interface {
	GetDeviceInfo() []DeviceInfo
}

type xrtWrapper struct {
	logger  *slog.Logger
	xrtPath string
}

func NewDeviceInfoRepository(logger *slog.Logger) DeviceInfoRepository {
	// TODO: verify whether xrtPath is valid

	return &cache{
		logger: logger,
		ttl:    *cacheTtl,
		xrt: &xrtWrapper{
			logger:  logger,
			xrtPath: *xrtPath,
		},
	}
}

func (x *xrtWrapper) GetDeviceInfo() []DeviceInfo {
	devices := make([]DeviceInfo, 0)

	deviceIds, err := x.getDeviceIds()
	if err != nil {
		x.logger.Error("Failed to retrieve XRT device ids", slog.Any("error", err))
		return devices
	}

	for _, id := range deviceIds {
		info, err := x.getSingleDeviceInfo(id)
		if err != nil {
			x.logger.Error("Failed to retrieve XRT device info", slog.String("device", id), slog.Any("error", err))
			continue
		}

		devices = append(devices, info)
	}

	return devices
}

func (x *xrtWrapper) getDeviceIds() ([]string, error) {
	x.logger.Info("Retrieving XRT device ids")

	f, err := os.CreateTemp("", "xbutil-output-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	cmd := exec.Command(filepath.Join(x.xrtPath, "bin", "xbutil"), "examine", "--format", "json", "--output", f.Name(), "--force")
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var info struct {
		System struct {
			Host HostInformation `json:"host"`
		} `json:"system"`
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	devices := make([]string, 0)
	for _, device := range info.System.Host.Devices {
		if device.IsReady {
			devices = append(devices, device.BDF)
		}
	}

	return devices, nil
}

func (x *xrtWrapper) getSingleDeviceInfo(id string) (DeviceInfo, error) {
	x.logger.Info("Retrieving XRT device info", slog.String("device", id))
	var empty DeviceInfo

	f, err := os.CreateTemp("", "xbutil-output-*")
	if err != nil {
		return empty, err
	}
	defer os.Remove(f.Name())

	cmd := exec.Command(filepath.Join(x.xrtPath, "bin", "xbutil"), "examine", "--device", id, "--report", "thermal", "--report", "electrical", "--report", "platform", "--format", "json", "--output", f.Name(), "--force")
	if err = cmd.Run(); err != nil {
		return empty, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return empty, err
	}

	var info struct {
		Devices []DeviceInfo `json:"devices"`
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return empty, err
	}

	for _, device := range info.Devices {
		if device.DeviceID == id {
			return device, nil
		}
	}

	return empty, fmt.Errorf("no device with id %s", id)
}

type cache struct {
	logger *slog.Logger
	xrt    *xrtWrapper
	ttl    time.Duration

	value  []DeviceInfo
	expiry time.Time
}

func (c *cache) GetDeviceInfo() []DeviceInfo {
	if time.Now().Before(c.expiry) {
		c.logger.Info("Using cached device info")
		return c.value
	}

	c.logger.Debug("Cached device info expired or not set, refreshing")

	c.value = c.xrt.GetDeviceInfo()
	c.expiry = time.Now().Add(c.ttl)
	c.logger.Debug(fmt.Sprintf("Cached device info expires at %s", c.expiry))
	return c.value
}
