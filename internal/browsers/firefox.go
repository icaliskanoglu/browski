package browsers

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"

	"browski/internal/platform"
)

var (
	firefox = map[platform.Platform][]Browser{
		platform.Linux: {
			{
				Name:            "Firefox",
				executable:      "firefox",
				configDirectory: "firefox",
			},
		},
		platform.Darwin: {
			{
				Name:            "Firefox",
				executable:      "/Applications/Firefox.app/Contents/MacOS/firefox",
				configDirectory: "Firefox",
			},
		},
		platform.Windows: {
			{
				Name:            "Firefox",
				executable:      "C:\\Program Files\\Mozilla Firefox\\firefox.exe",
				configDirectory: "Mozilla Firefox",
			},
		},
	}
)

func firefoxBrowsers() []Browser {
	var browsers = make([]Browser, 0)
	for _, browser := range firefox[platform.Platform(runtime.GOOS)] {
		// Find the actual executable path (may be in PATH on Linux)
		execPath := findFirefoxExecutable(browser.executable)
		if execPath == "" {
			continue
		}
		browser.executable = execPath

		browser.Type = "firefox"
		browser.Icon = "https://blog.mozilla.org/design/files/2019/10/Fx-Browser-icon-fullColor.svg"

		// Extract icon from .app bundle on macOS
		if runtime.GOOS == "darwin" {
			appPath := "/Applications/Firefox.app"
			browser.Icon = getAppIconOrDefault(appPath, browser.Icon)
		}

		browsers = append(browsers, browser)
	}
	return browsers
}

// findFirefoxExecutable finds Firefox executable similar to chromium browsers
func findFirefoxExecutable(execPath string) string {
	// On Windows and macOS with absolute paths, check if file exists
	if len(execPath) > 0 && (execPath[0] == '/' || (len(execPath) > 1 && execPath[1] == ':')) {
		if _, err := os.Stat(execPath); err == nil {
			return execPath
		}
		// On Windows, try alternative locations
		if runtime.GOOS == "windows" {
			// Try Program Files (x86)
			if _, err := os.Stat("C:\\Program Files (x86)\\Mozilla Firefox\\firefox.exe"); err == nil {
				return "C:\\Program Files (x86)\\Mozilla Firefox\\firefox.exe"
			}
		}
		return ""
	}

	// For command names (common on Linux), use LookPath
	if path, err := exec.LookPath(execPath); err == nil {
		return path
	}

	return ""
}

func firefoxOpen(request OpenRequest, url string) error {
	for _, browser := range browsers {
		if browser.Name == request.BrowserName {
			slog.Warn("Opening browser: " + browser.executable + " " + url)
			cmd := exec.Command(browser.executable, url)
			err := cmd.Start()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("browser not found: %s", request.BrowserName)
}
