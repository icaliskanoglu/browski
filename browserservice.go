package main

import (
	"log/slog"
	"sort"

	"browski/internal/browsers"
	"browski/internal/preferences"
	"browski/internal/updater"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// BrowserService provides browser selection functionality
type BrowserService struct {
	app            *application.App
	url            string
	prefs          *preferences.Preferences
	chooserWindow  *application.WebviewWindow
	prefsWindow    *application.WebviewWindow
}

// NewBrowserService creates a new BrowserService
func NewBrowserService(app *application.App) *BrowserService {
	prefs, err := preferences.Load()
	if err != nil {
		slog.With("error", err).Warn("Failed to load preferences, using defaults")
		prefs = &preferences.Preferences{
			HiddenBrowsers: make(map[string]bool),
		}
	}

	return &BrowserService{
		app:   app,
		prefs: prefs,
	}
}

// SetWindows sets references to the app windows
func (bs *BrowserService) SetWindows(chooser, prefs *application.WebviewWindow) {
	bs.chooserWindow = chooser
	bs.prefsWindow = prefs
}

// SetURL sets a new URL and emits the url-changed event
// The frontend will show the window when ready
func (bs *BrowserService) SetURL(newURL string) {
	slog.With("url", newURL).Info("SetURL called")
	bs.url = newURL
	bs.app.Event.Emit("url-changed", newURL)
	slog.Info("url-changed event emitted")
}

// ListBrowsers returns all detected browsers filtered by preferences
func (bs *BrowserService) ListBrowsers() []browsers.Browser {
	allBrowsers := browsers.ListBrowsers()

	// Filter based on preferences
	filteredBrowsers := make([]browsers.Browser, 0)
	for _, browser := range allBrowsers {
		// For browsers with profiles, the main browser entry is treated separately
		// For browsers without profiles, check if the main browser is hidden
		includeMainBrowser := !bs.prefs.IsBrowserHidden(browser.Name, "")

		// Filter profiles
		visibleProfiles := make([]browsers.Profile, 0)
		for _, profile := range browser.Profiles {
			if !bs.prefs.IsBrowserHidden(browser.Name, profile.Name) {
				visibleProfiles = append(visibleProfiles, profile)
			}
		}

		// Decide whether to include this browser:
		// - If browser has no profiles: include only if main browser is not hidden
		// - If browser has profiles: include if either main browser is visible OR at least one profile is visible
		if len(browser.Profiles) == 0 {
			// No profiles: include only if main browser is not hidden
			if !includeMainBrowser {
				continue
			}
		} else {
			// Has profiles: include if main browser is visible OR at least one profile is visible
			if !includeMainBrowser && len(visibleProfiles) == 0 {
				continue
			}
		}

		// Sort profiles alphabetically by name
		sort.Slice(visibleProfiles, func(i, j int) bool {
			return visibleProfiles[i].Name < visibleProfiles[j].Name
		})

		browser.Profiles = visibleProfiles
		browser.ShowMainEntry = includeMainBrowser
		filteredBrowsers = append(filteredBrowsers, browser)
	}

	// Sort browsers alphabetically by name
	sort.Slice(filteredBrowsers, func(i, j int) bool {
		return filteredBrowsers[i].Name < filteredBrowsers[j].Name
	})

	return filteredBrowsers
}

// ListAllBrowsers returns all detected browsers without filtering
func (bs *BrowserService) ListAllBrowsers() []browsers.Browser {
	return browsers.ListBrowsers()
}

// GetURL returns the current URL
func (bs *BrowserService) GetURL() string {
	slog.With("url", bs.url).Info("GetURL called from frontend")
	return bs.url
}

// Resize resizes the chooser window
func (bs *BrowserService) Resize(width int, height int) {
	if bs.chooserWindow != nil {
		bs.chooserWindow.SetSize(width, height)
	}
}

// Open opens the selected browser with the URL
func (bs *BrowserService) Open(request browsers.OpenRequest) {
	if err := browsers.Open(request, bs.url); err != nil {
		slog.With("error", err).Error("Failed to open browser")
	}
	if bs.chooserWindow != nil {
		bs.chooserWindow.Hide()
	}
}

// HideWindow hides the chooser window
func (bs *BrowserService) HideWindow() {
	if bs.chooserWindow != nil {
		bs.chooserWindow.Hide()
	}
}

// ShowWindow shows and centers the chooser window
func (bs *BrowserService) ShowWindow() {
	if bs.chooserWindow != nil {
		bs.chooserWindow.Show()
		bs.chooserWindow.Center()
		bs.chooserWindow.Focus()
	}
}

// GetPreferences returns the current preferences
func (bs *BrowserService) GetPreferences() *preferences.Preferences {
	return bs.prefs
}

// SavePreferences saves the preferences to disk
func (bs *BrowserService) SavePreferences(prefs *preferences.Preferences) error {
	if err := prefs.Save(); err != nil {
		slog.With("error", err).Error("Failed to save preferences")
		return err
	}
	bs.prefs = prefs
	// Emit event to update tray menu
	bs.app.Event.Emit("preferences-changed", struct{}{})
	return nil
}

// OpenPreferences shows the preferences window
func (bs *BrowserService) OpenPreferences() {
	if bs.prefsWindow != nil {
		bs.prefsWindow.Show()
		bs.prefsWindow.Center()
		bs.prefsWindow.Focus()
	}
	bs.app.Event.Emit("open-preferences", struct{}{})
}

// ClosePreferences hides the preferences window
func (bs *BrowserService) ClosePreferences() {
	if bs.prefsWindow != nil {
		bs.prefsWindow.Hide()
	}
}

// Quit quits the application
func (bs *BrowserService) Quit() {
	bs.app.Quit()
}

// CheckForUpdates checks if a new version is available
func (bs *BrowserService) CheckForUpdates() (*updater.UpdateInfo, error) {
	return updater.CheckForUpdates()
}

// GetAppVersion returns the current application version
func (bs *BrowserService) GetAppVersion() string {
	return updater.GetCurrentVersion()
}
