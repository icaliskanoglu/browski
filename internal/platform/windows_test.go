//go:build windows
// +build windows

package platform

import (
	"testing"
)

func TestGetCurrentDefaultBrowser(t *testing.T) {
	browser, err := getCurrentDefaultBrowser()
	if err != nil {
		t.Fatal(err)
	}
	println(browser)
}

func TestGetProgId(t *testing.T) {
	progID, err := getProgID("BrowskiHtml")
	if err != nil {
		// Skip test if BrowskiHtml is not registered (expected on CI or fresh systems)
		t.Skipf("BrowskiHtml not registered: %v", err)
	}
	println(progID)
}

func TestSetDefaultBrowser(t *testing.T) {
	err := setBrowskiDefaultBrowser("BrowskiHtml")
	if err != nil {
		t.Fatal(err)
	}
}
