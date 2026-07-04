package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"browski/internal/browsers"
	"browski/internal/updater"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

// getSocketPath returns the IPC socket path
func getSocketPath() string {
	return filepath.Join(os.TempDir(), "browski.sock")
}

// startIPCServer starts listening for URL messages from other instances
func startIPCServer(browserService *BrowserService) {
	sockPath := getSocketPath()

	// Remove any existing socket file
	os.Remove(sockPath)

	// Create Unix domain socket
	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		log.Printf("Failed to create IPC socket: %v", err)
		return
	}

	go func() {
		defer listener.Close()
		defer os.Remove(sockPath)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("IPC accept error: %v", err)
				return
			}

			go handleIPCConnection(conn, browserService)
		}
	}()
}

// handleIPCConnection handles an incoming IPC connection
func handleIPCConnection(conn net.Conn, browserService *BrowserService) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			log.Printf("Received URL from another instance: %s", url)
			browserService.SetURL(url)
		}
	}
}

// sendURLToExistingInstance sends a URL to an already-running instance
func sendURLToExistingInstance(urlArg string) bool {
	if urlArg == "" {
		return false
	}

	sockPath := getSocketPath()

	// Try to connect to the existing instance's socket
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return false
	}
	defer conn.Close()

	// Send the URL
	_, err = fmt.Fprintf(conn, "%s\n", urlArg)
	if err != nil {
		log.Printf("Failed to send URL to existing instance: %v", err)
		return false
	}

	log.Printf("Sent URL to existing instance: %s", urlArg)
	return true
}

func init() {
	// Register events
	application.RegisterEvent[string]("url-changed")
	application.RegisterEvent[struct{}]("open-preferences")
	application.RegisterEvent[struct{}]("preferences-changed")
}

func main() {
	// Get URL from command line arguments if provided
	var urlArg string
	if len(os.Args) > 1 {
		urlArg = os.Args[1]
	}

	// Check if another instance is already running
	if checkExistingInstance(urlArg) {
		// Another instance is running, exit silently
		return
	}

	// Acquire single instance lock
	lockFile, err := acquireLock()
	if err != nil {
		log.Printf("Failed to acquire lock: %v", err)
		return
	}
	defer releaseLock(lockFile)

	// Create application
	app := application.New(application.Options{
		Name:        "Browski",
		Description: "Browser Profile Picker",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
			ActivationPolicy: application.ActivationPolicyAccessory, // Hide from Dock
		},
	})

	// Create browser service
	browserService := NewBrowserService(app)

	// Register service
	app.RegisterService(application.NewService(browserService))

	// Start IPC server for handling URLs from other instances
	startIPCServer(browserService)

	// Create system tray
	systray := app.SystemTray.New()

	// Use template icon on macOS for automatic dark mode support
	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(icon)
	} else {
		systray.SetIcon(icon)
	}

	// Build tray menu with browser list (like Choosy)
	buildTrayMenu := func() *application.Menu {
		menu := app.Menu.New()

		// Add browser profiles to menu
		allBrowsers := browserService.ListBrowsers()
		for _, browser := range allBrowsers {
			// Add main browser only if ShowMainEntry is true
			if browser.ShowMainEntry {
				browserName := browser.Name
				browserType := browser.Type
				// Capture variables for closure
				bName := browserName
				bType := browserType
				menu.Add(browserName).OnClick(func(ctx *application.Context) {
					// Open browser with empty URL (opens default page)
					browserService.Open(browsers.OpenRequest{
						BrowserName: bName,
						Type:        bType,
						ProfileName: "",
					})
				})
			}

			// Add profiles for this browser
			for _, profile := range browser.Profiles {
				profileLabel := "  " + profile.Name + " (" + browser.Name + ")"
				// Capture variables for closure
				bName := browser.Name
				bType := browser.Type
				pName := profile.Name
				menu.Add(profileLabel).OnClick(func(ctx *application.Context) {
					// Open browser with profile
					browserService.Open(browsers.OpenRequest{
						BrowserName: bName,
						Type:        bType,
						ProfileName: pName,
					})
				})
			}
		}

		menu.AddSeparator()
		menu.Add("Preferences...").OnClick(func(ctx *application.Context) {
			browserService.OpenPreferences()
		})
		menu.AddSeparator()
		menu.Add("Quit").OnClick(func(ctx *application.Context) {
			app.Quit()
		})

		return menu
	}

	systray.SetMenu(buildTrayMenu())

	// Listen for preferences changes to rebuild menu
	app.Event.On("preferences-changed", func(event *application.CustomEvent) {
		systray.SetMenu(buildTrayMenu())
	})

	// Create browser chooser window (hidden by default)
	chooserWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Browski",
		Width:            420,
		Height:           180,
		Hidden:           true,
		Frameless:        true,
		AlwaysOnTop:      true,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0), // Transparent background for rounded corners
		URL:              "/",
	})

	// Create preferences window (hidden by default, native window with controls)
	prefsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Browski",
		Width:            500,
		Height:           600,
		Hidden:           true,
		Frameless:        false, // Native window with controls
		BackgroundColour: application.NewRGBA(30, 30, 30, 255),
		URL:              "/preferences.html",
	})

	// Intercept window close to hide instead of destroy
	prefsWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		e.Cancel()         // Prevent the window from being destroyed
		prefsWindow.Hide() // Hide it instead
	})

	// Set window references in service
	browserService.SetWindows(chooserWindow, prefsWindow)

	// Handle URLs from system (when app is launched via http/https links)
	// This must be registered AFTER windows are created
	app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) {
		url := e.Context().URL()
		if url != "" {
			log.Printf("Received URL from system: %s", url)
			browserService.SetURL(url)
		}
	})

	// Disable resizing for chooser window
	chooserWindow.SetResizable(false)

	// Don't attach window to tray - we want menu behavior like Choosy, not window toggling

	// Set URL if provided
	if urlArg != "" {
		browserService.SetURL(urlArg)
	}

	// Don't show the window here - let the frontend show it after content is loaded
	// This eliminates the loading delay flash

	// Check for updates in background (non-blocking)
	updater.CheckInBackground()

	runErr := app.Run()
	if runErr != nil {
		log.Fatal(runErr)
	}
}
