// +build darwin

package browsers

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// extractAppIcon extracts an icon from a macOS .app bundle and converts it to a base64 data URL
func extractAppIcon(appPath string) (string, error) {
	// Find the .icns file in the app bundle
	icnsPath := filepath.Join(appPath, "Contents", "Resources", "AppIcon.icns")

	// Check if the icon exists
	if _, err := os.Stat(icnsPath); os.IsNotExist(err) {
		// Try alternative icon names
		resourcesPath := filepath.Join(appPath, "Contents", "Resources")
		entries, err := os.ReadDir(resourcesPath)
		if err != nil {
			return "", fmt.Errorf("failed to read resources directory: %w", err)
		}

		// Find first .icns file
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) == ".icns" {
				icnsPath = filepath.Join(resourcesPath, entry.Name())
				break
			}
		}
	}

	// Create temporary PNG file
	tmpDir := os.TempDir()
	pngPath := filepath.Join(tmpDir, fmt.Sprintf("browski-icon-%d.png", os.Getpid()))
	defer os.Remove(pngPath)

	// Convert .icns to PNG using sips (macOS built-in image tool)
	// Extract 64x64 version which is good for web display
	cmd := exec.Command("sips", "-s", "format", "png", icnsPath, "--out", pngPath, "-Z", "64")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to convert icon: %w", err)
	}

	// Read the PNG file
	pngData, err := os.ReadFile(pngPath)
	if err != nil {
		return "", fmt.Errorf("failed to read PNG: %w", err)
	}

	// Convert to base64 data URL
	b64 := base64.StdEncoding.EncodeToString(pngData)
	dataURL := fmt.Sprintf("data:image/png;base64,%s", b64)

	return dataURL, nil
}

// getAppIconOrDefault attempts to extract the app icon, falls back to URL if it fails
func getAppIconOrDefault(appPath, defaultURL string) string {
	if icon, err := extractAppIcon(appPath); err == nil {
		return icon
	}
	return defaultURL
}
