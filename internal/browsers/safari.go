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
	safari = map[platform.Platform][]Browser{
		platform.Darwin: {
			{
				Name:       "Safari",
				executable: "/Applications/Safari.app/Contents/MacOS/Safari",
				Type:       "safari",
				Icon:       "https://upload.wikimedia.org/wikipedia/commons/5/52/Safari_browser_logo.svg",
			},
		},
	}
)

func safariBrowsers() []Browser {
	var browsers = make([]Browser, 0)

	// Safari is macOS only
	if runtime.GOOS != "darwin" {
		return browsers
	}

	for _, browser := range safari[platform.Platform(runtime.GOOS)] {
		// Check if browser executable exists
		if _, err := os.Stat(browser.executable); os.IsNotExist(err) {
			continue
		}

		// Extract icon from .app bundle on macOS
		appPath := "/System/Cryptexes/App/System/Applications/Safari.app"
		browser.Icon = getAppIconOrDefault(appPath, browser.Icon)

		browsers = append(browsers, browser)
	}
	return browsers
}

func safariOpen(request OpenRequest, url string) error {
	for _, browser := range browsers {
		if browser.Name == request.BrowserName {
			slog.Info("Opening Safari with URL: " + url)
			// On macOS, use 'open -a Safari' instead of calling the executable directly
			// This avoids sandboxing issues
			if runtime.GOOS == "darwin" {
				cmd := exec.Command("open", "-a", "Safari", url)
				err := cmd.Start()
				if err != nil {
					return err
				}
				return nil
			}
			// Fallback for other platforms (though Safari is macOS-only)
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
