# Release Process

This document describes how to create a new release of Browski using semantic versioning.

## Semantic Versioning

Browski follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backward compatible manner
- **PATCH** version for backward compatible bug fixes

Current version is tracked in the `VERSION` file at the project root.

## Creating a Release

### Using GitHub Actions (Recommended)

1. Go to the [Actions tab](../../actions) in the GitHub repository
2. Select the "Release" workflow
3. Click "Run workflow"
4. Select the version bump type:
   - **patch**: For bug fixes (1.0.0 → 1.0.1)
   - **minor**: For new features (1.0.0 → 1.1.0)
   - **major**: For breaking changes (1.0.0 → 2.0.0)
5. Optionally check "Skip build" if you only want to bump the version without building
6. Click "Run workflow"

The workflow will automatically:
- Read the current version from `VERSION` file
- Calculate the new version based on bump type
- Update `VERSION`, `build/config.yml`, and `build/windows/info.json`
- Commit the changes with message: `chore: bump version to X.Y.Z`
- Create and push a git tag (e.g., `v1.0.1`)
- Build for macOS (Intel + Apple Silicon), Windows, and Linux
- Create a GitHub release with all build artifacts

### Manual Release

If you prefer to create a release manually:

1. Update the version in the following files:
   - `VERSION`
   - `build/config.yml` (line 14)
   - `build/windows/info.json` (lines 3 and 7)

2. Commit the changes:
   ```bash
   git add VERSION build/config.yml build/windows/info.json
   git commit -m "chore: bump version to X.Y.Z"
   git push origin main
   ```

3. Create and push a tag:
   ```bash
   git tag vX.Y.Z
   git push origin vX.Y.Z
   ```

4. The Release workflow will automatically trigger on the new tag and build all artifacts.

## Release Artifacts

Each release includes the following artifacts:

### macOS
- `browski-macos-amd64.dmg` - Intel Macs
- `browski-macos-arm64.dmg` - Apple Silicon Macs
- Signed and notarized (if certificates configured)

### Windows
- `browski-arm64-installer.exe` - NSIS installer with WebView2 bootstrapper
- `browski-amd64.msix` - MSIX package (optional)

### Linux
- `browski-amd64.AppImage` - Portable AppImage
- `browski-amd64.deb` - Debian/Ubuntu package
- `browski-amd64.rpm` - Fedora/RHEL package

## Version Compatibility

- **Wails**: v3.0.0-alpha2.111+
- **Go**: 1.25+
- **Node.js**: 22+

## Release Notes

Release notes are automatically generated from the template in `.github/workflows/release.yml`. You can edit them after the release is created by going to the Releases page and clicking "Edit release".

## Troubleshooting

### Build fails on macOS
- Ensure code signing certificates are configured in repository secrets
- Check that `MACOS_CERTIFICATE`, `MACOS_CERTIFICATE_PWD`, and `MACOS_SIGNING_IDENTITY` are set

### Build fails on Windows
- NSIS is automatically installed via Chocolatey
- Ensure `build/windows/icon.ico` exists and is valid

### Version not updated
- Check that the workflow has write permissions to the repository
- Verify `GITHUB_TOKEN` has sufficient permissions

### Build succeeds but no release created
- Check that all platform builds completed successfully
- Verify the `create-release` job conditions are met
