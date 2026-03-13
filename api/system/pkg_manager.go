package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// Package represents a FreeBSD pkg entry
type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// ListInstalledPackages wraps 'pkg info'
func ListInstalledPackages() ([]Package, error) {
	cmd := exec.Command("pkg", "info")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("pkg info failed: %v", err)
	}

	var packages []Package
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line == "" { continue }
		// pkg info format is usually "name-version description" or just "name-version"
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			pkgInfo := fields[0]
			lastDash := strings.LastIndex(pkgInfo, "-")
			name := pkgInfo
			version := "unknown"
			if lastDash > 0 {
				name = pkgInfo[:lastDash]
				version = pkgInfo[lastDash+1:]
			}
			
			desc := ""
			if len(fields) > 1 {
				desc = strings.Join(fields[1:], " ")
			}
			
			packages = append(packages, Package{
				Name:        name,
				Version:     version,
				Description: desc,
			})
		}
	}
	return packages, nil
}

// InstallPackage adds a new extension
func InstallPackage(name string) error {
	cmd := exec.Command("pkg", "install", "-y", name)
	return cmd.Run()
}

// RemovePackage uninstalls an extension
func RemovePackage(name string) error {
	cmd := exec.Command("pkg", "delete", "-y", name)
	return cmd.Run()
}
