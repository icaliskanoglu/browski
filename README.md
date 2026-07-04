# Browski

A cross-platform browser profile picker built with Wails v3. Browski allows you to choose which browser and profile to open URLs with, making it perfect for users who manage multiple browser profiles or use different browsers for different tasks.

## Features

- **Default Browser Support**: Set Browski as your system's default browser
- **Multi-Browser Support**: Automatically detects Chrome, Firefox, Edge, Brave, and other popular browsers
- **Profile Selection**: Choose specific browser profiles for each URL
- **System Tray Integration**: Quick access from your menu bar/system tray
- **Keyboard Shortcuts**: Press 1-9 to quickly select browsers
- **Single Instance**: Intelligently handles multiple URL launches
- **Cross-Platform**: Works on macOS, Linux, and Windows

## Installation

### macOS

1. Download the latest release from [Releases](https://github.com/yourusername/browski/releases)
2. Unzip and move `Browski.app` to `/Applications/`
3. Launch Browski
4. To set as default browser:
   - Open **System Settings** → **Desktop & Dock**
   - Scroll to **Default web browser**
   - Select **Browski**

### Linux

1. Download the latest `.deb` or `.rpm` package
2. Install using your package manager
3. Launch Browski from your applications menu

### Windows

1. Download the latest installer
2. Run the installer
3. Launch Browski from the Start menu

## Usage

### As Default Browser

Once set as your default browser, Browski will automatically appear whenever you click a link:

1. Click any http/https link from any application
2. Browski appears with your available browsers
3. Select your preferred browser/profile
4. The URL opens in your chosen browser

### From System Tray

Click the Browski icon in your system tray to:
- Quickly launch any browser/profile
- Open preferences to hide unwanted browsers
- Quit the application

### Keyboard Shortcuts

When the browser picker appears:
- Press `1-9` to select the corresponding browser
- Press `ESC` to cancel and close the picker

## Building from Source

### Prerequisites

- Go 1.25+
- Node.js 18+
- Wails CLI v3: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

### Development

```bash
# Clone the repository
git clone https://github.com/yourusername/browski.git
cd browski

# Run in development mode
wails3 dev

# Build for production
wails3 build

# Package for macOS (creates .app bundle)
wails3 task darwin:package

# Install to /Applications for testing
cp -r bin/browski.app /Applications/
```

### Project Structure

```
browski/
├── main.go                 # Application entry point
├── browserservice.go       # Browser management service
├── internal/
│   ├── browsers/          # Browser detection logic
│   ├── platform/          # Platform-specific code
│   ├── preferences/       # User preferences
│   └── updater/          # Update checker
├── frontend/
│   └── src/
│       ├── components/    # React components
│       └── preferences.jsx
└── build/
    ├── config.yml         # Wails configuration
    └── darwin/           # macOS-specific assets
```

## Configuration

### Hiding Browsers

1. Right-click the Browski tray icon
2. Select "Preferences..."
3. Toggle browsers or profiles you want to hide
4. Click "Save"

### Preferences File

Preferences are stored in:
- **macOS**: `~/.browski/preferences.json`
- **Linux**: `~/.config/browski/preferences.json`
- **Windows**: `%APPDATA%\browski\preferences.json`

## How It Works

Browski registers itself as a handler for http/https URL schemes. When a URL is opened:

1. The operating system launches Browski with the URL
2. Browski scans for installed browsers and their profiles
3. A frameless window appears showing available options
4. User selects a browser/profile
5. The URL opens in the selected browser
6. Browski hides until the next URL

## Troubleshooting

For detailed troubleshooting, see the [Troubleshooting Guide](docs/troubleshooting.md).

Common issues:
- **Blank page on macOS**: Use packaged `.app` bundle, reinstall to `/Applications/`
- **Browser not detected**: Ensure browser is installed, check if hidden in Preferences
- **URL not opening**: Verify Browski is set as default browser, check logs

Full diagnostics and solutions in [docs/troubleshooting.md](docs/troubleshooting.md).

## Development

See [CLAUDE.md](CLAUDE.md) for detailed development documentation.

## Documentation

- [Installation Guide](docs/installation.md) - How to install on all platforms
- [Troubleshooting Guide](docs/troubleshooting.md) - Common issues and solutions
- [Contributing Guide](docs/contributing.md) - How to contribute to Browski
- [Roadmap](planning/roadmap.md) - Planned features and improvements
- [Distribution Guide](planning/distribution.md) - Packaging and release process

## License

Browski is open source software licensed under the [MIT License](LICENSE).

## Credits

Built with [Wails v3](https://wails.io)
