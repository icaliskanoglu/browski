# Distribution Implementation Guide

Complete guide for packaging, distributing, and releasing Browski across all platforms.

**Last Updated**: 2026-07-04

---

## Table of Contents

1. [Overview](#overview)
2. [Implementation Status](#implementation-status)
3. [Phase-by-Phase Plan](#phase-by-phase-plan)
4. [Platform-Specific Details](#platform-specific-details)
5. [CI/CD Integration](#cicd-integration)
6. [Testing & Verification](#testing--verification)
7. [Distribution Channels](#distribution-channels)

---

## Overview

Browski uses GitHub Actions for automated builds and releases across macOS, Windows, and Linux. The distribution strategy focuses on:

- **Native installers** for each platform (DMG, NSIS, deb/rpm)
- **Code signing** for security and trust
- **Package managers** for easy installation (Homebrew, Chocolatey, Snap, etc.)
- **Auto-updates** for seamless user experience
- **Single-instance** behavior with URL forwarding

### Quick Release Process

```bash
# Tag and push
git tag v1.0.0
git push origin v1.0.0

# Or use GitHub CLI
gh release create v1.0.0 --generate-notes
```

This triggers automated builds for all platforms.

---

## Implementation Status

### ✅ Already Completed

- [x] **Single-instance mode** - File locks (lock_unix.go, lock_windows.go) + IPC via Unix sockets
- [x] **URL intake abstraction** - HandleIncomingURL via SetURL in browserservice.go
- [x] **Version stamping** - ldflags injection in CI workflow
- [x] **Multi-platform CI** - Native runners for macOS, Windows, Linux
- [x] **macOS URL handlers** - CFBundleURLTypes for http/https in Info.plist
- [x] **macOS HTML viewer** - CFBundleDocumentTypes with LSItemContentTypes
- [x] **Basic release workflow** - .github/workflows/release.yml
- [x] **Binary naming** - Consistent "browski" across all platforms

### 🚧 In Progress / Pending

#### Cross-Platform
- [ ] Verify single-instance works on Windows
- [ ] Test URL forwarding on all platforms
- [ ] Create SHA256SUMS for all releases
- [ ] Set up GPG signing for checksums

#### macOS
- [ ] DMG packaging
- [ ] Code signing and notarization
- [ ] Homebrew tap repository
- [ ] Auto-update mechanism

#### Windows
- [ ] NSIS installer with registry keys
- [ ] Code signing (via SignPath.io)
- [ ] winget manifest automation

#### Linux
- [ ] .desktop file
- [ ] nfpm configuration (deb/rpm)
- [ ] AppImage build
- [ ] Snap/Flatpak packages
- [ ] AUR PKGBUILD

---

## Phase-by-Phase Plan

### Phase 1: Core Functionality ✅ → 🚧

**Goal**: Ensure Browski properly registers as a browser on all platforms.

**Status**: 2/4 tasks complete

| Task | Status | Priority |
|------|--------|----------|
| macOS Info.plist updates | ✅ DONE | HIGH |
| Linux binary name fix in CI | ✅ DONE | HIGH |
| Create Linux .desktop file | 🚧 PENDING | HIGH |
| Verify Windows single-instance | 🚧 PENDING | MEDIUM |

**Time Estimate**: 3-4 hours remaining

#### Task 3: Linux .desktop File

**Create**: `build/linux/browski.desktop`

```desktop
[Desktop Entry]
Name=Browski
Comment=Browser profile picker
Exec=browski %u
Icon=browski
Type=Application
Categories=Network;WebBrowser;
MimeType=text/html;x-scheme-handler/http;x-scheme-handler/https;
StartupNotify=true
```

**Install location**: /usr/share/applications/

#### Task 4: Windows Single-Instance Testing

**Steps**:
1. Build Windows binary
2. Launch first instance with URL
3. Launch second instance with different URL
4. Verify URL forwarded to first instance
5. Check only one window appears

**Files**: lock_windows.go, main.go

---

### Phase 2: Packaging 🚧

**Goal**: Create native installers for each platform.

**Status**: 0/4 tasks complete

| Task | Status | Priority | Estimate |
|------|--------|----------|----------|
| macOS DMG packaging | 🚧 PENDING | HIGH | 2-3h |
| Windows NSIS installer | 🚧 PENDING | HIGH | 3-4h |
| Linux nfpm (deb/rpm) | 🚧 PENDING | MEDIUM | 1-2h |
| Linux AppImage | 🚧 PENDING | LOW | 1-2h |

#### macOS DMG Task

**Create**: Taskfile task `darwin:package:dmg`

```yaml
darwin:package:dmg:
  desc: "Create macOS DMG installer"
  cmds:
    - |
      create-dmg \
        --volname "Browski" \
        --volicon "build/darwin/icons.icns" \
        --window-pos 200 120 \
        --window-size 600 400 \
        --icon-size 100 \
        --icon "Browski.app" 175 120 \
        --hide-extension "Browski.app" \
        --app-drop-link 425 120 \
        "bin/Browski-{{.VERSION}}.dmg" \
        "bin/Browski.app"
```

**Dependencies**: create-dmg tool

#### Windows NSIS Installer

**Current Issue**: `wails3 task windows:package:nsis` doesn't exist (CI line 163 fails)

**Create**: `build/windows/installer/installer.nsi`

**Registry Keys Required**:
```
HKLM\SOFTWARE\RegisteredApplications
  Browski = Software\Browski\Capabilities

HKLM\SOFTWARE\Browski\Capabilities
  ApplicationName = Browski
  ApplicationDescription = Browser profile picker

HKLM\SOFTWARE\Browski\Capabilities\URLAssociations
  http = BrowskiURL
  https = BrowskiURL

HKLM\SOFTWARE\Classes\BrowskiURL
  (Default) = "Browski URL"
  URL Protocol = ""

HKLM\SOFTWARE\Classes\BrowskiURL\shell\open\command
  (Default) = "C:\Program Files\Browski\browski.exe" "%1"
```

**Optional**: Deep-link to `ms-settings:defaultapps`

#### Linux nfpm Configuration

**Create**: `build/linux/nfpm.yaml`

**Outputs**: .deb and .rpm packages

**Contents**:
- Binary: /usr/bin/browski
- .desktop file: /usr/share/applications/browski.desktop
- Icon: /usr/share/pixmaps/browski.png
- postinst script: update-desktop-database

**Create**: `build/linux/postinst.sh`

```bash
#!/bin/sh
if [ -x /usr/bin/update-desktop-database ]; then
    update-desktop-database -q /usr/share/applications
fi
```

#### Linux AppImage

**Status**: Referenced in CI but task doesn't exist

**Note**: Can defer to later sprint if time is limited

---

### Phase 3: Code Signing 🔒

**Goal**: Sign binaries for security and to eliminate warning dialogs.

**Status**: 0/3 tasks complete

| Platform | Task | Cost | Time Estimate |
|----------|------|------|---------------|
| macOS | Notarization | $99/year | 4-6h |
| Windows | SignPath.io | Free (OSS) | 2-3h |
| Linux | GPG signing | Free | 1h |

#### macOS Signing & Notarization

**GitHub Secrets Required**:
- `MACOS_CERT_P12_BASE64`
- `MACOS_CERT_PASSWORD`
- `APPLE_ID`
- `APPLE_TEAM_ID`
- `APPLE_APP_SPECIFIC_PASSWORD`

**CI Steps**:
1. Import certificate
2. Sign with `--options runtime --timestamp`
3. Create zip for notarization
4. Submit to `notarytool`
5. Staple ticket to app
6. Create and sign DMG
7. Staple DMG
8. Verify with `spctl -a -vv`

**Cost**: $99/year (Apple Developer Program)

#### Windows Signing

**Service**: SignPath.io free OSS program

**Steps**:
1. Apply to SignPath.io
2. Set up GitHub Action integration
3. Sign both browski.exe and installer
4. Document SmartScreen bypass until approved

**Note**: Until signed, users see SmartScreen warning

#### Linux GPG Signing

**Steps**:
1. Generate GPG key for project
2. Set up GPG signing in CI
3. Create SHA256SUMS
4. Sign with GPG → SHA256SUMS.asc

**Upload to releases**:
- SHA256SUMS
- SHA256SUMS.asc

---

### Phase 4: Automation 🤖

**Goal**: Automate package manager updates and releases.

**Status**: 0/4 tasks complete

| Task | Tool/Service | Estimate |
|------|-------------|----------|
| Homebrew tap automation | GitHub Actions | 2h |
| winget manifest automation | vedantmgoyal9/winget-releaser | 1h |
| Checksums generation | GitHub Actions | 1h |
| Comprehensive testing | Manual/automated | 4h |

---

## Platform-Specific Details

### macOS Distribution

**Bundle Requirements**:
- ✅ CFBundleURLTypes for http/https
- ✅ CFBundleDocumentTypes for public.html
- ✅ LSItemContentTypes array
- ✅ Universal binary (arm64 + amd64)

**Distribution Channels**:
1. GitHub Releases (DMG)
2. Homebrew tap
3. Optional: Mac App Store (requires App Store Connect)

**Acceptance Criteria**:
- [ ] App appears in System Settings → Default Browser
- [ ] spctl shows "Notarized Developer ID"
- [ ] Link click opens chooser with URL
- [ ] Second link doesn't spawn second window

### Windows Distribution

**Registry Requirements** (see Phase 2 NSIS section)

**Distribution Channels**:
1. GitHub Releases (installer)
2. winget (Windows Package Manager)
3. Optional: Chocolatey, Microsoft Store

**Acceptance Criteria**:
- [ ] App in Settings → Default apps
- [ ] No SmartScreen warning (post-signing)
- [ ] Link click opens chooser
- [ ] Registry keys properly set
- [ ] Single-instance works

### Linux Distribution

**Desktop Integration**:
- .desktop file in /usr/share/applications/
- Icon in /usr/share/pixmaps/
- MIME type associations for http/https/html

**Package Formats**:
1. .deb (Debian/Ubuntu)
2. .rpm (Fedora/RHEL)
3. AppImage (universal)
4. Snap (Snap Store)
5. Flatpak (Flathub)
6. AUR (Arch User Repository)

**Dependencies**:
- libgtk-4-dev
- libwebkitgtk-6.0-dev

**Acceptance Criteria**:
- [ ] Ubuntu: xdg-settings works
- [ ] Fedora: Package installs and works
- [ ] AppImage runs on clean VM
- [ ] Single-instance works

---

## CI/CD Integration

### Current Workflow

`.github/workflows/release.yml` triggers on:
- Tag push: `git push origin v*`
- Manual dispatch: workflow_dispatch

**Jobs**:
1. `build-macos` - Builds for arm64 and amd64
2. `build-windows` - Builds for amd64
3. `build-linux` - Builds for amd64
4. `create-release` - Publishes artifacts

### Workflow Improvements Needed

**build-macos job**:
- [ ] Add DMG creation step
- [ ] Add notarization step
- [ ] Add universal binary build

**build-windows job**:
- [x] Fix binary renaming (browski → browski)
- [ ] Add NSIS packaging step
- [ ] Add code signing step

**build-linux job**:
- [x] Fix artifact path (browski → browski)
- [ ] Add .desktop file packaging
- [ ] Add deb/rpm generation
- [ ] Add AppImage creation

**create-release job**:
- [ ] Add checksum generation
- [ ] Add GPG signing
- [ ] Update release notes template

### Artifact Naming Convention

**Standard**: `browski-<version>-<platform>-<arch>.<ext>`

**Examples**:
- `browski-1.0.0-darwin-arm64.dmg`
- `browski-1.0.0-darwin-amd64.dmg`
- `browski-1.0.0-windows-amd64.exe`
- `browski-1.0.0-linux-amd64.deb`
- `browski-1.0.0-linux-amd64.rpm`
- `browski-1.0.0-linux-amd64.AppImage`

---

## Testing & Verification

### Pre-Release Testing

**macOS**:
```bash
wails3 task darwin:build
wails3 task darwin:package
wails3 task darwin:package:dmg  # Once implemented
# Install from DMG
# Test URL handling
# Verify in System Settings
```

**Windows**:
```bash
wails3 build
# Run NSIS installer (once implemented)
# Check registry keys
# Verify in Settings → Default Apps
# Test single-instance
```

**Linux**:
```bash
wails3 build
sudo install -Dm755 bin/browski /usr/local/bin/
sudo install -Dm644 build/linux/browski.desktop /usr/share/applications/
update-desktop-database
xdg-settings set default-web-browser browski.desktop
# Test URL handling
```

### Acceptance Criteria Checklist

**All Platforms**:
- [ ] --version prints git tag
- [ ] Second link doesn't spawn second window
- [ ] URL forwarding works between instances
- [ ] Browser detection finds installed browsers
- [ ] Profile selection works correctly

**Platform-Specific**: See sections above

---

## Distribution Channels

### GitHub Releases (Primary)

**Files Published**:
- macOS: `browski-<ver>-darwin-{arm64,amd64}.dmg`
- Windows: `browski-<ver>-windows-amd64.exe`
- Linux: `browski-<ver>-linux-amd64.{deb,rpm,AppImage}`
- Checksums: `SHA256SUMS`, `SHA256SUMS.asc`

**Pros**: Free, version control, easy
**Cons**: Manual download required

### Homebrew (macOS)

**Setup**:
```bash
# Create tap repository
gh repo create homebrew-browski --public

# Create formula
cat > Formula/browski.rb <<EOF
class Browski < Formula
  desc "Browser profile picker"
  homepage "https://github.com/icaliskanoglu/browski"
  version "1.0.0"

  if Hardware::CPU.arm?
    url "https://github.com/icaliskanoglu/browski/releases/download/v1.0.0/browski-1.0.0-darwin-arm64.dmg"
    sha256 "REPLACE_WITH_ACTUAL_SHA"
  else
    url "https://github.com/icaliskanoglu/browski/releases/download/v1.0.0/browski-1.0.0-darwin-amd64.dmg"
    sha256 "REPLACE_WITH_ACTUAL_SHA"
  end

  def install
    prefix.install Dir["*"]
  end
end
EOF
```

**Users install**:
```bash
brew tap icaliskanoglu/browski
brew install browski
```

**Automation**: Update formula on each release (Phase 4)

### winget (Windows)

**Automation**: Use `vedantmgoyal9/winget-releaser` GitHub Action

**Users install**:
```powershell
winget install Browski
```

### Chocolatey (Windows)

**Create**: `browski.nuspec`
**Submit to**: https://community.chocolatey.org/

**Users install**:
```powershell
choco install browski
```

### Snap Store (Linux)

**Create**: `snap/snapcraft.yaml`

**Build and publish**:
```bash
snapcraft
snapcraft login
snapcraft upload --release=stable browski_1.0.0_amd64.snap
```

**Users install**:
```bash
snap install browski
```

### AUR (Arch Linux)

**Create**: `PKGBUILD`
**Submit to**: AUR repository

**Users install**:
```bash
yay -S browski
```

---

## Budget & Resources

### Costs

| Item | Cost | Frequency | Required? |
|------|------|-----------|-----------|
| Apple Developer | $99 | Annual | Yes (for notarization) |
| Windows Code Signing | Free* | - | Yes (SignPath OSS) |
| Domain (optional) | $12 | Annual | No |
| **Total** | ~$111 | Annual | |

*SignPath.io is free for open-source projects

### Free Alternatives

- Skip Mac App Store (use Homebrew)
- Skip Microsoft Store (use winget)
- Use GitHub Pages for website

---

## Support & Marketing

### Support Channels

- **GitHub Issues** - Bug reports, feature requests
- **GitHub Discussions** - Q&A, community
- **Email** - Optional
- **Discord/Slack** - Community-driven

### Marketing

**Announce on**:
- Product Hunt
- Hacker News
- Reddit (/r/productivity, /r/browsers)
- Twitter/X
- Dev.to

**Create**:
- Demo video
- Screenshots
- Website (GitHub Pages)
- README badges

---

## Next Steps

### Immediate (Phase 1)

1. [ ] Create Linux .desktop file
2. [ ] Test Windows single-instance behavior
3. [ ] Verify all acceptance criteria for Phase 1

### Short-term (Phase 2)

1. [ ] Implement DMG packaging
2. [ ] Create NSIS installer
3. [ ] Set up nfpm for deb/rpm
4. [ ] Test all packages

### Medium-term (Phase 3)

1. [ ] Set up Apple Developer account
2. [ ] Configure macOS notarization
3. [ ] Apply for SignPath.io
4. [ ] Set up GPG signing

### Long-term (Phase 4)

1. [ ] Automate Homebrew updates
2. [ ] Automate winget updates
3. [ ] Create Snap package
4. [ ] Submit to AUR
5. [ ] Consider Flatpak

---

## References

- [Wails v3 Documentation](https://v3alpha.wails.io/)
- [macOS Notarization Guide](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Windows App Registration](https://docs.microsoft.com/en-us/windows/win32/shell/default-programs)
- [Desktop Entry Specification](https://specifications.freedesktop.org/desktop-entry-spec/latest/)
- [SignPath.io](https://signpath.io/)
- [create-dmg](https://github.com/create-dmg/create-dmg)
- [nfpm](https://nfpm.goreleaser.com/)

---

**Maintained by**: Browski Development Team
**Last Review**: 2026-07-04
