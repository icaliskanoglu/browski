# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important Rules

- **NEVER add Claude as co-author in git commits**. Do not include "Co-Authored-By: Claude <noreply@anthropic.com>" or similar attribution in commit messages.

## Project Overview

Browski is a cross-platform browser profile picker built with Wails v3 (Go + React). It can be set as your default browser and displays as an always-on-top, frameless window that shows available browsers and their profiles, allowing users to quickly select which browser/profile to open a URL with. The app runs as a system tray application and can handle http/https URLs from anywhere in the system.

## Technology Stack

- **Backend**: Go 1.25+ with Wails v3.0.0-alpha2.111+
- **Frontend**: React 18 with Vite
- **Platforms**: macOS (Darwin), Linux, Windows

## Development Commands

### System Requirements Check

Before starting development, verify your system has all required dependencies:

```bash
wails3 doctor
```

### Running the Application

```bash
# Live development mode with hot reload
wails3 dev

# Run with specific build flags
wails3 dev -tags dev
```

The Vite dev server runs at http://localhost:34115 where you can call Go methods from browser devtools.

### Building

```bash
# Build production binary for current platform
wails3 build

# Package into .app bundle (macOS) - required for URL scheme handling
wails3 task darwin:package

# Build for specific platform
wails3 build -platform darwin/universal  # macOS Intel + Apple Silicon
wails3 build -platform darwin/arm64      # macOS Apple Silicon only
wails3 build -platform darwin/amd64      # macOS Intel only
wails3 build -platform linux/amd64       # Linux
wails3 build -platform windows/amd64     # Windows

# Build without packaging (faster for testing)
wails3 build -noPackage

# Build with debug mode
wails3 build -debug
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run Vite dev server (standalone)
npm run dev

# Build frontend only
npm run build
```

### Testing

```bash
# Run Go tests
go test ./...

# Run specific package tests
go test ./internal/browsers

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Updating Dependencies

```bash
# Update Wails CLI to latest version
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# Update Go dependencies
go get -u ./...
go mod tidy

# Regenerate Wails bindings after Go updates
wails3 generate bindings

# Update frontend dependencies
cd frontend
npm update
npm audit fix  # Fix security vulnerabilities

# Check for outdated packages
npm outdated
```

### Common Development Tasks

```bash
# Generate TypeScript bindings for Go methods
wails3 generate bindings

# Regenerate build assets (after modifying build/config.yml)
wails3 task common:update:build-assets

# Clean build artifacts
rm -rf bin frontend/dist

# View application logs (macOS)
tail -f ~/.browski/log/browski.log

# Install to /Applications for testing (macOS)
cp -r bin/browski.app /Applications/
```

## Architecture

### Go Backend Structure

- **`main.go`**: Application entry point. Configures Wails with platform-specific options (transparent/translucent windows, frameless mode, always-on-top). Accepts URLs via:
  - Command line arguments (`browski https://example.com`)
  - URL scheme handlers (http/https via `ApplicationLaunchedWithUrl` event)
  - IPC from other instances (single instance enforcement)
  - Logs to `~/.browski/log/browski.log`

- **`browserservice.go`**: BrowserService struct that manages browser operations. Exposes methods to frontend:
  - `ListBrowsers()` - Returns detected browsers and profiles
  - `Open(OpenRequest)` - Opens URL in selected browser/profile
  - `GetURL()` - Returns current URL to be opened
  - `SetURL(url)` - Sets URL and emits url-changed event
  - `Resize(width, height)` - Dynamically resizes window
  - `ShowWindow()` / `HideWindow()` - Controls window visibility
  - `GetPreferences()` / `SavePreferences()` - Manages user preferences
  - `Quit()` - Closes application

- **`internal/browsers/`**: Browser detection and launching logic
  - `browser.go` - Core types (Browser, Profile, OpenRequest) and platform-agnostic browser listing
  - `chromium.go` - Detects Chromium-based browsers (Chrome, Edge, Brave, etc.) and their profiles
  - `firefox.go` - Firefox browser detection and profile handling
  - Platform-specific browser launching logic

- **`internal/platform/`**: Platform-specific integrations
  - `platform.go` - Platform constants and log directories
  - `darwin.go`, `linux.go`, `windows.go` - OS-specific implementations
  - `darwin/` - macOS-specific code including default browser detection using Objective-C (CGO)

### Frontend Structure

- **`frontend/src/components/Selector.jsx`**: Main browser selector component that:
  - Lists all browsers and profiles as clickable icons
  - Handles window resizing based on content
  - Quits on Escape key or mouse leave
  - Shows profile name on hover
  - Listens for url-changed events to show window when URLs arrive
  - Supports keyboard shortcuts (1-9 to select browser)

- **`frontend/src/preferences.jsx`**: Preferences management interface for hiding browsers/profiles

- **`frontend/bindings/`**: Auto-generated Wails bindings for calling Go methods from JavaScript

### Key Behaviors

1. **Browser Detection**: On startup, the app scans platform-specific config directories to discover installed browsers and their profiles (Chromium reads from Preferences JSON, Firefox reads from profiles.ini/installs.ini)

2. **URL Handling**: The app registers as a handler for http/https protocols via:
   - `CFBundleURLTypes` in Info.plist (macOS)
   - Desktop file associations (Linux)
   - Registry entries (Windows)
   - URLs are handled via `ApplicationLaunchedWithUrl` event (registered AFTER windows are created to avoid timing issues)

3. **Single Instance**:
   - Uses file locking to ensure only one instance runs
   - Platform-specific implementations (flock on Unix, exclusive file access on Windows)
   - IPC via Unix domain sockets to send URLs to existing instance

4. **System Tray**:
   - Runs as menu bar/system tray application
   - Hides from Dock on macOS (ActivationPolicyAccessory)
   - Tray menu shows all browsers and profiles for quick access

5. **Window Lifecycle**:
   - Chooser window is frameless, transparent, and always on top
   - Hidden by default, shown when URL is received
   - Dynamically resizes based on content after images load
   - Closes immediately when mouse leaves the window or Escape is pressed

6. **Profile Icons**: Each browser profile icon is extracted from the browser's config directory and displayed in the UI

## File Associations

The app registers itself as a viewer for:
- `.html` files (as "BrowskiHtml")
- `http://` and `https://` URL schemes

## Configuration

- **`build/config.yml`**: Wails v3 project configuration including:
  - App metadata (name, version, bundle ID)
  - Protocol handlers (http, https)
  - File associations (html)
  - Dev mode settings
  - After modifying, run `wails3 task common:update:build-assets` to regenerate platform assets

- **`Taskfile.yml`**: Build task configuration
  - `APP_NAME`: "browski" (binary name)
  - Platform-specific build tasks

- **`frontend/vite.config.js`**: Vite bundler configuration

- **`go.mod`**: Go dependencies (currently using Wails v3.0.0-alpha2.111)

## Wails-Specific Details

### Wails Runtime Context

The `App` struct in `app.go` receives a context during startup that provides access to Wails runtime methods:
- `application.Get().Window.SetSize()` - Resize the window
- `application.Get().Quit()` - Close the application
- Other runtime methods available in the `github.com/wailsapp/wails/v3/pkg/application` package

### Bindings

Go methods exposed to the frontend must be:
- Public (capitalized method names)
- Part of the bound struct (passed to `Bind` in `main.go`)
- Return values that are JSON-serializable

After adding/modifying Go methods, regenerate TypeScript bindings:
```bash
wails3 generate bindings
```

### Platform-Specific Behavior

The application uses platform-specific configurations in `main.go`:
- **macOS**: Translucent window with hidden title bar for frameless look
- **Linux**: Translucent window support
- **Windows**: Mica backdrop effect for modern Windows 11 appearance

### Debugging

- Use `wails3 dev` for hot reload during development
- Browser devtools available in the running application window
- Go logs output to `~/.browski/log/browski.log`
- Use `slog` for structured logging in Go code
- For URL scheme debugging: Check logs for "Received URL from system" messages
- Event timing: `ApplicationLaunchedWithUrl` handler must be registered AFTER `SetWindows()` is called

### Setting as Default Browser (macOS)

1. Build and package the app: `wails3 task darwin:package`
2. Install to Applications: `cp -r bin/browski.app /Applications/`
3. Open **System Settings** → **Desktop & Dock**
4. Scroll to **Default web browser** dropdown
5. Select **Browski**
6. Test by clicking any http/https link

### Known Issues

- The app uses a frameless window, which means standard window controls (minimize, maximize, close) must be implemented in the frontend if needed
- Window resizing is handled programmatically through `Resize()` method
- URL scheme handling requires proper .app bundle packaging (not just binary)
- `ApplicationLaunchedWithUrl` event must be registered AFTER windows are created to avoid blank page issue
