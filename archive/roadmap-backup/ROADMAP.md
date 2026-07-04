# Browski Roadmap

This document outlines planned features and improvements for Browski.

## Current Status

✅ **Completed (v1.0)**
- ✅ Cross-platform browser detection (Chrome, Firefox, Edge, Brave, Arc, etc.)
- ✅ Profile management and selection
- ✅ System tray integration (macOS, Windows, Linux)
- ✅ Default browser support via URL scheme handlers (http/https)
- ✅ Single instance with IPC communication
- ✅ Keyboard shortcuts (1-9 for selection, ESC to cancel)
- ✅ Preferences system for hiding browsers/profiles
- ✅ Basic CI/CD for building macOS, Windows, and Linux binaries
- ✅ GitHub Releases distribution

## Short Term (Next 1-2 Releases)

### Distribution & Packaging

**Priority: High**

- [ ] **macOS DMG Installer** - Create proper `.dmg` disk images for easier installation
  - [ ] Implement `darwin:package:dmg` task
  - [ ] Include license agreement
  - [ ] Custom background and app positioning
  - [ ] Update CI/CD workflow

- [ ] **Windows NSIS Installer** - Professional installer for Windows
  - [ ] Implement `windows:package:nsis` task
  - [ ] Add to Start Menu
  - [ ] Desktop shortcut (optional)
  - [ ] Uninstaller
  - [ ] Update CI/CD workflow

- [ ] **Linux Packages**
  - [ ] `.deb` packages (Debian/Ubuntu) - Implement `linux:package:deb` task
  - [ ] `.rpm` packages (Fedora/RHEL) - Implement `linux:package:rpm` task
  - [ ] `.AppImage` (Universal Linux) - Implement `linux:package:appimage` task
  - [ ] Update CI/CD workflow

### Package Manager Support

**Priority: Medium**

- [ ] **Homebrew (macOS)**
  - [ ] Create homebrew-browski tap repository
  - [ ] Automate formula updates on release
  - [ ] Test installation on different macOS versions

- [ ] **Chocolatey (Windows)**
  - [ ] Create `browski.nuspec` configuration
  - [ ] Submit to Chocolatey repository
  - [ ] Automate package updates

- [ ] **Snap Store (Linux)**
  - [ ] Create `snapcraft.yaml` configuration
  - [ ] Register on Snap Store
  - [ ] Set up confinement rules
  - [ ] Automate snap builds in CI/CD

- [ ] **Arch User Repository (AUR)**
  - [ ] Create `PKGBUILD` file
  - [ ] Submit to AUR
  - [ ] Documentation for AUR installation

### Code Signing

**Priority: High** (Required for production)

- [ ] **macOS Code Signing**
  - [ ] Apple Developer account setup
  - [ ] Certificate management in CI/CD
  - [ ] Notarization process
  - [ ] Gatekeeper compatibility

- [ ] **Windows Code Signing**
  - [ ] Acquire code signing certificate
  - [ ] Sign executables and installers
  - [ ] SmartScreen reputation building

- [ ] **Linux Code Signing**
  - [ ] GPG signing for .deb packages
  - [ ] GPG signing for .rpm packages
  - [ ] Secure key management

## Medium Term (3-6 Months)

### Features

**Priority: High**

- [ ] **Auto-Update System**
  - [ ] Update checker service
  - [ ] In-app update notifications
  - [ ] Sparkle framework integration (macOS)
  - [ ] Squirrel.Windows integration
  - [ ] Delta updates for faster downloads

- [ ] **Browser Rules & Automation**
  - [ ] URL pattern matching (e.g., always open work URLs in Work profile)
  - [ ] Domain-based browser selection
  - [ ] Time-based rules (work hours → work browser)
  - [ ] Rule editor UI

- [ ] **Enhanced UI**
  - [ ] Dark/light mode toggle
  - [ ] Custom themes
  - [ ] Profile avatars/icons
  - [ ] Browser usage statistics
  - [ ] Search/filter browsers

**Priority: Medium**

- [ ] **Import/Export**
  - [ ] Export preferences to file
  - [ ] Import preferences from file
  - [ ] Sync across machines (iCloud, Dropbox, etc.)

- [ ] **Advanced Browser Support**
  - [ ] Custom browser configuration
  - [ ] Portable browser support
  - [ ] Browser-specific launch flags
  - [ ] Profile nicknames/aliases

- [ ] **Performance Improvements**
  - [ ] Faster browser detection
  - [ ] Icon caching
  - [ ] Lazy loading for profiles
  - [ ] Memory usage optimization

### Distribution

**Priority: Medium**

- [ ] **Flatpak (Linux)**
  - [ ] Create Flatpak manifest
  - [ ] Submit to Flathub
  - [ ] Sandbox configuration

- [ ] **Microsoft Store (Windows)**
  - [ ] MSIX packaging
  - [ ] Store listing preparation
  - [ ] Compliance and certification

- [ ] **Mac App Store (macOS)**
  - [ ] App Store compliance review
  - [ ] Sandboxing adjustments
  - [ ] In-app purchases for pro features (if applicable)

## Long Term (6-12 Months)

### Advanced Features

**Priority: Medium**

- [ ] **Browser Extensions Integration**
  - [ ] Communication with browser extensions
  - [ ] Enhanced profile detection
  - [ ] Sync settings via extension

- [ ] **Mobile Companion Apps**
  - [ ] iOS app to send links to desktop
  - [ ] Android app to send links to desktop
  - [ ] QR code link sharing

- [ ] **Team/Enterprise Features**
  - [ ] Centralized policy management
  - [ ] Organization-wide defaults
  - [ ] Usage analytics for IT admins
  - [ ] MDM integration

**Priority: Low**

- [ ] **Plugin System**
  - [ ] Plugin API for custom integrations
  - [ ] Community plugin marketplace
  - [ ] JavaScript/Go plugin support

- [ ] **Cloud Sync**
  - [ ] Browski Cloud service
  - [ ] Cross-device preference sync
  - [ ] Usage history sync

## Community & Documentation

**Priority: Medium**

- [ ] **Enhanced Documentation**
  - [ ] Video tutorials
  - [ ] FAQ section
  - [ ] Use case examples
  - [ ] Developer documentation for contributors

- [ ] **Community Building**
  - [ ] Discord server
  - [ ] Contributing guidelines
  - [ ] Code of conduct
  - [ ] Issue templates

- [ ] **Localization**
  - [ ] i18n framework integration
  - [ ] Translation management system
  - [ ] Community translation contributions
  - [ ] Support for RTL languages

## Technical Debt & Infrastructure

**Priority: Varies**

- [ ] **Testing**
  - [ ] Unit tests for Go backend (coverage >80%)
  - [ ] Integration tests for browser detection
  - [ ] E2E tests for URL handling
  - [ ] UI tests for React components
  - [ ] Cross-platform testing automation

- [ ] **CI/CD Improvements**
  - [ ] Automated release notes generation
  - [ ] Beta/canary release channels
  - [ ] Performance benchmarking in CI
  - [ ] Security scanning (CodeQL, Snyk)

- [ ] **Build System**
  - [ ] Migrate to Wails v3 stable (when released)
  - [ ] Reproducible builds
  - [ ] Build time optimization
  - [ ] Dependency management automation

## Research & Experimentation

**Priority: Low**

- [ ] **AI-Powered Features**
  - [ ] Smart browser selection based on content
  - [ ] Automatic profile switching suggestions
  - [ ] URL categorization

- [ ] **Browser Integration**
  - [ ] Native browser APIs exploration
  - [ ] Deep linking capabilities
  - [ ] Session restoration

## Breaking Changes & Major Updates

### v2.0 (Future)

Potential breaking changes being considered:

- Configuration file format migration
- New preferences storage backend
- API redesign for better extensibility
- UI framework upgrade

---

## How to Contribute

Want to help implement these features?

1. Check the [Issues](https://github.com/icaliskanoglu/browski/issues) page
2. Look for issues tagged with `help wanted` or `good first issue`
3. Comment on an issue to claim it
4. Submit a pull request with your implementation

## Feedback

Have ideas not on this roadmap?

- Open a [GitHub Discussion](https://github.com/icaliskanoglu/browski/discussions)
- Create a [Feature Request](https://github.com/icaliskanoglu/browski/issues/new?template=feature_request.md)
- Upvote existing feature requests to show interest

---

**Last Updated:** 2026-07-04

This roadmap is subject to change based on user feedback, contributor availability, and project priorities.
