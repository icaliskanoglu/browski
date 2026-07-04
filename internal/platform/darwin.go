//go:build darwin
// +build darwin

package platform

import (
	"browski/internal/platform/darwin"
)

func SetDefaultBrowser() {

	darwin.SetDefaultBrowser("com.ihsanc.browski")
}
