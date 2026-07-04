# Distribution Implementation Tasks

Based on the distribution specification, here's the implementation checklist:

## ✅ Already Implemented

- [x] Single-instance mode (via file locks and IPC - lock_unix.go, lock_windows.go, main.go)
- [x] URL intake abstraction (HandleIncomingURL via SetURL in browserservice.go)
- [x] Version stamping infrastructure (ldflags in CI)
- [x] Native runners in CI (macos-latest, windows-latest, ubuntu-latest)
- [x] macOS CFBundleURLTypes for http/https (build/darwin/Info.plist)
- [x] Basic GitHub Actions release workflow (.github/workflows/release.yml)

## 🚧 In Progress / Needs Updates

### 0. Cross-platform Prerequisites
- [ ] Verify single-instance works on Windows (currently using file locking, may need adjustment)
- [ ] Test URL forwarding on all platforms
- [x] Version injection in builds (already in CI)

### 1. macOS Distribution

#### 1.1 Bundle Configuration
- [ ] Add CFBundleDocumentTypes for public.html to Info.plist
  - File: `build/darwin/Info.plist`
  - Need: Entry for `public.html` with role Viewer
  - Why: Required to appear in System Settings default browser dropdown

- [x] Universal binary build (already configured in CI: darwin/arm64, darwin/amd64)
  - Note: CI builds separate arch binaries, need to add universal build

#### 1.2 Signing & Notarization
- [ ] Set up GitHub Actions secrets:
  - [ ] MACOS_CERT_P12_BASE64
  - [ ] MACOS_CERT_PASSWORD  
  - [ ] APPLE_ID
  - [ ] APPLE_TEAM_ID
  - [ ] APPLE_APP_SPECIFIC_PASSWORD

- [ ] Update macOS CI job in release.yml:
  - [x] Import cert (partially done, needs keychain setup)
  - [ ] Sign with --options runtime --timestamp
  - [ ] Create zip for notarization
  - [ ] Submit to notarytool
  - [ ] Staple ticket to app
  - [ ] Create DMG
  - [ ] Sign and staple DMG
  - [ ] Verify with spctl

#### 1.3 Distribution
- [ ] Create DMG packaging (use create-dmg or custom script)
- [ ] Set up Homebrew tap repository
- [ ] Automate Homebrew cask updates on release

### 2. Windows Distribution

#### 2.1 Default Browser Registration
- [ ] Extend NSIS installer (build/windows/installer/project.nsi if exists)
  - [ ] Add registry keys for RegisteredApplications
  - [ ] Add Capabilities\URLAssociations for http/https
  - [ ] Add ProgID (AppNameURL)
  - [ ] Add uninstall cleanup
  - [ ] Optional: deep-link to ms-settings:defaultapps

- [ ] Verify NSIS task exists (release.yml line 163 fails currently)
  - Current status: `wails3 task windows:package:nsis` doesn't exist
  - Need to: Create the task or use different packaging method

#### 2.2 Code Signing
- [ ] Apply to SignPath.io free OSS program
- [ ] Set up SignPath GitHub Action integration
- [ ] Sign both App.exe and installer
- [ ] Document SmartScreen bypass until signing approved

#### 2.3 Distribution
- [ ] Set up winget manifest automation (vedantmgoyal9/winget-releaser)

### 3. Linux Distribution

#### 3.1 Desktop Integration
- [ ] Create appname.desktop file
  - File: `build/linux/browski.desktop`
  - Include: Exec with %u, MimeType for http/https/text/html
  - Install location: /usr/share/applications/

#### 3.2 Packaging
- [ ] Set up nfpm configuration
  - File: `build/linux/nfpm.yaml` (may exist, need to verify)
  - Outputs: .deb and .rpm
  - Include: .desktop file, icon, postinst script

- [ ] Create AppImage build task
  - File: Task in build/linux/Taskfile.yml
  - Currently referenced in CI but doesn't exist

- [ ] Add postinst script for update-desktop-database

- [ ] Document WebKitGTK dependencies

#### 3.3 Distribution  
- [ ] Generate GPG key for project releases
- [ ] Set up GPG signing for artifacts
- [ ] Create SHA256SUMS and SHA256SUMS.asc

### 4. GitHub Actions Updates

#### 4.1 Workflow Structure
- [x] Multi-platform build jobs (already exists)
- [x] Release creation job (already exists)
- [ ] Add universal build for macOS
- [ ] Add signing steps for all platforms
- [ ] Add checksums generation

#### 4.2 Artifact Management
- [ ] Ensure consistent naming: App-<ver>-<platform>-<arch>.<ext>
- [ ] Fix artifact path in Linux job (currently uses bin/browski-v3)
- [ ] Add SHA256SUMS generation
- [ ] Add GPG signing of checksums

#### 4.3 Automation
- [ ] Pin Wails CLI version (currently using env var)
- [ ] Add caching for Go modules
- [ ] Add caching for frontend node_modules
- [ ] Set release as failed if signing/notarization fails
- [ ] Add Homebrew cask bump automation
- [ ] Add winget manifest bump automation

## 5. Testing & Verification

### Acceptance Criteria Checklist
- [ ] macOS: App appears in System Settings default browser list
- [ ] macOS: spctl -a -vv shows Notarized Developer ID
- [ ] macOS: Link click opens chooser with URL
- [ ] Windows: App in Settings → Default apps
- [ ] Windows: No SmartScreen warning (post-signing)
- [ ] Windows: Link click opens chooser
- [ ] Linux Ubuntu: xdg-settings works
- [ ] Linux Fedora: Package installs and works
- [ ] Linux: AppImage runs on clean VM
- [ ] All platforms: Second link doesn't spawn second window
- [ ] All platforms: --version prints git tag

## Priority Order

### Phase 1: Core Functionality (Current Sprint)
1. Add CFBundleDocumentTypes to macOS Info.plist
2. Fix Linux binary name in CI (browski-v3 → browski)
3. Create .desktop file for Linux
4. Verify single-instance URL forwarding on Windows

### Phase 2: Packaging (Next Sprint)
1. Create DMG packaging for macOS
2. Set up NSIS installer tasks for Windows
3. Set up nfpm for Linux .deb/.rpm
4. Create AppImage build task

### Phase 3: Code Signing (Following Sprint)
1. Set up macOS signing and notarization
2. Apply for SignPath.io and set up Windows signing
3. Set up GPG signing for Linux artifacts

### Phase 4: Automation (Final Sprint)
1. Homebrew tap automation
2. Winget manifest automation
3. Checksums and signatures in releases
4. Comprehensive testing on all platforms

---

**Next Steps:**
1. Start with Phase 1 tasks
2. Update this document as tasks are completed
3. Create GitHub issues for major items
4. Test each platform thoroughly before moving to next phase
