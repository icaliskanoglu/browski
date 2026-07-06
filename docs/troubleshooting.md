# Troubleshooting Guide

Common issues and solutions for Browski across all platforms.

---

## Quick Diagnostics

### Check Logs

Browski logs to platform-specific locations:

**macOS**:
```bash
tail -f ~/.browski/log/browski.log
```

**Linux**:
```bash
tail -f ~/.browski/log/browski.log
# or
tail -f ~/.config/browski/browski.log
```

**Windows**:
```powershell
Get-Content $env:APPDATA\browski\browski.log -Tail 20 -Wait
```

### Check Version

```bash
# macOS
/Applications/Browski.app/Contents/MacOS/browski --version

# Linux
browski --version

# Windows
& "C:\Program Files\Browski\browski.exe" --version
```

---

## Installation Issues

### macOS: "App is damaged and can't be opened"

**Cause**: macOS Gatekeeper blocking unsigned apps

**Solution**:
```bash
xattr -cr /Applications/Browski.app
```

Then right-click Browski.app → Open → Open anyway

### Windows: Installer blocked by Windows Defender

**Cause**: Unsigned executable triggering SmartScreen

**Solution**:
1. Click "More info"
2. Click "Run anyway"

**Note**: This warning will disappear once the app is code-signed

### Linux: Permission denied

**Cause**: Binary not executable

**Solution**:
```bash
chmod +x /path/to/browski
```

---

## Default Browser Setup

### Can't Find Browski in Default Browser List

#### macOS

**Symptoms**: Browski doesn't appear in System Settings → Desktop & Dock → Default web browser

**Solutions**:
1. **Ensure proper location**:
   ```bash
   # App must be in /Applications/, not Downloads
   sudo mv ~/Downloads/Browski.app /Applications/
   ```

2. **Verify bundle packaging**:
   ```bash
   # Check Info.plist has URL types
   /usr/libexec/PlistBuddy -c "Print :CFBundleURLTypes" /Applications/Browski.app/Contents/Info.plist
   ```

3. **Rebuild with packaging** (if building from source):
   ```bash
   wails3 task darwin:package
   cp -r bin/Browski.app /Applications/
   ```

4. **Restart your Mac**

5. **Check logs**:
   ```bash
   tail -f ~/.browski/log/browski.log
   ```

#### Windows

**Solutions**:
1. Launch Browski at least once before setting as default
2. Open Settings → Apps → Default apps
3. If not listed, reinstall with proper registry keys

**Check registry** (if you're technical):
```powershell
reg query "HKLM\SOFTWARE\RegisteredApplications" /v Browski
```

#### Linux

**Solutions**:
1. **Create .desktop file manually**:
   ```bash
   cat > ~/.local/share/applications/browski.desktop <<EOF
   [Desktop Entry]
   Name=Browski
   Comment=Browser profile picker
   Exec=/usr/local/bin/browski %u
   Terminal=false
   Type=Application
   Icon=web-browser
   Categories=Network;WebBrowser;
   MimeType=x-scheme-handler/http;x-scheme-handler/https;text/html;
   EOF
   ```

2. **Update desktop database**:
   ```bash
   update-desktop-database ~/.local/share/applications/
   ```

3. **Set as default**:
   ```bash
   xdg-settings set default-web-browser browski.desktop
   ```

4. **Verify**:
   ```bash
   xdg-settings get default-web-browser
   ```

---

## URL Handling Issues

### Blank Page When Clicking Links (macOS)

**Symptoms**: Browski opens but shows blank window when clicking http/https links

**Solutions**:
1. **Ensure using packaged .app bundle**, not raw binary:
   ```bash
   # Wrong: Running binary directly
   ./bin/browski https://example.com

   # Right: Using packaged app
   /Applications/Browski.app/Contents/MacOS/browski https://example.com
   ```

2. **Reinstall to /Applications/**:
   ```bash
   rm -rf /Applications/Browski.app
   cp -r path/to/Browski.app /Applications/
   ```

3. **Re-select as default browser** in System Settings

4. **Check logs for errors**:
   ```bash
   tail -f ~/.browski/log/browski.log
   ```

### macOS: URL Scheme Not Working

**Symptoms**: Clicking links doesn't launch Browski

**Solutions**:
1. **Rebuild with proper packaging**:
   ```bash
   wails3 task darwin:package
   ```

2. **Reinstall to /Applications/**

3. **Reset default browser**:
   - System Settings → Desktop & Dock → Default web browser
   - Select Browski again

4. **Clear launch services cache** (if still not working):
   ```bash
   /System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill -r -domain local -domain system -domain user
   ```

### URLs Don't Open in Selected Browser

**Symptoms**: Browski opens but URL doesn't launch in chosen browser

**Solutions**:
1. **Check that Browski is set as default browser**

2. **Test from command line**:
   ```bash
   # macOS
   /Applications/Browski.app/Contents/MacOS/browski "https://example.com"

   # Linux
   browski "https://example.com"

   # Windows
   & "C:\Program Files\Browski\browski.exe" "https://example.com"
   ```

3. **Check logs for error messages**

4. **Verify browser paths** - Browser might have been uninstalled or moved

---

## Browser Detection Issues

### Browser Not Detected

**Symptoms**: Installed browser doesn't appear in Browski's list

**Solutions**:
1. **Ensure browser is properly installed**:
   - macOS: Check `/Applications/`
   - Windows: Check `C:\Program Files\` and `C:\Program Files (x86)\`
   - Linux: Check `which chrome` / `which firefox` etc.

2. **Check if hidden in Preferences**:
   - Right-click Browski tray icon → Preferences
   - Verify browser isn't toggled off

3. **Check supported browsers**:
   - Chrome, Edge, Brave, Arc, Vivaldi, Opera (Chromium-based)
   - Firefox
   - Safari (macOS only)

4. **Report issue** with:
   - Browser name and version
   - Install location
   - Operating system

### Profiles Not Showing

**Symptoms**: Browser detected but profiles missing

**Solutions**:
1. **Ensure profiles exist** - Launch browser and check profiles

2. **Check profile directories**:
   - **macOS Chrome**: `~/Library/Application Support/Google/Chrome/`
   - **Windows Chrome**: `%LOCALAPPDATA%\Google\Chrome\User Data\`
   - **Linux Chrome**: `~/.config/google-chrome/`
   - **Firefox**: Look for `profiles.ini`

3. **Restart Browski** to re-scan

---

## Single-Instance Issues

### Multiple Windows Appear

**Symptoms**: Clicking multiple links spawns multiple Browski windows

**Expected**: Only one window, subsequent URLs forwarded to existing instance

**Solutions**:
1. **Check lock file permissions**:
   ```bash
   # macOS/Linux
   ls -la /tmp/browski.lock

   # Windows
   dir %TEMP%\browski.lock
   ```

2. **Check IPC socket** (macOS/Linux):
   ```bash
   ls -la /tmp/browski.sock
   ```

3. **Clear lock files and restart**:
   ```bash
   # macOS/Linux
   rm /tmp/browski.lock /tmp/browski.sock

   # Windows
   del %TEMP%\browski.lock
   ```

4. **Report issue** if persists on your platform

---

## Performance Issues

### Slow Browser Detection

**Symptoms**: Browski takes long time to load browser list

**Cause**: Scanning large profile directories

**Solutions**:
1. **Hide unused browsers/profiles** in Preferences
2. **Clear browser cache** (reduces profile scanning time)
3. Consider reporting for optimization

### High CPU Usage

**Symptoms**: Browski using significant CPU when idle

**Solutions**:
1. **Check logs** for errors causing retries
2. **Restart Browski**
3. **Report issue** with:
   - OS and version
   - Number of browsers/profiles detected
   - Logs from `~/.browski/log/browski.log`

---

## System Tray Issues

### Tray Icon Not Appearing

**macOS**:
- Check menu bar (upper right)
- Browski uses `ActivationPolicyAccessory` so it won't appear in Dock
- Press Cmd+Space, type "Browski", press Enter to relaunch

**Windows**:
- Check system tray (lower right, may be hidden)
- Click up arrow to show hidden icons

**Linux**:
- System tray support varies by desktop environment
- GNOME requires extension: `sudo apt install gnome-shell-extension-appindicator`
- KDE: Should work out of box

---

## Build Issues (Developers)

### Wails Build Fails

**Check system requirements**:
```bash
wails3 doctor
```

**Common issues**:
- **macOS**: Install Xcode Command Line Tools
- **Windows**: Install WebView2 runtime
- **Linux**: Install GTK4 and WebKitGTK:
  ```bash
  sudo apt install libgtk-4-dev libwebkitgtk-6.0-dev
  ```

### Frontend Build Fails

**Solutions**:
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Binary Name Wrong (browski)

**Cause**: Old Taskfile configuration

**Solution**: Check `Taskfile.yml` has `APP_NAME: "browski"` (not "browski")

---

## Getting Help

If your issue isn't listed here:

1. **Check logs** (see Quick Diagnostics above)
2. **Search [GitHub Issues](https://github.com/icaliskanoglu/browski/issues)**
3. **Create new issue** with:
   - Operating system and version
   - Browski version (`browski --version`)
   - Installation method (GitHub Release, Homebrew, built from source)
   - Steps to reproduce
   - Relevant logs
   - Screenshots if applicable

---

## See Also

- [Installation Guide](installation.md)
- [Contributing Guide](contributing.md)
- [README](../README.md)
