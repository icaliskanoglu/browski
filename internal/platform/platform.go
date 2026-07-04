package platform

type Platform string

var (
	Linux   Platform = "linux"
	Windows Platform = "windows"
	Darwin  Platform = "darwin"

	LogFile      = "browski.log"
	LogDirectory = map[Platform]string{
		Linux:   "/var/log/browski",
		Windows: "C:\\ProgramData\\browski\\logs",
		Darwin:  "/var/log/browski",
	}
)
