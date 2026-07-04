# Installation Guide

This guide covers the currently supported installation methods for Browski.

> **Note:** Additional installation methods (Homebrew, Chocolatey, Snap, etc.) are planned for future releases. See [ROADMAP.md](../roadmap/ROADMAP.md) for details.

## macOS

### Download from GitHub Releases

1. Go to the [Releases page](https://github.com/icaliskanoglu/browski/releases)
2. Download the latest version for your Mac:
   - Apple Silicon (M1/M2/M3): Download the artifact containing `arm64`
   - Intel Macs: Download the artifact containing `amd64`
3. Extract the downloaded archive
4. Find `Browski.app` in the extracted files
5. Move `Browski.app` to `/Applications/`
6. First launch: Right-click the app → **Open** (due to macOS security for unsigned apps)

### Setting as Default Browser

1. Launch Browski (it will appear in the menu bar)
2. Open **System Settings** (or System Preferences on older macOS)
3. Navigate to **Desktop & Dock**
4. Scroll down to **Default web browser**
5. Select **Browski** from the dropdown
6. Test by clicking any http/https link

### Uninstalling

```bash
rm -rf /Applications/Browski.app
rm -rf ~/.browski  # Remove preferences and logs (optional)
```

## Windows

### Download from GitHub Releases

1. Go to the [Releases page](https://github.com/icaliskanoglu/browski/releases)
2. Download the Windows artifact (`browski-windows-amd64`)
3. Extract the downloaded archive
4. Run `browski.exe`

> **Note:** Windows installer (NSIS/MSI) and Chocolatey support are planned for future releases. See [ROADMAP.md](../roadmap/ROADMAP.md).

### Setting as Default Browser

1. Launch Browski (it will appear in the system tray)
2. Open **Settings** → **Apps** → **Default apps**
3. Click on **Web browser**
4. Select **Browski** from the list
5. Test by clicking any http/https link

### Uninstalling

1. Delete the `browski.exe` file
2. Remove preferences (optional):
   - Press `Win + R`
   - Type `%APPDATA%` and press Enter
   - Delete the `browski` folder if it exists

## Linux

### Download Binary

1. Go to the [Releases page](https://github.com/icaliskanoglu/browski/releases)
2. Download the Linux artifact (`browski-linux-amd64`)
3. Extract the archive and find the `browski` binary
4. Make it executable:
   ```bash
   chmod +x browski
   ```
5. Move to a directory in your PATH (optional):
   ```bash
   sudo mv browski /usr/local/bin/
   ```

> **Note:** Package formats (.deb, .rpm, .AppImage) and distribution via Snap Store, AUR, and Flatpak are planned for future releases. See [ROADMAP.md](../roadmap/ROADMAP.md).

### Setting as Default Browser

**GNOME:**
1. Open Settings → Default Applications
2. Set Web to Browski

**KDE:**
1. Open System Settings → Applications → Default Applications
2. Set Web Browser to Browski

**Command line:**
```bash
xdg-settings set default-web-browser browski.desktop
```

### Uninstalling

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/browski

# Remove preferences and logs (optional)
rm -rf ~/.config/browski
rm -rf ~/.browski
```

## Building from Source

### Prerequisites

- **Go**: 1.25 or later
- **Node.js**: 18 or later
- **npm**: Latest version
- **Wails CLI**: v3

Install Wails CLI:
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Verify installation:
```bash
wails3 doctor
```

### Build Steps

```bash
# Clone the repository
git clone https://github.com/icaliskanoglu/browski.git
cd browski

# Install frontend dependencies
cd frontend
npm install
cd ..

# Build for your platform
wails3 build

# Package (macOS only - creates .app bundle)
wails3 task darwin:package

# The binary will be in:
# - macOS: bin/browski.app
# - Linux: bin/browski
# - Windows: bin/browski.exe
```

### Platform-Specific Builds

**macOS:**
```bash
# Apple Silicon
wails3 build -platform darwin/arm64

# Intel
wails3 build -platform darwin/amd64

# Universal (both architectures)
wails3 build -platform darwin/universal

# Don't forget to package
wails3 task darwin:package
```

**Windows:**
```bash
wails3 build -platform windows/amd64
```

**Linux:**
```bash
wails3 build -platform linux/amd64
```

### Installing Built Binary

**macOS:**
```bash
cp -r bin/browski.app /Applications/
```

**Linux:**
```bash
sudo cp bin/browski /usr/local/bin/
```

**Windows:**
```powershell
# Copy to Program Files
Copy-Item bin\browski.exe "C:\Program Files\Browski\"
```

## Verifying Installation

After installation, verify Browski is working:

```bash
# macOS
/Applications/Browski.app/Contents/MacOS/browski --version

# Linux
browski --version

# Windows (in PowerShell)
& "C:\Program Files\Browski\browski.exe" --version
```

Test with a URL:
```bash
# macOS
open -a Browski "https://example.com"

# Linux
browski "https://example.com"

# Windows
Start-Process "browski://https://example.com"
```

## Troubleshooting

### macOS: "App is damaged and can't be opened"

This happens because the app isn't signed. Fix it:
```bash
xattr -cr /Applications/Browski.app
```

### macOS: URL scheme not working

1. Rebuild with proper packaging:
   ```bash
   wails3 task darwin:package
   ```
2. Reinstall to /Applications/
3. Re-select as default browser in System Settings

### Linux: Desktop file not found

Create it manually:
```bash
cat > ~/.local/share/applications/browski.desktop <<EOF
[Desktop Entry]
Name=Browski
Exec=/usr/local/bin/browski %u
Terminal=false
Type=Application
Icon=web-browser
Categories=Network;WebBrowser;
MimeType=x-scheme-handler/http;x-scheme-handler/https;text/html;
EOF

# Update database
update-desktop-database ~/.local/share/applications/
```

### Windows: Installer blocked by Windows Defender

This is normal for unsigned executables. Click "More info" → "Run anyway".

For a permanent solution, the developer needs to code sign the application.

### Can't find Browski in default browser list

**macOS:**
- Ensure app is in `/Applications/` (not Downloads or Desktop)
- Verify the app was packaged correctly with URL scheme support
- Restart your Mac
- Check logs: `tail -f ~/.browski/log/browski.log`

**Windows:**
- Launch Browski at least once before setting as default
- Check Windows Settings → Apps → Default apps

**Linux:**
- Browski needs a `.desktop` file to appear in the default browser list
- Create one manually (see Desktop file section in troubleshooting above)
- Update database: `update-desktop-database ~/.local/share/applications/`

## Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Search [GitHub Issues](https://github.com/icaliskanoglu/browski/issues)
3. Create a new issue with:
   - Your OS and version
   - Browski version
   - Installation method
   - Error messages or logs from `~/.browski/log/browski.log`

## Next Steps

After installation:

1. [Set up Browski as your default browser](#setting-as-default-browser)
2. Customize which browsers/profiles appear in [Preferences](../README.md#configuration)
3. Learn about usage in the [README](../README.md#usage)
4. Check [Troubleshooting Guide](troubleshooting.md) if you encounter issues

## Future Installation Methods

See [Roadmap](../planning/roadmap.md) for planned improvements including:
- Native installers (DMG, MSI, NSIS)
- Package managers (Homebrew, Chocolatey, Snap, Flatpak, AUR)
- Auto-update support
- Code signing for all platforms
