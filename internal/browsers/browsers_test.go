package browsers

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func prettyPrint(i interface{}) string {
	s, _ := yaml.Marshal(i)
	return string(s)
}

func TestListBrowsers(t *testing.T) {
	// t.Skip("Skipping test as it requires a chrome installation")
	browsers := ListBrowsers()

	println(prettyPrint(browsers))
}

func TestChromeOpen(t *testing.T) {
	// t.Skip("Skipping test as it requires a chrome installation")
	chromiumOpen(OpenRequest{
		BrowserName: "Edge",
		Type:        "chromium",
	}, "https://www.google.com")
}

/*func TestListDesktopFiles(t *testing.T) {
	// t.Skip("Skipping test as it requires a chrome installation")
	ListAppBySchema()
}
*/
