//go:build windows
// +build windows

package platform

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

func SetDefaultBrowser() {

	println("SetDefaultBrowser")
}

// getCurrentDefaultBrowser retrieves the ProgID of the current default browser.
func getCurrentDefaultBrowser() (string, error) {
	httpKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`, registry.READ)
	if err != nil {
		return "", err
	}
	defer httpKey.Close()

	progID, _, err := httpKey.GetStringValue("ProgId")
	if err != nil {
		return "", err
	}

	return progID, nil
}

/*
func registerBrowski() error {
	yourBrowskiProgID := "YourBrowskiProgID"
	browserIconPath := "C:\\Path\\To\\Your\\Browser\\Icon.ico"

	// Create registry entries for your ProgID
	progIDKey, err := registry.CreateKey(registry.CLASSES_ROOT, yourBrowskiProgID)
	if err != nil {
		return err
	}
	defer progIDKey.Close()

	// Set default values for ProgID
	if err := progIDKey.SetStringValue("", "BrowskiURLHandler"); err != nil {
		return err
	}

	// Create registry entries for URL handling
	urlHandlerKey, err := registry.CreateKey(progIDKey, "BrowskiURLHandler")
	if err != nil {
		return err
	}
	defer urlHandlerKey.Close()

	if err := urlHandlerKey.SetStringValue("", "URL: Browski Protocol"); err != nil {
		return err
	}

	if err := urlHandlerKey.SetStringValue("URL Protocol", ""); err != nil {
		return err
	}

	// Create registry entries for icon
	iconKey, err := registry.CreateKey(urlHandlerKey, "DefaultIcon")
	if err != nil {
		return err
	}
	defer iconKey.Close()

	if err := iconKey.SetStringValue("", browserIconPath); err != nil {
		return err
	}

	// Create registry entries for command
	commandKey, err := registry.CreateKey(urlHandlerKey, "shell\\open\\command")
	if err != nil {
		return err
	}
	defer commandKey.Close()

	// Specify the command without setting it as the default browser
	if err := commandKey.SetStringValue("", fmt.Sprintf(`"C:\\Path\\To\\Your\\Browser\\Browski.exe" "%%1"`)); err != nil {
		return err
	}

	return nil
}*/

// setBrowskiDefaultBrowser sets "browski" as the default browser by updating the ProgID in the registry.
func setBrowskiDefaultBrowser(progId string) error {

	// Set ProgID for HTTP protocol
	httpKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer httpKey.Close()
	if err := httpKey.SetStringValue("ProgId", progId); err != nil {
		return err
	}

	// Set ProgID for HTTPS protocol
	httpsKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\Shell\Associations\UrlAssociations\https\UserChoice`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer httpsKey.Close()
	if err := httpsKey.SetStringValue("ProgId", progId); err != nil {
		return err
	}

	return nil
}

// getProgID retrieves the ProgID of a program based on its executable path.
func getProgID(name string) (string, error) {

	// Open the registry key for the file extension
	extKey, err := registry.OpenKey(registry.CLASSES_ROOT, name, registry.READ)
	if err != nil {
		return "", err
	}
	defer extKey.Close()

	// Read the ProgID associated with the file extension
	progID, _, err := extKey.GetStringValue("")
	if err != nil {
		return "", err
	}

	// Open the registry key for the ProgID
	progIDKey, err := registry.OpenKey(registry.CLASSES_ROOT, progID, registry.READ)
	if err != nil {
		return "", err
	}
	defer progIDKey.Close()

	// Read the friendly name of the ProgID (optional)
	friendlyName, _, err := progIDKey.GetStringValue("")
	if err != nil {
		// Ignore errors when reading friendly name
		friendlyName = ""
	}

	return fmt.Sprintf("%s (%s)", progID, friendlyName), nil
}
