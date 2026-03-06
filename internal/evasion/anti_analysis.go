package evasion

import (
	"net"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// CheckAll runs all anti-analysis checks and returns true if any analysis is detected.
func CheckAll() bool {
	if DetectDebuggerProcesses() {
		return true
	}
	if DetectVirtualization() {
		return true
	}
	return false
}

// DetectDebuggerProcesses checks for common analysis and debugger processes.
func DetectDebuggerProcesses() bool {
	debuggers := []string{
		"x64dbg.exe",
		"wireshark.exe",
		"processhacker.exe",
		"fiddler.exe", // User said 'piddler', likely meant 'fiddler'
		"vmtoolsd.exe",
		"vboxservice.exe",
		"vboxtray.exe",
	}

	procs, err := process.Processes()
	if err != nil {
		return false
	}

	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}
		name = strings.ToLower(name)
		for _, dbg := range debuggers {
			if strings.Contains(name, dbg) {
				return true
			}
		}
	}
	return false
}

// DetectVirtualization checks for common virtualization indicators.
func DetectVirtualization() bool {
	// MAC Address prefix checks
	prefixes := []string{
		"08:00:27", // VirtualBox
		"00:05:69", // VMware
		"00:0C:29", // VMware
		"00:50:56", // VMware
		"00:1C:42", // Parallels
		"00:15:5D", // Hyper-V
		"52:54:00", // QEMU/KVM
	}

	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			mac := iface.HardwareAddr.String()
			mac = strings.ToUpper(mac)
			for _, prefix := range prefixes {
				if strings.HasPrefix(mac, prefix) {
					return true
				}
			}
		}
	}

	return false
}
