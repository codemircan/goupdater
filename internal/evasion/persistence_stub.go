//go:build !windows
package evasion

// InstallPersistence is a stub for non-Windows platforms.
func InstallPersistence() error {
	return nil
}
