# GitBro macOS Scripts

This folder contains macOS build/install helper scripts for GitBro.

## Files

- `build-macos.sh` - builds Intel and Apple Silicon executables from source
- `install-macos.sh` - installs the correct binary automatically
- `uninstall-macos.sh` - removes GitBro

## Build

From the repository root:

```bash
bash macos/build-macos.sh
```

This creates:

```bash
macos/gitbro-macos-intel
macos/gitbro-macos-arm64
```

## Install

After building, run:

```bash
bash macos/install-macos.sh
```

The installer detects your Mac architecture and copies GitBro to:

```bash
/usr/local/bin/gitbro
```

If macOS asks for a password, enter your Mac login password.

Close and reopen Terminal, then verify:

```bash
gitbro
```

## Set Gemini API Key

GitBro works without an API key using rule-based suggestions, but AI suggestions need `GEMINI_API_KEY`.

For the current terminal session:

```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
```

For Zsh permanently, which is the default shell on modern macOS:

```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

For Bash permanently:

```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.bash_profile
source ~/.bash_profile
```

Verify:

```bash
echo "$GEMINI_API_KEY"
```

## Use

Inside any Git repository:

```bash
git add .
gitbro
```

## Uninstall

```bash
bash macos/uninstall-macos.sh
```