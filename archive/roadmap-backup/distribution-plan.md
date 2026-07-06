# Distribution Implementation Plan
**Date**: 2026-07-04

## Goal
Complete Phase 1 (Core Functionality) of the distribution implementation and begin Phase 2 (Packaging).

## Phase 1: Core Functionality ✅ → 🚧

### Task 1: macOS Info.plist Updates ✅
**Status**: COMPLETED
- ✅ Add CFBundleDocumentTypes for public.html with LSItemContentTypes
- ✅ Add .htm extension support
- ✅ Update type name to "HTML Document"

**Files Modified**:
- `build/darwin/Info.plist`

### Task 2: Fix Linux Binary Name in CI ✅
**Status**: COMPLETED
- ✅ Update release.yml line 255 from `bin/browski` to `bin/browski`

**Files Modified**:
- `.github/workflows/release.yml`

### Task 3: Create .desktop File for Linux 🚧
**Status**: PENDING
- [ ] Create `build/linux/browski.desktop` file
- [ ] Include proper Exec, Icon, MimeType entries
- [ ] Add categories and keywords
- [ ] Configure for /usr/share/applications/ installation

**Implementation Details**:
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

**Files to Create**:
- `build/linux/browski.desktop`

### Task 4: Verify Single-Instance URL Forwarding on Windows 🚧
**Status**: PENDING (Requires Windows machine or VM)
- [ ] Test file locking on Windows (lock_windows.go)
- [ ] Verify IPC communication works
- [ ] Test URL forwarding between instances
- [ ] Document any Windows-specific issues

**Testing Steps**:
1. Build Windows binary
2. Launch first instance with URL
3. Launch second instance with different URL
4. Verify second URL is forwarded to first instance
5. Check that only one window appears

---

## Phase 2: Packaging (Begin Today if Time Permits)

### Task 1: Create DMG Packaging for macOS 🚧
**Priority**: HIGH
**Status**: PENDING

**Implementation Steps**:
1. Install create-dmg tool or use custom DMG script
2. Create Taskfile task: `darwin:package:dmg`
3. Configure DMG appearance (background, icon positions)
4. Set up volume name and window size
5. Update CI workflow to build DMG

**Files to Create/Modify**:
- `Taskfile.yml` (add darwin:package:dmg task)
- `build/darwin/dmg-background.png` (optional)
- `build/darwin/dmg-spec.json` (DMG configuration)
- `.github/workflows/release.yml` (update macOS job)

**Task Definition** (to add to Taskfile.yml):
```yaml
darwin:package:dmg:
  desc: "Create macOS DMG installer"
  cmds:
    - |
      # Create DMG using create-dmg
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

### Task 2: Set Up NSIS Installer for Windows 🚧
**Priority**: HIGH
**Status**: PENDING (CI currently fails on this task)

**Current Issue**:
- `.github/workflows/release.yml` line 163 calls `wails3 task windows:package:nsis`
- This task does not exist in Taskfile.yml

**Implementation Steps**:
1. Create `build/windows/installer/installer.nsi` script
2. Add Windows registry keys for default browser registration
3. Add Taskfile task: `windows:package:nsis`
4. Include uninstaller with registry cleanup
5. Optional: Add deep-link to ms-settings:defaultapps

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

**Files to Create/Modify**:
- `build/windows/installer/installer.nsi`
- `Taskfile.yml` (add windows:package:nsis task)

### Task 3: Set Up nfpm for Linux Packages 🚧
**Priority**: MEDIUM
**Status**: PENDING

**Implementation Steps**:
1. Create `build/linux/nfpm.yaml` configuration
2. Configure .deb and .rpm outputs
3. Add .desktop file to package
4. Create postinst script for update-desktop-database
5. Add Taskfile tasks for deb/rpm building

**Files to Create/Modify**:
- `build/linux/nfpm.yaml`
- `build/linux/postinst.sh` (post-installation script)
- `Taskfile.yml` (add linux:package:deb and linux:package:rpm tasks)

### Task 4: Create AppImage Build Task 🚧
**Priority**: LOW (can defer)
**Status**: PENDING

**Note**: Referenced in CI but task doesn't exist. Can be deferred to later sprint.

---

## Time Estimates

### Phase 1 Remaining Tasks (3-4 hours)
- Task 3 (Linux .desktop): 30 mins
- Task 4 (Windows testing): 2-3 hours (requires Windows environment)

### Phase 2 Tasks (6-8 hours)
- DMG packaging: 2-3 hours
- NSIS installer: 3-4 hours
- nfpm setup: 1-2 hours
- AppImage: 1-2 hours (deferred)

---

## Today's Realistic Goals

### Must Complete ✅
1. ✅ macOS Info.plist updates (DONE)
2. ✅ Linux binary name fix (DONE)
3. Linux .desktop file creation
4. DMG packaging task

### Should Complete 🎯
5. NSIS installer setup
6. nfpm configuration

### Nice to Have ⭐
7. Windows single-instance testing
8. AppImage build task

---

## Testing Plan

After implementation, test each platform:

### macOS Testing
- [ ] Build with `wails3 task darwin:build`
- [ ] Package with `wails3 task darwin:package`
- [ ] Create DMG with `wails3 task darwin:package:dmg`
- [ ] Install from DMG
- [ ] Verify appears in System Settings → Default Browser
- [ ] Test URL handling

### Windows Testing (if environment available)
- [ ] Build with `wails3 build`
- [ ] Run NSIS installer
- [ ] Check registry keys created
- [ ] Verify in Settings → Default Apps
- [ ] Test single-instance behavior
- [ ] Test URL forwarding

### Linux Testing
- [ ] Build with `wails3 build`
- [ ] Install .desktop file
- [ ] Run `update-desktop-database`
- [ ] Test `xdg-settings set default-web-browser browski.desktop`
- [ ] Verify in applications menu
- [ ] Test URL handling

---

## Notes

- Phase 1 Tasks 1-2 were completed before plan creation
- Windows testing may require VM setup or GitHub Actions runs
- Code signing (Phase 3) requires certificate setup - deferred
- Package manager automation (Phase 4) requires packages to exist first

---

## Success Criteria for Today

By end of day, we should have:
1. ✅ macOS properly configured to appear in default browser list
2. ✅ Linux binary name consistent across all builds
3. Linux .desktop file created and tested
4. DMG packaging working and producing signed (or signable) DMGs
5. NSIS installer creating proper Windows installer with registry keys
6. CI workflows updated to use new packaging tasks

If all above are complete, Browski will be ready for Phase 3 (code signing) and Phase 4 (distribution automation).
