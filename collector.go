package main

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func powerConsumptionWarningValue(info DeviceInfo) float64 {
	if info.Electrical.PowerConsumptionWarning {
		return 1
	}
	return 0
}

type collector struct {
	temperature             *prometheus.Desc
	voltage                 *prometheus.Desc
	current                 *prometheus.Desc
	powerConsumption        *prometheus.Desc
	powerConsumptionMax     *prometheus.Desc
	powerConsumptionWarning *prometheus.Desc
}

func NewCollector() prometheus.Collector {
	return &collector{
		temperature: prometheus.NewDesc("xrt_device_temperature",
			"Temperature of the device in degrees Celsius",
			[]string{"device_id", "serial", "location_id", "description"},
			nil,
		),
		voltage: prometheus.NewDesc("xrt_device_voltage",
			"Voltage of the device in Volts",
			[]string{"device_id", "serial", "location_id", "description"},
			nil,
		),
		current: prometheus.NewDesc("xrt_device_current",
			"Current of the device in Amperes",
			[]string{"device_id", "serial", "location_id", "description"},
			nil,
		),
		powerConsumption: prometheus.NewDesc("xrt_device_power_consumption",
			"Power consumption of the device in Watts",
			[]string{"device_id", "serial"},
			nil,
		),
		powerConsumptionMax: prometheus.NewDesc("xrt_device_power_consumption_max",
			"Maximum power consumption of the device in watts",
			[]string{"device_id", "serial"},
			nil,
		),
		powerConsumptionWarning: prometheus.NewDesc("xrt_device_power_consumption_warning",
			"Whether the power consumption of the device has crossed a threshold",
			[]string{"device_id", "serial"},
			nil,
		),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperature
	ch <- c.voltage
	ch <- c.current
	ch <- c.powerConsumption
	ch <- c.powerConsumptionMax
	ch <- c.powerConsumptionWarning
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	devices, err := GetDevices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve XRT devices: %s\n", err)
		return
	}

	for _, device := range devices {
		info, err := GetDeviceInfo(device)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to retrieve XRT device info for device %s: %s\n", device, err)
			continue
		}

		serial := info.Platforms[0].Controller.CardMgmtController.SerialNumber

		for _, t := range info.Thermals {
			if t.IsPresent {
				ch <- prometheus.MustNewConstMetric(c.temperature, prometheus.GaugeValue, t.TempC, info.DeviceID, serial, t.LocationID, t.Description)
			}
		}

		for _, p := range info.Electrical.PowerRails {
			if p.Voltage.IsPresent {
				ch <- prometheus.MustNewConstMetric(c.voltage, prometheus.GaugeValue, p.Voltage.Volts, info.DeviceID, serial, p.Id, p.Description)
			}

			if p.Current.IsPresent {
				ch <- prometheus.MustNewConstMetric(c.current, prometheus.GaugeValue, p.Current.Amps, info.DeviceID, serial, p.Id, p.Description)
			}
		}

		ch <- prometheus.MustNewConstMetric(c.powerConsumption, prometheus.GaugeValue, info.Electrical.PowerConsumptionWatts, info.DeviceID, serial)
		ch <- prometheus.MustNewConstMetric(c.powerConsumptionMax, prometheus.GaugeValue, info.Electrical.PowerConsumptionMaxWatts, info.DeviceID, serial)
		ch <- prometheus.MustNewConstMetric(c.powerConsumptionWarning, prometheus.GaugeValue, powerConsumptionWarningValue(info), info.DeviceID, serial)
	}
}
