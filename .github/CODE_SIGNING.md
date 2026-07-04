# Code Signing Setup

## macOS Code Signing

### Prerequisites
1. Apple Developer Account ($99/year)
2. Developer ID Application certificate
3. Developer ID Installer certificate (for .pkg)

### Steps

#### 1. Create Certificates in Apple Developer Portal
1. Go to https://developer.apple.com/account/resources/certificates/
2. Create "Developer ID Application" certificate
3. Download and install in Keychain Access

#### 2. Export Certificates
```bash
# Export as .p12 file from Keychain Access
# Use a strong password
```

#### 3. Convert to Base64 for GitHub Secrets
```bash
base64 -i DeveloperID_Application.p12 -o cert.base64.txt
cat cert.base64.txt | pbcopy  # Copy to clipboard
```

#### 4. Add GitHub Secrets
Go to your repository Settings → Secrets and variables → Actions:

- `MACOS_CERTIFICATE`: Base64-encoded .p12 certificate
- `MACOS_CERTIFICATE_PWD`: Certificate password
- `APPLE_ID`: Your Apple ID email
- `APPLE_ID_PASSWORD`: App-specific password (create at appleid.apple.com)
- `APPLE_TEAM_ID`: Your Team ID (find in Apple Developer account)

#### 5. Update GitHub Actions Workflow
The release workflow will automatically:
- Import certificates
- Sign the .app bundle
- Notarize with Apple
- Create signed .dmg

### Notarization

Notarization is required for macOS apps to pass Gatekeeper. The workflow handles this automatically using:
```bash
xcrun notarytool submit \
  --apple-id "$APPLE_ID" \
  --password "$APPLE_ID_PASSWORD" \
  --team-id "$APPLE_TEAM_ID" \
  --wait \
  browski.dmg
```

## Windows Code Signing

### Prerequisites
1. Code signing certificate (e.g., from DigiCert, Sectigo)
2. Hardware token or cloud signing service

### Steps

#### 1. Obtain Certificate
Purchase from a trusted CA like:
- DigiCert
- Sectigo
- GlobalSign

#### 2. Add GitHub Secrets
- `WINDOWS_CERTIFICATE`: Base64-encoded .pfx certificate
- `WINDOWS_CERTIFICATE_PWD`: Certificate password

#### 3. Sign Binary
```bash
signtool sign /f certificate.pfx \
  /p PASSWORD \
  /tr http://timestamp.digicert.com \
  /td sha256 \
  /fd sha256 \
  browski.exe
```

## Linux Code Signing

Linux doesn't require code signing, but you can:
1. Sign packages with GPG
2. Provide checksums (SHA256)
3. Publish to trusted repositories

### Generate Checksums
```bash
sha256sum browski > browski.sha256
sha256sum browski.AppImage > browski.AppImage.sha256
```

## Testing Signed Builds

### macOS
```bash
# Verify code signature
codesign -dv --verbose=4 Browski.app

# Verify notarization
spctl -a -vvv -t execute Browski.app

# Check Gatekeeper approval
xattr -l Browski.app
```

### Windows
```bash
# Verify signature
signtool verify /pa browski.exe
```

## Troubleshooting

### macOS: "App is damaged and can't be opened"
- App not notarized
- Run: `xattr -cr Browski.app` to remove quarantine (testing only)

### Windows: "Windows protected your PC"
- Binary not signed
- Sign with EV certificate for instant SmartScreen reputation

### macOS: Notarization fails
- Check entitlements
- Ensure all binaries in bundle are signed
- Use hardened runtime

## Cost Summary

| Item | Cost | Notes |
|------|------|-------|
| Apple Developer Program | $99/year | Required for macOS |
| Windows Code Signing | $200-400/year | Optional but recommended |
| Linux | Free | No signing required |

## Resources

- [Apple Code Signing Guide](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Windows Signing Guide](https://learn.microsoft.com/en-us/windows/win32/seccrypto/signing-code)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
