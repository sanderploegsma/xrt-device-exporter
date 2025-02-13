package main

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

func bool2gauge(value bool) float64 {
	if value {
		return 1
	}

	return 0
}

type collector struct {
	logger                        *slog.Logger
	xrt                           Xrt
	xrtStatus                     *prometheus.Desc
	deviceReady                   *prometheus.Desc
	deviceTemperature             *prometheus.Desc
	deviceVoltage                 *prometheus.Desc
	deviceCurrent                 *prometheus.Desc
	devicePowerConsumption        *prometheus.Desc
	devicePowerConsumptionMax     *prometheus.Desc
	devicePowerConsumptionWarning *prometheus.Desc
}

func NewCollector(logger *slog.Logger) prometheus.Collector {
	return &collector{
		logger: logger,
		xrt:    NewXrt(logger),
		xrtStatus: prometheus.NewDesc("xrt_status",
			"Whether the XRT is available",
			[]string{"version", "branch"},
			nil,
		),
		deviceReady: prometheus.NewDesc("xrt_device_ready",
			"Whether the device is ready",
			[]string{"device_id", "shell"},
			nil,
		),
		deviceTemperature: prometheus.NewDesc("xrt_device_temperature",
			"Temperature of the device in degrees Celsius",
			[]string{"device_id", "serial", "shell", "location_id", "description"},
			nil,
		),
		deviceVoltage: prometheus.NewDesc("xrt_device_voltage",
			"Voltage of the device in Volts",
			[]string{"device_id", "serial", "shell", "location_id", "description"},
			nil,
		),
		deviceCurrent: prometheus.NewDesc("xrt_device_current",
			"Current of the device in Amperes",
			[]string{"device_id", "serial", "shell", "location_id", "description"},
			nil,
		),
		devicePowerConsumption: prometheus.NewDesc("xrt_device_power_consumption",
			"Power consumption of the device in Watts",
			[]string{"device_id", "serial", "shell"},
			nil,
		),
		devicePowerConsumptionMax: prometheus.NewDesc("xrt_device_power_consumption_max",
			"Maximum power consumption of the device in watts",
			[]string{"device_id", "serial", "shell"},
			nil,
		),
		devicePowerConsumptionWarning: prometheus.NewDesc("xrt_device_power_consumption_warning",
			"Whether the power consumption of the device has crossed a threshold",
			[]string{"device_id", "serial", "shell"},
			nil,
		),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.xrtStatus
	ch <- c.deviceReady
	ch <- c.deviceTemperature
	ch <- c.deviceVoltage
	ch <- c.deviceCurrent
	ch <- c.devicePowerConsumption
	ch <- c.devicePowerConsumptionMax
	ch <- c.devicePowerConsumptionWarning
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Info("Retrieving host info")
	hostInfo, err := c.xrt.GetHostInfo()
	if err != nil {
		c.logger.Error("Failed to retrieve XRT host info", slog.Any("error", err))
		ch <- prometheus.MustNewConstMetric(c.xrtStatus, prometheus.GaugeValue, 0, "unknown", "unknown")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.xrtStatus, prometheus.GaugeValue, 1, hostInfo.Xrt.Version, hostInfo.Xrt.Branch)

	for _, device := range hostInfo.Devices {
		deviceLogger := c.logger.With(slog.String("device", device.BDF))

		ch <- prometheus.MustNewConstMetric(c.deviceReady, prometheus.GaugeValue, bool2gauge(device.IsReady), device.BDF, device.VBNV)

		if !device.IsReady {
			deviceLogger.Info("Device not ready, not retrieving info")
			continue
		}

		deviceLogger.Info("Retrieving device info")
		deviceInfo, err := c.xrt.GetDeviceInfo(device.BDF)
		if err != nil {
			deviceLogger.Error("Failed to retrieve XRT device info", slog.Any("error", err))
			continue
		}

		serial := deviceInfo.Platforms[0].Controller.CardMgmtController.SerialNumber
		shell := deviceInfo.Platforms[0].StaticRegion.VBNV

		for _, t := range deviceInfo.Thermals {
			if t.IsPresent {
				ch <- prometheus.MustNewConstMetric(c.deviceTemperature, prometheus.GaugeValue, t.TempC, deviceInfo.DeviceID, serial, shell, t.LocationID, t.Description)
			}
		}

		for _, p := range deviceInfo.Electrical.PowerRails {
			if p.Voltage.IsPresent {
				ch <- prometheus.MustNewConstMetric(c.deviceVoltage, prometheus.GaugeValue, p.Voltage.Volts, deviceInfo.DeviceID, serial, shell, p.Id, p.Description)
			}

			if p.Current.IsPresent {
				ch <- prometheus.MustNewConstMetric(c.deviceCurrent, prometheus.GaugeValue, p.Current.Amps, deviceInfo.DeviceID, serial, shell, p.Id, p.Description)
			}
		}

		ch <- prometheus.MustNewConstMetric(c.devicePowerConsumption, prometheus.GaugeValue, deviceInfo.Electrical.PowerConsumptionWatts, deviceInfo.DeviceID, serial, shell)
		ch <- prometheus.MustNewConstMetric(c.devicePowerConsumptionMax, prometheus.GaugeValue, deviceInfo.Electrical.PowerConsumptionMaxWatts, deviceInfo.DeviceID, serial, shell)
		ch <- prometheus.MustNewConstMetric(c.devicePowerConsumptionWarning, prometheus.GaugeValue, bool2gauge(deviceInfo.Electrical.PowerConsumptionWarning), deviceInfo.DeviceID, serial, shell)
	}
}
