# GitBro for Windows

This folder contains the Windows GitBro package.

## Files

- `gitbro.exe` - GitBro executable
- `install.bat` - installs GitBro system-wide
- `uninstall.bat` - removes GitBro

## Install

1. Extract the package folder first.
2. Keep `gitbro.exe` and `install.bat` in the same folder.
3. Double-click `install.bat`.
4. Accept the Administrator permission prompt.
5. Close all Command Prompt and PowerShell windows.
6. Open a fresh terminal and run:

```bat
gitbro
```

The installer copies GitBro to:

```bat
C:\Program Files\GitBro\bin\gitbro.exe
```

It also adds this folder to the system PATH:

```bat
C:\Program Files\GitBro\bin
```

Go does not need to be installed.

## Set Gemini API Key

GitBro works without an API key using rule-based suggestions, but AI suggestions need `GEMINI_API_KEY`.

Set the key from Command Prompt:

```bat
setx GEMINI_API_KEY "your-gemini-api-key-here"
```

Then close and reopen Command Prompt.

Verify:

```bat
echo %GEMINI_API_KEY%
```

## Use

Inside any Git repository:

```bat
git add .
gitbro
```

## Uninstall

Double-click:

```bat
uninstall.bat
```

Accept the Administrator permission prompt. After uninstall, open a fresh Command Prompt and verify:

```bat
where gitbro
```