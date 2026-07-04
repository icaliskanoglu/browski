package preferences

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Preferences stores user preferences for the application
type Preferences struct {
	DefaultBrowser  string            `json:"defaultBrowser"`  // Browser name + profile name (e.g., "Google Chrome:Work")
	HiddenBrowsers  map[string]bool   `json:"hiddenBrowsers"`  // Map of browser+profile combinations to hide
	AlwaysUseDefault bool             `json:"alwaysUseDefault"` // If true, always open with default without showing UI
}

// GetPreferencesPath returns the path to the preferences file
func GetPreferencesPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".browski", "preferences.json"), nil
}

// Load reads preferences from disk
func Load() (*Preferences, error) {
	path, err := GetPreferencesPath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return default preferences
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Preferences{
			HiddenBrowsers: make(map[string]bool),
		}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var prefs Preferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, err
	}

	// Ensure HiddenBrowsers map is initialized
	if prefs.HiddenBrowsers == nil {
		prefs.HiddenBrowsers = make(map[string]bool)
	}

	return &prefs, nil
}

// Save writes preferences to disk
func (p *Preferences) Save() error {
	path, err := GetPreferencesPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetBrowserKey returns a unique key for a browser+profile combination
func GetBrowserKey(browserName, profileName string) string {
	if profileName == "" {
		return browserName
	}
	return browserName + ":" + profileName
}

// IsBrowserHidden checks if a browser+profile combination is hidden
func (p *Preferences) IsBrowserHidden(browserName, profileName string) bool {
	key := GetBrowserKey(browserName, profileName)
	return p.HiddenBrowsers[key]
}

// SetBrowserHidden sets whether a browser+profile combination is hidden
func (p *Preferences) SetBrowserHidden(browserName, profileName string, hidden bool) {
	key := GetBrowserKey(browserName, profileName)
	if hidden {
		p.HiddenBrowsers[key] = true
	} else {
		delete(p.HiddenBrowsers, key)
	}
}
