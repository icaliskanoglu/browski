package browsers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"browski/internal/platform"

	"github.com/PaesslerAG/jsonpath"
)

var (
	chrome = map[platform.Platform][]Browser{
		platform.Linux: {
			{
				Name:            "Google Chrome",
				executable:      "google-chrome-stable",
				configDirectory: "google-chrome",
				Type:            "chromium",
				Icon:            "https://www.google.com/chrome/static/images/chrome-logo.svg",
			},
			{
				Name:            "Brave",
				executable:      "brave",
				configDirectory: "BraveSoftware/Brave-Browser",
				Type:            "chromium",
				Icon:            "https://brave.com/static-assets/images/brave-logo-sans-text.svg",
			},
			{
				Name:            "Microsoft Edge",
				executable:      "microsoft-edge-stable",
				configDirectory: "microsoft-edge",
				Type:            "chromium",
				Icon:            "https://upload.wikimedia.org/wikipedia/commons/9/98/Microsoft_Edge_logo_%282019%29.svg",
			},
			{
				Name:            "Opera",
				executable:      "opera",
				configDirectory: "opera",
				Type:            "chromium",
				Icon:            "https://upload.wikimedia.org/wikipedia/commons/4/49/Opera_2015_icon.svg",
			},
			{
				Name:            "Vivaldi",
				executable:      "vivaldi-stable",
				configDirectory: "vivaldi",
				Type:            "chromium",
				Icon:            "https://vivaldi.com/wp-content/themes/vivaldicom-theme/img/press/icons/vivaldi_icon.png",
			},
		},
		platform.Darwin: {
			{
				Name:            "Google Chrome",
				executable:      "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
				configDirectory: "Google/Chrome",
				Type:            "chromium",
				Icon:            "https://www.google.com/chrome/static/images/chrome-logo.svg",
			},
			{
				Name:            "Microsoft Edge",
				executable:      "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
				configDirectory: "Microsoft Edge",
				Type:            "chromium",
				Icon:            "https://upload.wikimedia.org/wikipedia/commons/9/98/Microsoft_Edge_logo_%282019%29.svg",
			},
			{
				Name:            "Brave Browser",
				executable:      "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
				configDirectory: "BraveSoftware/Brave-Browser",
				Type:            "chromium",
				Icon:            "https://brave.com/static-assets/images/brave-logo-sans-text.svg",
			},
			{
				Name:            "Arc",
				executable:      "/Applications/Arc.app/Contents/MacOS/Arc",
				configDirectory: "Arc",
				Type:            "chromium",
				Icon:            "https://arc.net/icon.png",
			},
			{
				Name:            "Opera",
				executable:      "/Applications/Opera.app/Contents/MacOS/Opera",
				configDirectory: "com.operasoftware.Opera",
				Type:            "chromium",
				Icon:            "https://upload.wikimedia.org/wikipedia/commons/4/49/Opera_2015_icon.svg",
			},
			{
				Name:            "Vivaldi",
				executable:      "/Applications/Vivaldi.app/Contents/MacOS/Vivaldi",
				configDirectory: "Vivaldi",
				Type:            "chromium",
				Icon:            "https://vivaldi.com/wp-content/themes/vivaldicom-theme/img/press/icons/vivaldi_icon.png",
			},
		},
		platform.Windows: {
			{
				Type:            "chromium",
				Name:            "Google Chrome",
				executable:      "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
				configDirectory: "Google\\Chrome\\User Data",
				Icon:            "https://www.google.com/chrome/static/images/chrome-logo.svg",
			},
			{
				Type:            "chromium",
				Name:            "Microsoft Edge",
				executable:      "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
				configDirectory: "Microsoft\\Edge\\User Data",
				Icon:            "https://upload.wikimedia.org/wikipedia/commons/9/98/Microsoft_Edge_logo_%282019%29.svg",
			},
			{
				Type:            "chromium",
				Name:            "Brave",
				executable:      "C:\\Program Files\\BraveSoftware\\Brave-Browser\\Application\\brave.exe",
				configDirectory: "BraveSoftware\\Brave-Browser\\User Data",
				Icon:            "https://brave.com/static-assets/images/brave-logo-sans-text.svg",
			},
			{
				Type:            "chromium",
				Name:            "Opera",
				executable:      "C:\\Program Files\\Opera\\launcher.exe",
				configDirectory: "Opera Software\\Opera Stable",
				Icon:            "https://upload.wikimedia.org/wikipedia/commons/4/49/Opera_2015_icon.svg",
			},
			{
				Type:            "chromium",
				Name:            "Vivaldi",
				executable:      "C:\\Program Files\\Vivaldi\\Application\\vivaldi.exe",
				configDirectory: "Vivaldi\\User Data",
				Icon:            "https://vivaldi.com/wp-content/themes/vivaldicom-theme/img/press/icons/vivaldi_icon.png",
			},
		},
	}
)

func chromiumBrowsers(userConfigDir string) []Browser {
	var browsers = make([]Browser, 0)
	for _, browser := range chrome[platform.Platform(runtime.GOOS)] {
		// Find the actual executable path (may be in PATH on Linux/Windows)
		execPath := findExecutable(browser.executable)
		if execPath == "" {
			continue
		}
		browser.executable = execPath

		// Extract icon from .app bundle on macOS
		if runtime.GOOS == "darwin" {
			// Determine the .app path from the executable path
			// e.g., "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" -> "/Applications/Google Chrome.app"
			exePath := browser.executable
			if len(exePath) > 0 && exePath[0] == '/' {
				// Find the .app directory
				for i := 0; i < len(exePath); i++ {
					if i+4 < len(exePath) && exePath[i:i+4] == ".app" {
						appPath := exePath[:i+4]
						browser.Icon = getAppIconOrDefault(appPath, browser.Icon)
						break
					}
				}
			}
		}

		browser.Profiles = chromeProfiles(userConfigDir, browser)
		browsers = append(browsers, browser)
	}
	return browsers
}

// findExecutable finds the full path to an executable
// For absolute paths (Windows/macOS), checks if the file exists
// For command names (Linux), uses exec.LookPath to find it in PATH
func findExecutable(execPath string) string {
	// On Windows and macOS with absolute paths, check if file exists
	if filepath.IsAbs(execPath) {
		if _, err := os.Stat(execPath); err == nil {
			return execPath
		}
		// On Windows, also check alternative locations
		if runtime.GOOS == "windows" {
			// Try Program Files (x86)
			alt := strings.Replace(execPath, "Program Files", "Program Files (x86)", 1)
			if _, err := os.Stat(alt); err == nil {
				return alt
			}
			// Try user's local AppData
			if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
				// Extract just the app name and path after Program Files
				if idx := strings.Index(execPath, "Program Files\\"); idx != -1 {
					relPath := execPath[idx+len("Program Files\\"):]
					localPath := filepath.Join(localAppData, relPath)
					if _, err := os.Stat(localPath); err == nil {
						return localPath
					}
				}
			}
		}
		return ""
	}

	// For relative paths or command names (common on Linux), use LookPath
	if path, err := exec.LookPath(execPath); err == nil {
		return path
	}

	return ""
}

func chromiumOpen(request OpenRequest, url string) error {
	for _, browser := range browsers {
		if browser.Name == request.BrowserName {
			var profileDirectory string
			if request.ProfileName == "" {
				profileDirectory = "Default"
			} else {
				for _, profile := range browser.Profiles {
					if profile.Name == request.ProfileName {
						profileDirectory = profile.Directory
						break
					}
				}
			}
			cmd := exec.Command(browser.executable, "--profile-directory="+profileDirectory, url)
			_, err := cmd.Output()
			if err != nil {
				slog.With(err).Error("Could not run command")
				return fmt.Errorf("browser open url with profile: %s, url: %s", request.BrowserName, url)
			}
			return nil
		}
	}
	return fmt.Errorf("browser not found: %s", request.BrowserName)
}
func chromeProfiles(userConfigDir string, browser Browser) []Profile {

	var profiles []Profile
	statePath := filepath.Join(userConfigDir, browser.configDirectory, "Local State")
	if fileContent, err := os.ReadFile(statePath); err == nil {
		var p interface{}
		if err := json.Unmarshal(fileContent, &p); err == nil {
			if profilesMap, err := jsonpath.Get("$.profile.info_cache", p); err == nil {
				for s, i := range profilesMap.(map[string]interface{}) {
					profile := Profile{
						Name:      i.(map[string]interface{})["name"].(string),
						Directory: s,
					}
					if iconUrl, ok := i.(map[string]interface{})["last_downloaded_gaia_picture_url_with_size"]; ok && iconUrl != "" {
						profile.Icon = iconUrl.(string)
					} else {
						profile.Icon = browser.Icon
					}
					profiles = append(profiles, profile)
				}
			} else {
				slog.With(err).Error("Error getting profile info_cache")
			}
		} else {
			slog.With(err).Error("Error unmarshalling json")
		}

	} else {
		// Only log error if it's not a "file not found" error
		// (file not found is normal for browsers that haven't been used yet)
		if !os.IsNotExist(err) {
			slog.With(err).Error("Error reading file" + statePath)
		}
	}
	return profiles
}
