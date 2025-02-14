# xrt-device-exporter

Export XRT device information as Prometheus metrics.

> [!NOTE]
> This exporter has been tested with Xilinx XRT 2022.2 combined with the Alveo U55C.
> Other XRT versions and FPGAs _should_ be supported, but I have no way of testing this.

## Prerequisites

Make sure the Xilinx XRT is installed.

## Usage

    $ xrt-device-exporter --help
    usage: xrt-device-exporter [<flags>]


    Flags:
    -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
        --xrt.path=/opt/xilinx/xrt  
                                    Path to the XRT installation directory ($XILINX_XRT)
        --xrt.cache-ttl=5s         Time to cache XRT device information
        --web.telemetry-path="/metrics"  
                                    Path under which to expose metrics.
        --[no-]web.systemd-socket  Use systemd socket activation listeners instead of port listeners (Linux only).
        --web.listen-address=:9101 ...  
                                    Addresses on which to expose metrics and web interface. Repeatable for multiple addresses. Examples: `:9100` or `[::1]:9100` for http, `vsock://:9100` for vsock
        --web.config.file=""       Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
        --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
        --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
    -v, --[no-]version             Show application version.

## Example output

    # HELP xrt_device_current Current of the device in Amperes
    # TYPE xrt_device_current gauge
    xrt_device_current{description="12 Volts Auxillary",device_id="0000:d1:00.1",location_id="12v_aux",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 1.848
    xrt_device_current{description="12 Volts PCI Express",device_id="0000:d1:00.1",location_id="12v_pex",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 1.865
    xrt_device_current{description="3.3 Volts PCI Express",device_id="0000:d1:00.1",location_id="3v3_pex",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 1.344
    xrt_device_current{description="Internal FPGA Vcc",device_id="0000:d1:00.1",location_id="vccint",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 33.4
    xrt_device_current{description="Internal FPGA Vcc IO",device_id="0000:d1:00.1",location_id="vccint_io",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 3.9
    # HELP xrt_device_power_consumption Power consumption of the device in Watts
    # TYPE xrt_device_power_consumption gauge
    xrt_device_power_consumption{device_id="0000:d1:00.1",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 49.658128
    # HELP xrt_device_power_consumption_max Maximum power consumption of the device in watts
    # TYPE xrt_device_power_consumption_max gauge
    xrt_device_power_consumption_max{device_id="0000:d1:00.1",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 225
    # HELP xrt_device_power_consumption_warning Whether the power consumption of the device has crossed a threshold
    # TYPE xrt_device_power_consumption_warning gauge
    xrt_device_power_consumption_warning{device_id="0000:d1:00.1",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 0
    # HELP xrt_device_ready Whether the device is ready
    # TYPE xrt_device_ready gauge
    xrt_device_ready{device_id="0000:d1:00.1",shell="xilinx_u55c_gen3x16_xdma_base_3"} 1
    # HELP xrt_device_temperature Temperature of the device in degrees Celsius
    # TYPE xrt_device_temperature gauge
    xrt_device_temperature{description="Cage0",device_id="0000:d1:00.1",location_id="cage_temp_0",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 37
    xrt_device_temperature{description="FPGA",device_id="0000:d1:00.1",location_id="fpga0",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 57
    xrt_device_temperature{description="FPGA HBM",device_id="0000:d1:00.1",location_id="fpga_hbm",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 52
    xrt_device_temperature{description="Int Vcc",device_id="0000:d1:00.1",location_id="int_vcc",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 50
    xrt_device_temperature{description="PCB Top Front",device_id="0000:d1:00.1",location_id="pcb_top_front",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 42
    xrt_device_temperature{description="PCB Top Rear",device_id="0000:d1:00.1",location_id="pcb_top_rear",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 39
    # HELP xrt_device_voltage Voltage of the device in Volts
    # TYPE xrt_device_voltage gauge
    xrt_device_voltage{description="0.9 Volts Vcc",device_id="0000:d1:00.1",location_id="0v9_vcc",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 0.904
    xrt_device_voltage{description="1.2 Volts HBM",device_id="0000:d1:00.1",location_id="hbm_1v2",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 1.202
    xrt_device_voltage{description="1.8 Volts Top",device_id="0000:d1:00.1",location_id="1v8_top",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 1.80884
    xrt_device_voltage{description="12 Volts Auxillary",device_id="0000:d1:00.1",location_id="12v_aux",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 12.192
    xrt_device_voltage{description="12 Volts PCI Express",device_id="0000:d1:00.1",location_id="12v_pex",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 12.176
    xrt_device_voltage{description="3.3 Volts PCI Express",device_id="0000:d1:00.1",location_id="3v3_pex",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 3.288
    xrt_device_voltage{description="3.3 Volts Vcc",device_id="0000:d1:00.1",location_id="3v3_vcc",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 3.361
    xrt_device_voltage{description="5.5 Volts System",device_id="0000:d1:00.1",location_id="5v5_system",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 5.008
    xrt_device_voltage{description="Internal FPGA Vcc",device_id="0000:d1:00.1",location_id="vccint",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 0.853
    xrt_device_voltage{description="Internal FPGA Vcc IO",device_id="0000:d1:00.1",location_id="vccint_io",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 0.855
    xrt_device_voltage{description="Mgt Vtt",device_id="0000:d1:00.1",location_id="mgt_vtt",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 1.2
    xrt_device_voltage{description="Vpp 2.5 Volts",device_id="0000:d1:00.1",location_id="vpp2v5",serial="000000000000",shell="xilinx_u55c_gen3x16_xdma_base_3"} 2.489
    # HELP xrt_info Information about the Xilinx XRT environment
    # TYPE xrt_info gauge
    xrt_info{branch="2022.2",version="2.14.354"} 1
