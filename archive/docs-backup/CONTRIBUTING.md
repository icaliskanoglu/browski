# Contributing to Browski

Thank you for your interest in contributing to Browski! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment, trolling, or discriminatory language
- Publishing others' private information
- Any conduct that could reasonably be considered inappropriate

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Node.js 22 or higher
- Wails v3 (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- Git

### Platform-Specific Requirements

**macOS:**
- Xcode Command Line Tools

**Linux:**
```bash
sudo apt-get install -y \
  libgtk-4-dev \
  libwebkitgtk-6.0-dev \
  build-essential \
  pkg-config
```

**Windows:**
- Visual Studio Build Tools or MinGW-w64

## Development Setup

1. **Fork the repository**
   ```bash
   # Click "Fork" on GitHub, then:
   git clone https://github.com/YOUR_USERNAME/browski.git
   cd browski
   ```

2. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/icaliskanoglu/browski.git
   ```

3. **Install dependencies**
   ```bash
   cd frontend
   npm install
   cd ..
   go mod download
   ```

4. **Run in development mode**
   ```bash
   wails3 dev
   ```

5. **Build the application**
   ```bash
   wails3 build
   ```

## How to Contribute

### Reporting Bugs

1. **Check existing issues** to avoid duplicates
2. **Use the bug report template** when creating an issue
3. **Include**:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshots if applicable
   - OS and version information
   - Browski version

### Suggesting Features

1. **Check existing issues** for similar suggestions
2. **Use the feature request template**
3. **Explain**:
   - The problem you're trying to solve
   - Your proposed solution
   - Alternative solutions you've considered
   - Why this would be useful to others

### Code Contributions

1. **Find an issue to work on** or create one
2. **Comment on the issue** to let others know you're working on it
3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

## Pull Request Process

### Before Submitting

1. **Sync with upstream**
   ```bash
   git fetch upstream
   git rebase upstream/redesign
   ```

2. **Test your changes**
   ```bash
   # Run tests
   go test ./...

   # Build and test manually
   wails3 build
   ./bin/browski-v3 https://github.com
   ```

3. **Format your code**
   ```bash
   # Go code
   go fmt ./...
   gofmt -s -w .

   # Frontend code
   cd frontend
   npm run format  # if available
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add support for Firefox containers"
   ```

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(browsers): add Firefox container support
fix(safari): resolve icon extraction on macOS 15
docs(readme): update installation instructions
refactor(preferences): simplify save logic
```

### Submitting the PR

1. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request on GitHub**
   - Use a clear title describing the change
   - Fill out the PR template completely
   - Link related issues
   - Add screenshots/videos if UI changed

3. **Respond to feedback**
   - Address review comments promptly
   - Make requested changes
   - Push updates (they'll automatically appear in the PR)

### PR Review Process

- Maintainers will review within 1-3 days
- CI must pass (tests, builds)
- At least one approval required
- May request changes or ask questions
- Once approved, maintainers will merge

## Coding Standards

### Go Code

- Follow standard Go conventions
- Use `gofmt` and `golangci-lint`
- Write meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

```go
// Good
func DetectBrowsers() []Browser {
    // Implementation
}

// Bad
func db() []Browser {
    // Implementation
}
```

### JavaScript/React

- Use functional components with hooks
- Follow existing code style
- Use meaningful component/variable names
- Add JSDoc comments for complex functions

```jsx
// Good
function BrowserCard({ browser, onSelect }) {
    // Implementation
}

// Bad
function BC({ b, os }) {
    // Implementation
}
```

### File Organization

```
browski/
├── main.go              # Application entry point
├── browserservice.go    # Main service layer
├── internal/
│   ├── browsers/       # Browser detection logic
│   ├── preferences/    # Preferences management
│   └── updater/        # Auto-update system
└── frontend/
    └── src/
        ├── components/ # React components
        └── ...
```

## Testing

### Go Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/browsers
```

### Manual Testing Checklist

- [ ] App launches successfully
- [ ] System tray icon appears
- [ ] Browser detection works
- [ ] Profile selection works
- [ ] Preferences save correctly
- [ ] Window closes on mouse leave
- [ ] Single instance enforcement works
- [ ] Updates check in background

### Platform-Specific Testing

Test on all supported platforms when possible:
- macOS (Intel + Apple Silicon)
- Windows 10/11
- Ubuntu/Debian Linux

## Documentation

### Code Documentation

- Add comments for exported functions
- Explain complex logic
- Document parameters and return values
- Include usage examples

```go
// ExtractAppIcon extracts an icon from a macOS .app bundle and converts
// it to a base64-encoded PNG data URL.
//
// Parameters:
//   - appPath: Path to the .app bundle (e.g., "/Applications/Safari.app")
//
// Returns:
//   - data URL string in format "data:image/png;base64,..."
//   - error if extraction fails
func ExtractAppIcon(appPath string) (string, error) {
    // Implementation
}
```

### User Documentation

Update README.md if your changes affect:
- Installation process
- Usage instructions
- Supported browsers/features
- Configuration options

## Recognition

Contributors will be:
- Listed in GitHub contributors
- Mentioned in release notes for significant contributions
- Thanked in the README (optional)

## Questions?

- Open a [Discussion](https://github.com/icaliskanoglu/browski/discussions)
- Join our community chat (if available)
- Email: (add if you want)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for making Browski better! 🎉
