// +build windows

package browsers

// getAppIconOrDefault returns the default URL on Windows (icon extraction not implemented)
func getAppIconOrDefault(appPath, defaultURL string) string {
	return defaultURL
}
