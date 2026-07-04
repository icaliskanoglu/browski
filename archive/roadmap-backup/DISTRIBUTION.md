# Distribution Guide

## Overview

Browski uses GitHub Actions for automated builds and releases across all platforms.

## Release Process

### 1. Create a Release

```bash
# Tag and push
git tag v1.0.0
git push origin v1.0.0

# Or use GitHub CLI
gh release create v1.0.0 --generate-notes
```

This triggers the release workflow which:
- Builds for macOS (Intel + Apple Silicon)
- Builds for Windows (x64)
- Builds for Linux (x64)
- Signs binaries (if certificates configured)
- Creates installers
- Publishes to GitHub Releases

### 2. Manual Trigger

You can also trigger builds manually:
```bash
# Via GitHub Actions UI
Actions → Release → Run workflow → Enter version

# Via GitHub CLI
gh workflow run release.yml -f version=v1.0.0
```

## Distribution Channels

### GitHub Releases (Recommended)
**Pros:** Free, easy, version control
**Cons:** Users must download manually

Files published:
- `browski-darwin-arm64.zip` (macOS Apple Silicon)
- `browski-darwin-amd64.zip` (macOS Intel)
- `browski-darwin-universal.zip` (macOS Universal)
- `browski-windows-amd64.exe` (Windows installer)
- `browski-linux-amd64` (Linux binary)
- `browski-linux-amd64.deb`
- `browski-linux-amd64.rpm`

### Homebrew (macOS)

#### Setup Tap
```bash
# Create homebrew-browski repository
gh repo create homebrew-browski --public

# Copy formula
cp homebrew/browski.rb ../homebrew-browski/Formula/
cd ../homebrew-browski
git add Formula/browski.rb
git commit -m "Add Browski formula"
git push
```

#### Update Formula After Release
```bash
# Generate SHA256 checksums
shasum -a 256 browski-darwin-arm64.zip
shasum -a 256 browski-darwin-amd64.zip

# Update formula with new URLs and SHA256s
vim Formula/browski.rb
git commit -am "Update to v1.0.0"
git push
```

#### Users Install Via
```bash
brew tap icaliskanoglu/browski
brew install browski
```

### Chocolatey (Windows)

Create `browski.nuspec`:
```xml
<?xml version="1.0"?>
<package xmlns="http://schemas.microsoft.com/packaging/2011/08/nuspec.xsd">
  <metadata>
    <id>browski</id>
    <version>1.0.0</version>
    <title>Browski</title>
    <authors>Your Name</authors>
    <description>Browser Profile Picker</description>
    <projectUrl>https://github.com/icaliskanoglu/browski</projectUrl>
    <licenseUrl>https://github.com/icaliskanoglu/browski/blob/master/LICENSE</licenseUrl>
    <tags>browser picker profile</tags>
  </metadata>
</package>
```

Submit to https://community.chocolatey.org/

Users install:
```powershell
choco install browski
```

### Snap Store (Linux)

Create `snap/snapcraft.yaml`:
```yaml
name: browski
version: '1.0.0'
summary: Browser Profile Picker
description: |
  A tool for quickly selecting which browser and profile to open URLs with.
grade: stable
confinement: strict

apps:
  browski:
    command: bin/browski
    plugs: [network, desktop, desktop-legacy, wayland, x11]

parts:
  browski:
    plugin: dump
    source: .
    stage-packages:
      - libgtk-3-0
      - libwebkit2gtk-4.0-37
```

Build and publish:
```bash
snapcraft
snapcraft login
snapcraft upload --release=stable browski_1.0.0_amd64.snap
```

### AUR (Arch Linux)

Create `PKGBUILD`:
```bash
pkgname=browski
pkgver=1.0.0
pkgrel=1
pkgdesc="Browser Profile Picker"
arch=('x86_64')
url="https://github.com/icaliskanoglu/browski"
license=('MIT')
depends=('gtk3' 'webkit2gtk')
source=("$pkgname-$pkgver.tar.gz::https://github.com/icaliskanoglu/browski/archive/v$pkgver.tar.gz")
sha256sums=('REPLACE')

package() {
    cd "$srcdir/$pkgname-$pkgver"
    install -Dm755 bin/browski "$pkgdir/usr/bin/browski"
}
```

Submit to AUR.

## Auto-Updates

### macOS with Sparkle

Requires implementing Sparkle framework for auto-updates. See SPARKLE_SETUP.md.

### Windows with Squirrel

Windows auto-updates via Squirrel.Windows. See SQUIRREL_SETUP.md.

### Linux

Linux users typically update via package manager:
```bash
# Homebrew
brew upgrade browski

# Snap
snap refresh browski

# Flatpak
flatpak update com.browski.app
```

## Analytics (Optional)

Track downloads with:
1. GitHub API (releases)
2. Google Analytics (website)
3. Homebrew analytics

## Marketing

### Announce On:
- Product Hunt
- Hacker News
- Reddit (/r/productivity, /r/browsers)
- Twitter/X
- Dev.to
- Your blog

### Create:
- Demo video
- Screenshots
- Website (GitHub Pages)
- README badges (downloads, version, etc.)

## Metrics to Track

- **GitHub**: Stars, forks, issues, downloads per release
- **Homebrew**: `brew info browski` shows install count
- **Snap**: Dashboard shows active installs
- **Website**: Visitors, conversion rate

## Budget

| Item | Cost | Frequency |
|------|------|-----------|
| Apple Developer | $99 | Annual |
| Code Signing (Windows) | $200-400 | Annual |
| Domain (optional) | $12 | Annual |
| **Total** | ~$311 | Annual |

Free alternatives:
- Skip Mac App Store (use GitHub Releases + Homebrew)
- Skip Windows Store (use GitHub Releases + Chocolatey)
- Use GitHub Pages for website

## Support

Provide support via:
- GitHub Issues (bug reports, feature requests)
- GitHub Discussions (Q&A, community)
- Email (optional)
- Discord/Slack (community driven)

## Next Steps

1. [ ] Push code to GitHub
2. [ ] Create first release (v1.0.0)
3. [ ] Set up Homebrew tap
4. [ ] Submit to Chocolatey (optional)
5. [ ] Create website with GitHub Pages
6. [ ] Announce on Product Hunt
7. [ ] Set up code signing certificates
8. [ ] Implement auto-updates
