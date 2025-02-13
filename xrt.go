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

type HostInfo struct {
	Xrt struct {
		Version string `json:"version"`
		Branch  string `json:"branch"`
	} `json:"xrt"`
	Devices []struct {
		BDF     string `json:"bdf"`
		VBNV    string `json:"vbnv"`
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

type Xrt interface {
	GetHostInfo() (HostInfo, error)
	GetDeviceInfo(id string) (DeviceInfo, error)
}

type xrtWrapper struct {
	logger  *slog.Logger
	xrtPath string
}

func NewXrt(logger *slog.Logger) Xrt {
	// TODO: verify whether xrtPath is valid

	return &cache{
		logger:          logger,
		ttl:             *cacheTtl,
		deviceInfoCache: make(map[string]cachedDeviceInfo),
		xrt: &xrtWrapper{
			logger:  logger,
			xrtPath: *xrtPath,
		},
	}
}

func (x *xrtWrapper) GetHostInfo() (HostInfo, error) {
	x.logger.Debug("Creating temporary file")
	f, err := os.CreateTemp("", "xbutil-output-*")
	if err != nil {
		return HostInfo{}, err
	}
	defer func() {
		x.logger.Debug("Removing temporary file", slog.String("file", f.Name()))
		os.Remove(f.Name())
	}()

	cmd := exec.Command(filepath.Join(x.xrtPath, "bin", "xbutil"), "examine", "--format", "json", "--output", f.Name(), "--force")
	x.logger.Debug(fmt.Sprintf("Running command: %s", cmd.Args))
	if err = cmd.Run(); err != nil {
		return HostInfo{}, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return HostInfo{}, err
	}

	var info struct {
		System struct {
			Host HostInfo `json:"host"`
		} `json:"system"`
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return HostInfo{}, err
	}

	return info.System.Host, nil
}

func (x *xrtWrapper) GetDeviceInfo(id string) (DeviceInfo, error) {
	logger := x.logger.With(slog.String("device", id))

	logger.Debug("Creating temporary file")
	f, err := os.CreateTemp("", "xbutil-output-*")
	if err != nil {
		return DeviceInfo{}, err
	}
	defer func() {
		logger.Debug("Removing temporary file", slog.String("file", f.Name()))
		os.Remove(f.Name())
	}()

	cmd := exec.Command(filepath.Join(x.xrtPath, "bin", "xbutil"), "examine", "--device", id, "--report", "thermal", "--report", "electrical", "--report", "platform", "--format", "json", "--output", f.Name(), "--force")
	logger.Debug(fmt.Sprintf("Running command: %s", cmd.Args))
	if err = cmd.Run(); err != nil {
		return DeviceInfo{}, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return DeviceInfo{}, err
	}

	var info struct {
		Devices []DeviceInfo `json:"devices"`
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return DeviceInfo{}, err
	}

	for _, device := range info.Devices {
		if device.DeviceID == id {
			return device, nil
		}
	}

	return DeviceInfo{}, fmt.Errorf("no device with id %s", id)
}

type cachedDeviceInfo struct {
	info   DeviceInfo
	expiry time.Time
}

type cachedHostInfo struct {
	info   HostInfo
	expiry time.Time
}

type cache struct {
	logger *slog.Logger
	xrt    *xrtWrapper
	ttl    time.Duration

	deviceInfoCache map[string]cachedDeviceInfo
	hostInfoCache   cachedHostInfo
}

func (c *cache) GetHostInfo() (HostInfo, error) {
	if time.Now().Before(c.hostInfoCache.expiry) {
		c.logger.Debug("Using cached host info")
		return c.hostInfoCache.info, nil
	}

	c.logger.Debug("Cached host info expired or not set, refreshing")
	info, err := c.xrt.GetHostInfo()
	if err != nil {
		return info, err
	}

	expiry := time.Now().Add(c.ttl)
	c.hostInfoCache = cachedHostInfo{
		info:   info,
		expiry: expiry,
	}

	c.logger.Debug(fmt.Sprintf("Cached host info expires at: %s", expiry))
	return info, nil
}

func (c *cache) GetDeviceInfo(id string) (DeviceInfo, error) {
	logger := c.logger.With(slog.String("device", id))

	if value, ok := c.deviceInfoCache[id]; ok && time.Now().Before(value.expiry) {
		logger.Debug("Using cached device info")
		return value.info, nil
	}

	logger.Debug("Cached device info expired or not set, refreshing")
	info, err := c.xrt.GetDeviceInfo(id)
	if err != nil {
		return info, err
	}

	expiry := time.Now().Add(c.ttl)
	c.deviceInfoCache[id] = cachedDeviceInfo{
		info:   info,
		expiry: expiry,
	}

	logger.Debug(fmt.Sprintf("Cached device info expires at %s", expiry))
	return info, nil
}
