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
    xrt_device_current{description="12 Volts Auxillary",device_id="0000:56:00.1",location_id="12v_aux",serial="AAAAAAAAAAAA"} 1.865
    xrt_device_current{description="12 Volts Auxillary",device_id="0000:57:00.1",location_id="12v_aux",serial="BBBBBBBBBBBB"} 1.64
    xrt_device_current{description="12 Volts Auxillary",device_id="0000:ce:00.1",location_id="12v_aux",serial="CCCCCCCCCCCC"} 1.312
    xrt_device_current{description="12 Volts Auxillary",device_id="0000:d1:00.1",location_id="12v_aux",serial="DDDDDDDDDDDD"} 1.344
    xrt_device_current{description="12 Volts PCI Express",device_id="0000:56:00.1",location_id="12v_pex",serial="AAAAAAAAAAAA"} 1.904
    xrt_device_current{description="12 Volts PCI Express",device_id="0000:57:00.1",location_id="12v_pex",serial="BBBBBBBBBBBB"} 1.808
    xrt_device_current{description="12 Volts PCI Express",device_id="0000:ce:00.1",location_id="12v_pex",serial="CCCCCCCCCCCC"} 1.704
    xrt_device_current{description="12 Volts PCI Express",device_id="0000:d1:00.1",location_id="12v_pex",serial="DDDDDDDDDDDD"} 1.721
    xrt_device_current{description="3.3 Volts PCI Express",device_id="0000:56:00.1",location_id="3v3_pex",serial="AAAAAAAAAAAA"} 1.344
    xrt_device_current{description="3.3 Volts PCI Express",device_id="0000:57:00.1",location_id="3v3_pex",serial="BBBBBBBBBBBB"} 1.328
    xrt_device_current{description="3.3 Volts PCI Express",device_id="0000:ce:00.1",location_id="3v3_pex",serial="CCCCCCCCCCCC"} 1.552
    xrt_device_current{description="3.3 Volts PCI Express",device_id="0000:d1:00.1",location_id="3v3_pex",serial="DDDDDDDDDDDD"} 1.56
    xrt_device_current{description="Internal FPGA Vcc",device_id="0000:56:00.1",location_id="vccint",serial="AAAAAAAAAAAA"} 32.401
    xrt_device_current{description="Internal FPGA Vcc",device_id="0000:57:00.1",location_id="vccint",serial="BBBBBBBBBBBB"} 30.3
    xrt_device_current{description="Internal FPGA Vcc",device_id="0000:ce:00.1",location_id="vccint",serial="CCCCCCCCCCCC"} 21.2
    xrt_device_current{description="Internal FPGA Vcc",device_id="0000:d1:00.1",location_id="vccint",serial="DDDDDDDDDDDD"} 23.2
    xrt_device_current{description="Internal FPGA Vcc IO",device_id="0000:56:00.1",location_id="vccint_io",serial="AAAAAAAAAAAA"} 3.8
    xrt_device_current{description="Internal FPGA Vcc IO",device_id="0000:57:00.1",location_id="vccint_io",serial="BBBBBBBBBBBB"} 3.8
    xrt_device_current{description="Internal FPGA Vcc IO",device_id="0000:ce:00.1",location_id="vccint_io",serial="CCCCCCCCCCCC"} 3.5
    xrt_device_current{description="Internal FPGA Vcc IO",device_id="0000:d1:00.1",location_id="vccint_io",serial="DDDDDDDDDDDD"} 3.9
    # HELP xrt_device_power_consumption Power consumption of the device in Watts
    # TYPE xrt_device_power_consumption gauge
    xrt_device_power_consumption{device_id="0000:56:00.1",serial="AAAAAAAAAAAA"} 50.355176
    xrt_device_power_consumption{device_id="0000:57:00.1",serial="BBBBBBBBBBBB"} 46.388672
    xrt_device_power_consumption{device_id="0000:ce:00.1",serial="CCCCCCCCCCCC"} 41.820736
    xrt_device_power_consumption{device_id="0000:d1:00.1",serial="DDDDDDDDDDDD"} 42.470224
    # HELP xrt_device_power_consumption_max Maximum power consumption of the device in watts
    # TYPE xrt_device_power_consumption_max gauge
    xrt_device_power_consumption_max{device_id="0000:56:00.1",serial="AAAAAAAAAAAA"} 225
    xrt_device_power_consumption_max{device_id="0000:57:00.1",serial="BBBBBBBBBBBB"} 225
    xrt_device_power_consumption_max{device_id="0000:ce:00.1",serial="CCCCCCCCCCCC"} 225
    xrt_device_power_consumption_max{device_id="0000:d1:00.1",serial="DDDDDDDDDDDD"} 225
    # HELP xrt_device_power_consumption_warning Whether the power consumption of the device has crossed a threshold
    # TYPE xrt_device_power_consumption_warning gauge
    xrt_device_power_consumption_warning{device_id="0000:56:00.1",serial="AAAAAAAAAAAA"} 0
    xrt_device_power_consumption_warning{device_id="0000:57:00.1",serial="BBBBBBBBBBBB"} 0
    xrt_device_power_consumption_warning{device_id="0000:ce:00.1",serial="CCCCCCCCCCCC"} 0
    xrt_device_power_consumption_warning{device_id="0000:d1:00.1",serial="DDDDDDDDDDDD"} 0
    # HELP xrt_device_temperature Temperature of the device in degrees Celsius
    # TYPE xrt_device_temperature gauge
    xrt_device_temperature{description="Cage0",device_id="0000:56:00.1",location_id="cage_temp_0",serial="AAAAAAAAAAAA"} 37
    xrt_device_temperature{description="Cage0",device_id="0000:57:00.1",location_id="cage_temp_0",serial="BBBBBBBBBBBB"} 31
    xrt_device_temperature{description="Cage0",device_id="0000:ce:00.1",location_id="cage_temp_0",serial="CCCCCCCCCCCC"} 33
    xrt_device_temperature{description="Cage0",device_id="0000:d1:00.1",location_id="cage_temp_0",serial="DDDDDDDDDDDD"} 35
    xrt_device_temperature{description="FPGA",device_id="0000:56:00.1",location_id="fpga0",serial="AAAAAAAAAAAA"} 55
    xrt_device_temperature{description="FPGA",device_id="0000:57:00.1",location_id="fpga0",serial="BBBBBBBBBBBB"} 53
    xrt_device_temperature{description="FPGA",device_id="0000:ce:00.1",location_id="fpga0",serial="CCCCCCCCCCCC"} 49
    xrt_device_temperature{description="FPGA",device_id="0000:d1:00.1",location_id="fpga0",serial="DDDDDDDDDDDD"} 52
    xrt_device_temperature{description="FPGA HBM",device_id="0000:56:00.1",location_id="fpga_hbm",serial="AAAAAAAAAAAA"} 50
    xrt_device_temperature{description="FPGA HBM",device_id="0000:57:00.1",location_id="fpga_hbm",serial="BBBBBBBBBBBB"} 48
    xrt_device_temperature{description="FPGA HBM",device_id="0000:ce:00.1",location_id="fpga_hbm",serial="CCCCCCCCCCCC"} 44
    xrt_device_temperature{description="FPGA HBM",device_id="0000:d1:00.1",location_id="fpga_hbm",serial="DDDDDDDDDDDD"} 47
    xrt_device_temperature{description="Int Vcc",device_id="0000:56:00.1",location_id="int_vcc",serial="AAAAAAAAAAAA"} 51
    xrt_device_temperature{description="Int Vcc",device_id="0000:57:00.1",location_id="int_vcc",serial="BBBBBBBBBBBB"} 45
    xrt_device_temperature{description="Int Vcc",device_id="0000:ce:00.1",location_id="int_vcc",serial="CCCCCCCCCCCC"} 44
    xrt_device_temperature{description="Int Vcc",device_id="0000:d1:00.1",location_id="int_vcc",serial="DDDDDDDDDDDD"} 45
    xrt_device_temperature{description="PCB Top Front",device_id="0000:56:00.1",location_id="pcb_top_front",serial="AAAAAAAAAAAA"} 40
    xrt_device_temperature{description="PCB Top Front",device_id="0000:57:00.1",location_id="pcb_top_front",serial="BBBBBBBBBBBB"} 35
    xrt_device_temperature{description="PCB Top Front",device_id="0000:ce:00.1",location_id="pcb_top_front",serial="CCCCCCCCCCCC"} 39
    xrt_device_temperature{description="PCB Top Front",device_id="0000:d1:00.1",location_id="pcb_top_front",serial="DDDDDDDDDDDD"} 39
    xrt_device_temperature{description="PCB Top Rear",device_id="0000:56:00.1",location_id="pcb_top_rear",serial="AAAAAAAAAAAA"} 39
    xrt_device_temperature{description="PCB Top Rear",device_id="0000:57:00.1",location_id="pcb_top_rear",serial="BBBBBBBBBBBB"} 33
    xrt_device_temperature{description="PCB Top Rear",device_id="0000:ce:00.1",location_id="pcb_top_rear",serial="CCCCCCCCCCCC"} 34
    xrt_device_temperature{description="PCB Top Rear",device_id="0000:d1:00.1",location_id="pcb_top_rear",serial="DDDDDDDDDDDD"} 36
    # HELP xrt_device_voltage Voltage of the device in Volts
    # TYPE xrt_device_voltage gauge
    xrt_device_voltage{description="0.9 Volts Vcc",device_id="0000:56:00.1",location_id="0v9_vcc",serial="AAAAAAAAAAAA"} 0.903
    xrt_device_voltage{description="0.9 Volts Vcc",device_id="0000:57:00.1",location_id="0v9_vcc",serial="BBBBBBBBBBBB"} 0.898
    xrt_device_voltage{description="0.9 Volts Vcc",device_id="0000:ce:00.1",location_id="0v9_vcc",serial="CCCCCCCCCCCC"} 0.9
    xrt_device_voltage{description="0.9 Volts Vcc",device_id="0000:d1:00.1",location_id="0v9_vcc",serial="DDDDDDDDDDDD"} 0.906
    xrt_device_voltage{description="1.2 Volts HBM",device_id="0000:56:00.1",location_id="hbm_1v2",serial="AAAAAAAAAAAA"} 1.203
    xrt_device_voltage{description="1.2 Volts HBM",device_id="0000:57:00.1",location_id="hbm_1v2",serial="BBBBBBBBBBBB"} 1.202
    xrt_device_voltage{description="1.2 Volts HBM",device_id="0000:ce:00.1",location_id="hbm_1v2",serial="CCCCCCCCCCCC"} 1.204
    xrt_device_voltage{description="1.2 Volts HBM",device_id="0000:d1:00.1",location_id="hbm_1v2",serial="DDDDDDDDDDDD"} 1.202
    xrt_device_voltage{description="1.8 Volts Top",device_id="0000:56:00.1",location_id="1v8_top",serial="AAAAAAAAAAAA"} 1.796
    xrt_device_voltage{description="1.8 Volts Top",device_id="0000:57:00.1",location_id="1v8_top",serial="BBBBBBBBBBBB"} 1.803
    xrt_device_voltage{description="1.8 Volts Top",device_id="0000:ce:00.1",location_id="1v8_top",serial="CCCCCCCCCCCC"} 1.814
    xrt_device_voltage{description="1.8 Volts Top",device_id="0000:d1:00.1",location_id="1v8_top",serial="DDDDDDDDDDDD"} 1.808
    xrt_device_voltage{description="12 Volts Auxillary",device_id="0000:56:00.1",location_id="12v_aux",serial="AAAAAAAAAAAA"} 12.2
    xrt_device_voltage{description="12 Volts Auxillary",device_id="0000:57:00.1",location_id="12v_aux",serial="BBBBBBBBBBBB"} 12.2
    xrt_device_voltage{description="12 Volts Auxillary",device_id="0000:ce:00.1",location_id="12v_aux",serial="CCCCCCCCCCCC"} 12.192
    xrt_device_voltage{description="12 Volts Auxillary",device_id="0000:d1:00.1",location_id="12v_aux",serial="DDDDDDDDDDDD"} 12.192
    xrt_device_voltage{description="12 Volts PCI Express",device_id="0000:56:00.1",location_id="12v_pex",serial="AAAAAAAAAAAA"} 12.176
    xrt_device_voltage{description="12 Volts PCI Express",device_id="0000:57:00.1",location_id="12v_pex",serial="BBBBBBBBBBBB"} 12.176
    xrt_device_voltage{description="12 Volts PCI Express",device_id="0000:ce:00.1",location_id="12v_pex",serial="CCCCCCCCCCCC"} 12.168
    xrt_device_voltage{description="12 Volts PCI Express",device_id="0000:d1:00.1",location_id="12v_pex",serial="DDDDDDDDDDDD"} 12.176
    xrt_device_voltage{description="3.3 Volts PCI Express",device_id="0000:56:00.1",location_id="3v3_pex",serial="AAAAAAAAAAAA"} 3.288
    xrt_device_voltage{description="3.3 Volts PCI Express",device_id="0000:57:00.1",location_id="3v3_pex",serial="BBBBBBBBBBBB"} 3.288
    xrt_device_voltage{description="3.3 Volts PCI Express",device_id="0000:ce:00.1",location_id="3v3_pex",serial="CCCCCCCCCCCC"} 3.28
    xrt_device_voltage{description="3.3 Volts PCI Express",device_id="0000:d1:00.1",location_id="3v3_pex",serial="DDDDDDDDDDDD"} 3.288
    xrt_device_voltage{description="3.3 Volts Vcc",device_id="0000:56:00.1",location_id="3v3_vcc",serial="AAAAAAAAAAAA"} 3.348
    xrt_device_voltage{description="3.3 Volts Vcc",device_id="0000:57:00.1",location_id="3v3_vcc",serial="BBBBBBBBBBBB"} 3.357
    xrt_device_voltage{description="3.3 Volts Vcc",device_id="0000:ce:00.1",location_id="3v3_vcc",serial="CCCCCCCCCCCC"} 3.348
    xrt_device_voltage{description="3.3 Volts Vcc",device_id="0000:d1:00.1",location_id="3v3_vcc",serial="DDDDDDDDDDDD"} 3.36
    xrt_device_voltage{description="5.5 Volts System",device_id="0000:56:00.1",location_id="5v5_system",serial="AAAAAAAAAAAA"} 4.989
    xrt_device_voltage{description="5.5 Volts System",device_id="0000:57:00.1",location_id="5v5_system",serial="BBBBBBBBBBBB"} 5
    xrt_device_voltage{description="5.5 Volts System",device_id="0000:ce:00.1",location_id="5v5_system",serial="CCCCCCCCCCCC"} 5.034
    xrt_device_voltage{description="5.5 Volts System",device_id="0000:d1:00.1",location_id="5v5_system",serial="DDDDDDDDDDDD"} 5.01
    xrt_device_voltage{description="Internal FPGA Vcc",device_id="0000:56:00.1",location_id="vccint",serial="AAAAAAAAAAAA"} 0.853
    xrt_device_voltage{description="Internal FPGA Vcc",device_id="0000:57:00.1",location_id="vccint",serial="BBBBBBBBBBBB"} 0.854
    xrt_device_voltage{description="Internal FPGA Vcc",device_id="0000:ce:00.1",location_id="vccint",serial="CCCCCCCCCCCC"} 0.854
    xrt_device_voltage{description="Internal FPGA Vcc",device_id="0000:d1:00.1",location_id="vccint",serial="DDDDDDDDDDDD"} 0.852
    xrt_device_voltage{description="Internal FPGA Vcc IO",device_id="0000:56:00.1",location_id="vccint_io",serial="AAAAAAAAAAAA"} 0.854
    xrt_device_voltage{description="Internal FPGA Vcc IO",device_id="0000:57:00.1",location_id="vccint_io",serial="BBBBBBBBBBBB"} 0.855
    xrt_device_voltage{description="Internal FPGA Vcc IO",device_id="0000:ce:00.1",location_id="vccint_io",serial="CCCCCCCCCCCC"} 0.856
    xrt_device_voltage{description="Internal FPGA Vcc IO",device_id="0000:d1:00.1",location_id="vccint_io",serial="DDDDDDDDDDDD"} 0.854
    xrt_device_voltage{description="Mgt Vtt",device_id="0000:56:00.1",location_id="mgt_vtt",serial="AAAAAAAAAAAA"} 1.205
    xrt_device_voltage{description="Mgt Vtt",device_id="0000:57:00.1",location_id="mgt_vtt",serial="BBBBBBBBBBBB"} 1.204
    xrt_device_voltage{description="Mgt Vtt",device_id="0000:ce:00.1",location_id="mgt_vtt",serial="CCCCCCCCCCCC"} 1.205
    xrt_device_voltage{description="Mgt Vtt",device_id="0000:d1:00.1",location_id="mgt_vtt",serial="DDDDDDDDDDDD"} 1.201
    xrt_device_voltage{description="Vpp 2.5 Volts",device_id="0000:56:00.1",location_id="vpp2v5",serial="AAAAAAAAAAAA"} 2.504
    xrt_device_voltage{description="Vpp 2.5 Volts",device_id="0000:57:00.1",location_id="vpp2v5",serial="BBBBBBBBBBBB"} 2.513
    xrt_device_voltage{description="Vpp 2.5 Volts",device_id="0000:ce:00.1",location_id="vpp2v5",serial="CCCCCCCCCCCC"} 2.507
    xrt_device_voltage{description="Vpp 2.5 Volts",device_id="0000:d1:00.1",location_id="vpp2v5",serial="DDDDDDDDDDDD"} 2.489
