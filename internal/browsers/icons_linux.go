// +build linux

package browsers

// getAppIconOrDefault returns the default URL on Linux (icon extraction not implemented)
func getAppIconOrDefault(appPath, defaultURL string) string {
	return defaultURL
}
