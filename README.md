# GitBro - Cross-Platform Git Commit Message Generator

A professional CLI tool for generating intelligent git commit messages using AI or rule-based suggestions.

## Features

- AI-powered suggestions with Google Gemini API
- Conventional Commit style messages
- Color-coded CLI output
- Cross-platform support for Windows, Linux, and macOS
- Rule-based fallback when no API key is configured

## Project Layout

```text
.
|-- main.go
|-- go.mod
|-- builds/
|   |-- windows/
|   |   |-- gitbro.exe
|   |   |-- install.bat
|   |   `-- uninstall.bat
|   |-- linux/
|   |   |-- gitbro
|   |   |-- install.sh
|   |   `-- uninstall.sh
|   `-- macos/
|       |-- gitbro-intel
|       |-- gitbro-apple-silicon
|       |-- install.sh
|       `-- uninstall.sh
|-- windows/
|   |-- gitbro.exe
|   |-- install.bat
|   `-- uninstall.bat
|-- linux/
|   |-- build-linux.sh
|   |-- install-linux.sh
|   `-- uninstall-linux.sh
`-- macos/
    |-- build-macos.sh
    |-- install-macos.sh
    `-- uninstall-macos.sh
```

## Share Build Packages

If you only want to give developers ready-to-run files, share the matching platform folder from `builds/`:

- Windows: `builds/windows/`
- Linux: `builds/linux/`
- macOS: `builds/macos/`

Each folder includes the binary plus installer/uninstaller files. The installer copies `gitbro` to a PATH location so users can run `gitbro` from any terminal folder. Windows installs to `C:\Program Files\GitBro\bin` and does not require Go to be installed.

You do not need to share source files, `go.mod`, or `go.sum` for normal usage.

## Quick Start

### Windows

```bat
cd windows
install.bat
```

Close and reopen your terminal, then run:

```bash
gitbro
```

### Linux

```bash
bash linux/build-linux.sh
bash linux/install-linux.sh
```

If prompted, add `~/.local/bin` to your PATH:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Then run:

```bash
gitbro
```

### macOS

```bash
bash macos/build-macos.sh
bash macos/install-macos.sh
```

Then run:

```bash
gitbro
```

## Setup - Get Free API Key

To use AI-powered suggestions, get a free Gemini API key:

1. Visit [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click "Create API Key"
3. Set environment variable:

**Windows:**

```cmd
setx GEMINI_API_KEY "your-key-here"
```

**Linux/macOS:**

```bash
export GEMINI_API_KEY="your-key-here"
```

Without an API key, GitBro uses built-in rule-based suggestions.

## Usage

```bash
# 1. Stage your changes
git add .

# 2. Run GitBro
gitbro

# 3. Select a commit message
```

## Building from Source

### Prerequisites

- Go 1.19+ installed
- Git

### Build Commands

**Windows:**

```bash
go build -buildvcs=false -o windows/gitbro.exe .
```

**Linux:**

```bash
GOOS=linux GOARCH=amd64 go build -buildvcs=false -o linux/gitbro-linux .
```

**macOS Intel:**

```bash
GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o macos/gitbro-macos-intel .
```

**macOS Apple Silicon:**

```bash
GOOS=darwin GOARCH=arm64 go build -buildvcs=false -o macos/gitbro-macos-arm64 .
```

Or use provided build scripts:

```bash
bash linux/build-linux.sh
bash macos/build-macos.sh
```

## Distribution

### Windows Package

- `windows/gitbro.exe`
- `windows/install.bat`
- `windows/uninstall.bat`
- `README.md`

### Linux Package

- `linux/gitbro-linux`
- `linux/install-linux.sh`
- `linux/uninstall-linux.sh`
- `README.md`

### macOS Package

- `macos/gitbro-macos-intel`
- `macos/gitbro-macos-arm64`
- `macos/install-macos.sh`
- `macos/uninstall-macos.sh`
- `README.md`

## Uninstallation

**Windows:**

```bat
windows\uninstall.bat
```

**Linux:**

```bash
bash linux/uninstall-linux.sh
```

**macOS:**

```bash
bash macos/uninstall-macos.sh
```

## Troubleshooting

### "gitbro: command not found"

- Windows: close all terminal windows and open a fresh one
- Linux: run `source ~/.bashrc` and ensure `~/.local/bin` is in PATH
- macOS: close terminal and open a new one

### API Key Issues

- Ensure `GEMINI_API_KEY` is set correctly
- Restart your terminal after setting the key
- Without an API key, GitBro uses built-in rule-based suggestions

## License

MIT

## Author

Built for developers everywhere
