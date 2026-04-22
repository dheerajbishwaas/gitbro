# GitBro

A professional CLI tool for generating conventional commits from staged Git changes.

## Build Instructions

1. **Install Go**: Download and install Go from [https://golang.org/dl/](https://golang.org/dl/).

2. **Build the binary**:
   ```
   go mod init gitbro
   go mod tidy
   go build -ldflags="-s -w" -o gitbro.exe main.go
   ```

3. **Install Inno Setup**: Download and install Inno Setup from [https://jrsoftware.org/isinfo.php](https://jrsoftware.org/isinfo.php).

4. **Compile the installer**:
   - Open `installer.iss` in Inno Setup Compiler.
   - Click "Compile" to generate `gitbro-setup.exe`.

5. **Upload to GitHub Releases**:
   - Create a new release on GitHub.
   - Upload `gitbro-setup.exe` as the release asset.

## Usage

After installation, users can run `gitbro` from any command prompt to generate and commit conventional commits based on staged changes.

## Cross-Compile (if needed)

For cross-compilation to Windows from another OS:
```
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o gitbro.exe main.go
```