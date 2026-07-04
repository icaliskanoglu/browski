package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"log/slog"
)

const (
	// GitHub API endpoint for latest release
	releaseURL = "https://api.github.com/repos/icaliskanoglu/browski/releases/latest"

	// Current version - update this with each release
	currentVersion = "1.0.0"
)

// Release represents a GitHub release
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Body        string    `json:"body"`
	Assets      []Asset   `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	Available      bool
	Version        string
	ReleaseNotes   string
	DownloadURL    string
	CurrentVersion string
}

// CheckForUpdates checks if a new version is available
func CheckForUpdates() (*UpdateInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", releaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add user agent to avoid GitHub rate limiting
	req.Header.Set("User-Agent", "Browski-Update-Checker")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Remove 'v' prefix if present
	latestVersion := release.TagName
	if len(latestVersion) > 0 && latestVersion[0] == 'v' {
		latestVersion = latestVersion[1:]
	}

	info := &UpdateInfo{
		CurrentVersion: currentVersion,
		Version:        latestVersion,
		ReleaseNotes:   release.Body,
		Available:      latestVersion > currentVersion,
	}

	// Find the appropriate download URL for this platform
	if info.Available {
		info.DownloadURL = findDownloadURL(&release)
	}

	return info, nil
}

// findDownloadURL finds the appropriate download URL for the current platform
func findDownloadURL(release *Release) string {
	var pattern string

	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			pattern = "macos-arm64.dmg"
		} else {
			pattern = "macos-amd64.dmg"
		}
	case "windows":
		pattern = "windows-amd64-setup.exe"
	case "linux":
		pattern = "linux-amd64.AppImage"
	default:
		return release.HTMLURL
	}

	// Find matching asset
	for _, asset := range release.Assets {
		if contains(asset.Name, pattern) {
			return asset.BrowserDownloadURL
		}
	}

	// Fallback to release page
	return release.HTMLURL
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CheckInBackground checks for updates in the background and logs results
func CheckInBackground() {
	go func() {
		time.Sleep(5 * time.Second) // Wait 5 seconds after app start

		info, err := CheckForUpdates()
		if err != nil {
			slog.With("error", err).Debug("Failed to check for updates")
			return
		}

		if info.Available {
			slog.With(
				"current", info.CurrentVersion,
				"latest", info.Version,
				"url", info.DownloadURL,
			).Info("Update available")
		} else {
			slog.With("version", info.CurrentVersion).Debug("Running latest version")
		}
	}()
}

// GetCurrentVersion returns the current application version
func GetCurrentVersion() string {
	return currentVersion
}
