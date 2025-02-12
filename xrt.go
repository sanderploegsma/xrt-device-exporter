package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
)

var (
	xrtPath = kingpin.Flag(
		"xrt.path",
		"Path to the XRT installation directory",
	).Default("/opt/xilinx/xrt").Envar("XILINX_XRT").ExistingDir()
)

func getDeviceJSON(opts ...string) ([]byte, error) {
	f, err := os.CreateTemp("", "xbutil-output-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	executable := filepath.Join(*xrtPath, "bin", "xbutil")
	args := []string{"examine"}
	args = append(args, opts...)
	args = append(args, "--format", "json", "--output", f.Name(), "--force")

	cmd := exec.Command(executable, args...)
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	return io.ReadAll(f)
}

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

func GetDevices() ([]string, error) {
	f, err := os.CreateTemp("", "xbutil-output-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	cmd := exec.Command(*xrtPath+"/bin/xbutil", "examine", "--format", "json", "--output", f.Name(), "--force")
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	data, err := getDeviceJSON()
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

func GetDeviceInfo(id string) (DeviceInfo, error) {
	var empty DeviceInfo

	data, err := getDeviceJSON("--device", id, "--report", "thermal", "--report", "electrical", "--report", "platform")
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
