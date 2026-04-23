# GitBro Linux Scripts

This folder contains Linux build/install helper scripts for GitBro.

## Files

- `build-linux.sh` - builds the Linux executable from source
- `install-linux.sh` - installs `gitbro-linux` to `~/.local/bin/gitbro`
- `uninstall-linux.sh` - removes GitBro

## Build

From the repository root:

```bash
bash linux/build-linux.sh
```

This creates:

```bash
linux/gitbro-linux
```

## Install

After building, run:

```bash
bash linux/install-linux.sh
```

The installer copies GitBro to:

```bash
~/.local/bin/gitbro
```

Restart your terminal, or run the command shown by the installer, then verify:

```bash
gitbro
```

## Set Gemini API Key

GitBro works without an API key using rule-based suggestions, but AI suggestions need `GEMINI_API_KEY`.

For the current terminal session:

```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
```

For Bash permanently:

```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.bashrc
source ~/.bashrc
```

For Zsh permanently:

```bash
echo 'export GEMINI_API_KEY="your-gemini-api-key-here"' >> ~/.zshrc
source ~/.zshrc
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
bash linux/uninstall-linux.sh
```