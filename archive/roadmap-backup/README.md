# Roadmap & Planning Documents

This directory contains implementation plans, distribution tasks, and roadmap documentation for Browski development.

## Documents

### [ROADMAP.md](ROADMAP.md)
Long-term feature roadmap organized by release phases:
- Current Status (v1.0)
- Short Term (1-2 releases)
- Medium Term (3-6 months)
- Long Term (6-12 months)
- Community & Documentation
- Technical Debt

### [distribution-plan.md](distribution-plan.md)
Detailed implementation plan for distribution infrastructure:
- Phase 1: Core Functionality (macOS Info.plist, Linux .desktop, binary naming)
- Phase 2: Packaging (DMG, NSIS, nfpm, AppImage)
- Phase 3: Code Signing (macOS notarization, Windows SignPath, Linux GPG)
- Phase 4: Automation (package manager automation, checksums)

Includes task definitions, time estimates, testing procedures, and success criteria.

### [DISTRIBUTION_TASKS.md](DISTRIBUTION_TASKS.md)
Comprehensive checklist based on the distribution specification:
- Already implemented features (✅)
- Pending tasks organized by platform (🚧)
- CI/CD workflow updates
- Acceptance criteria for each platform

### [DISTRIBUTION.md](DISTRIBUTION.md)
Original distribution specification document outlining requirements for:
- macOS distribution (bundle config, signing, notarization, DMG, Homebrew)
- Windows distribution (NSIS installer, registry keys, code signing, winget)
- Linux distribution (.desktop files, package formats, AppImage, Snap, Flatpak)
- GitHub Actions automation
- Single-instance behavior and URL handling

## Usage

These documents are primarily for developers and contributors working on Browski's distribution and packaging infrastructure.

For user-facing installation instructions, see [docs/INSTALLATION.md](../docs/INSTALLATION.md).

For development guidance, see [CLAUDE.md](../CLAUDE.md).
