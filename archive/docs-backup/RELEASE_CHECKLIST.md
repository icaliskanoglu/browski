# Release Checklist

## Pre-Release Setup (One-Time)

### 1. GitHub Secrets Configuration
Set up the following secrets in your GitHub repository (Settings → Secrets → Actions):

**For macOS Code Signing** (Optional but recommended):
- `MACOS_CERTIFICATE` - Base64-encoded .p12 certificate
- `MACOS_CERTIFICATE_PWD` - Certificate password
- `MACOS_SIGNING_IDENTITY` - Certificate Common Name
- `APPLE_ID` - Your Apple ID email
- `APPLE_ID_PASSWORD` - App-specific password
- `APPLE_TEAM_ID` - Your Apple Developer Team ID

**Notes:**
- Without these, builds will still work but won't be signed/notarized
- Cost: $99/year for Apple Developer Program

### 2. Update Version Numbers

Before each release, update version in:
- [ ] `build/config.yml` → `info.version`
- [ ] `internal/updater/updater.go` → `currentVersion`
- [ ] `homebrew/browski.rb` → `version`

### 3. Update Changelog
- [ ] Document new features
- [ ] List bug fixes
- [ ] Note breaking changes

## Release Process

### Step 1: Test Locally
```bash
# Build and test
wails3 build
./bin/browski https://github.com

# Verify functionality
# - Browser detection works
# - Profile selection works
# - System tray appears
# - Preferences save correctly
```

### Step 2: Commit Changes
```bash
git add .
git commit -m "Release v1.0.0"
git push origin redesign
```

### Step 3: Create Release Tag
```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0

# This triggers the GitHub Actions release workflow
```

### Step 4: Monitor Build
1. Go to Actions tab in GitHub
2. Watch the "Release" workflow
3. Verify all platforms build successfully:
   - macOS (arm64 + amd64)
   - Windows (amd64)
   - Linux (amd64)

### Step 5: Verify Release Artifacts
Check that GitHub Releases page has:
- [ ] browski-macos-arm64.dmg
- [ ] browski-macos-amd64.dmg
- [ ] browski-windows-amd64-setup.exe
- [ ] browski-linux-amd64.AppImage
- [ ] browski-linux-amd64.deb
- [ ] browski-linux-amd64.rpm

### Step 6: Update Distribution Channels

#### Homebrew (if set up)
```bash
# Calculate SHA256 checksums
shasum -a 256 browski-macos-arm64.dmg
shasum -a 256 browski-macos-amd64.dmg

# Update homebrew-browski repo
cd ../homebrew-browski
vim Formula/browski.rb  # Update version and SHA256s
git commit -am "Update to v1.0.0"
git push
```

#### Chocolatey (if set up)
Update `browski.nuspec` and submit to community repository

#### Snap Store (if set up)
```bash
snapcraft upload --release=stable browski_1.0.0_amd64.snap
```

### Step 7: Announce Release
- [ ] Twitter/X
- [ ] Product Hunt (for major releases)
- [ ] Hacker News
- [ ] Reddit (/r/productivity, /r/browsers)
- [ ] Dev.to or personal blog

## Post-Release

### Monitor
- [ ] GitHub Issues for bug reports
- [ ] Download statistics
- [ ] User feedback

### Update Documentation
- [ ] Update README if needed
- [ ] Update screenshots
- [ ] Update demo video

## Quick Reference Commands

```bash
# Check current version
grep "currentVersion" internal/updater/updater.go

# Test update checker locally
go run main.go

# Build all platforms via GitHub Actions
gh workflow run release.yml -f version=v1.0.0

# View release on GitHub
gh release view v1.0.0

# Delete tag if needed
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

## Troubleshooting

### Build Fails on macOS
- Check that Wails v3 is installed in Actions
- Verify frontend dependencies install correctly
- Check Go version matches (1.21+)

### Code Signing Fails
- Verify GitHub secrets are set correctly
- Check certificate hasn't expired
- Ensure Apple ID password is app-specific, not main password

### Linux Build Fails
- Verify system dependencies are installed in workflow
- Check that AppImage tools are available

## Version Numbering

Follow semantic versioning:
- **Major** (1.0.0): Breaking changes
- **Minor** (0.1.0): New features, backwards compatible
- **Patch** (0.0.1): Bug fixes

Example progression:
- 1.0.0 - Initial release
- 1.1.0 - Added Firefox support
- 1.1.1 - Fixed Safari icon bug
- 2.0.0 - Redesigned UI (breaking change)
