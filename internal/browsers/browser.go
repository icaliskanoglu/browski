package browsers

import (
	"errors"
	"log/slog"
	"os"
	"runtime"
)

func ListBrowsers() []Browser {
	return browsers
}

var browsers []Browser

func userConfigDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LocalAppData")
		if dir == "" {
			return "", errors.New("%LocalAppData% is not defined")
		}

	case "darwin", "ios":
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir += "/Library/Application Support"

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}
		dir += "/lib"

	default: // Unix
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
			}
			dir += "/.config"
		}
	}

	return dir, nil
}

func init() {
	dir, err := userConfigDir()
	if err != nil {
		slog.With(err).Error("Error getting user config dir")
		os.Exit(1)
	}
	browsers = append(browsers, chromiumBrowsers(dir)...)
	browsers = append(browsers, firefoxBrowsers()...)
	browsers = append(browsers, safariBrowsers()...)
}

func Open(request OpenRequest, url string) error {
	switch request.Type {
	case "chromium":
		return chromiumOpen(request, url)
	case "firefox":
		return firefoxOpen(request, url)
	case "safari":
		return safariOpen(request, url)
	default:
		return errors.New("browser type not found " + request.Type)
	}
}

type Profile struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
	Icon      string `json:"icon"`
}

type Browser struct {
	Type            string    `json:"type"`
	Name            string    `json:"name"`
	Profiles        []Profile `json:"profiles"`
	Icon            string    `json:"icon"`
	ShowMainEntry   bool      `json:"showMainEntry"` // Whether to show the main browser entry (without profile)
	executable      string
	configDirectory string
}

type OpenRequest struct {
	BrowserName string `json:"browser"`
	ProfileName string `json:"profile"`
	Type        string `json:"type"`
}

/*func listDesktopFiles() {
	directories := []string{"/usr/share/applications/", os.Getenv("HOME") + "/.local/share/applications/"}

	for _, dir := range directories {
		files, err := filepath.Glob(filepath.Join(dir, "*.desktop"))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, file := range files {
			execLine, icon, found := checkDesktopFile(file)
			if found {
				fmt.Printf("Found in: %s, Exec: %s, Icon: %s\n", file, execLine, icon)
			}
		}
	}
}
*/
/*func checkDesktopFile(filePath string) (string, string, bool) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", "", false
	}
	defer file.Close()

	var execLine, iconName string
	supportsHttps := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MimeType=") && strings.Contains(line, "x-scheme-handler/https") {
			supportsHttps = true
		}
		if supportsHttps {
			if strings.HasPrefix(line, "Exec=") {
				execLine = strings.TrimPrefix(line, "Exec=")
			}
			if strings.HasPrefix(line, "Icon=") {
				iconName = strings.TrimPrefix(line, "Icon=")
				themeNew, err := gtk.IconThemeNew()
				if err != nil {
					fmt.Println("Error:", err)
					return "", "", false
				}
				icon := themeNew.HasIcon(iconName)
				println(icon)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return "", "", false
	}

	return execLine, iconName, supportsHttps
}*/

/*func ListAppBySchema() {
	def := glib.SettingsSchemaSourceGetDefault()
	nonReolcatable, relocatable := def.ListSchemas(false)
	for _, schema := range nonReolcatable {
		fmt.Println(schema)
	}

	for _, schema := range relocatable {
		fmt.Println(schema)
	}
}*/
