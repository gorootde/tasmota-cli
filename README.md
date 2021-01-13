# tasmota-cli
CLI utility to manage tasmota devices from the command line.

![CI](https://github.com/gorootde/tasmota-cli/workflows/CI/badge.svg) [![Open Source Love](https://badges.frapsoft.com/os/gpl/gpl.svg?v=102)](https://github.com/ellerbrock/open-source-badge/)

# Usage
```bash
Usage: tasmota-cli [OPTIONS] <command> <ip> (<ip> <ip> ...)

  <command> 	Any tasmota command. See https://tasmota.github.io/docs/Commands/
  <ip> 		List of the IPs of the tasmota devices to execute the command on

CLI Arguments always take precedence over environment eariables
  TASMOTACLI_USERNAME  Username
  TASMOTACLI_PASSWORD  Password

Flags:
  -la
    	Enable legacy authentication mode (for tasmota versions <= 9.2.0)
  -p string
    	The password used to authenticate to tasmota
  -u string
    	The username used to authenticate to tasmota
  -v	Enable verbose mode
```

# Features
- Execute [commands](https://tasmota.github.io/docs/Commands/) on one or more tasmota devices with a single CLI call

** Coming soon **
- Scan for Tasmota devices in your network (Discovery)
