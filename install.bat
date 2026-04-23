@echo off
setlocal enabledelayedexpansion

REM Check if running as administrator
net session >nul 2>&1
if %ERRORLEVEL% neq 0 (
    color 0E
    cls
    echo.
    echo  ========================================================
    echo  .                                                      .
    echo  .     Requesting Administrator Privileges...          .
    echo  .                                                      .
    echo  ========================================================
    echo.
    echo  This installer needs administrator rights to copy
    echo  GitBro to the Go bin directory.
    echo.
    echo  A new window will open with administrator privileges.
    echo.
    timeout /t 2 /nobreak
    
    REM Re-run as administrator
    powershell -NoProfile -Command "Start-Process '%~0' -Verb RunAs"
    exit /b 0
)

color 0A
cls
setlocal enabledelayedexpansion

title GitBro Simple Installer

echo.
echo  ========================================================
echo  .                                                      .
echo  .         GitBro Installation                         .
echo  .                                                      .
echo  ========================================================
echo.
echo  This will copy GitBro to your Go bin directory.
echo.
pause

color 0B
cls
echo.
echo  Installing GitBro...
echo.

REM Go bin directory is already in PATH
set "INSTALL_DIR=C:\Program Files\Go\bin"

if not exist "%INSTALL_DIR%" (
    color 0C
    echo.
    echo  ERROR: Go bin directory not found at:
    echo  %INSTALL_DIR%
    echo.
    echo  Please install Go first from: https://golang.org
    echo.
    pause
    exit /b 1
)

echo  [*] Copying gitbro.exe...
copy /Y "%~dp0gitbro.exe" "%INSTALL_DIR%\gitbro.exe" >nul 2>&1

if %ERRORLEVEL% equ 0 (
    echo  [OK] Installation successful!
    color 0A
    cls
    echo.
    echo.
    echo  ========================================================
    echo  .                                                      .
    echo  .         Installation Complete!                      .
    echo  .                                                      .
    echo  ========================================================
    echo.
    echo  GitBro is now installed and ready to use!
    echo.
    echo  Location: %INSTALL_DIR%\gitbro.exe
    echo.
    echo  IMPORTANT - Next Steps:
    echo  1. Close all Command Prompt and PowerShell windows
    echo  2. Open a FRESH Command Prompt
    echo  3. Type: gitbro (from any folder)
    echo.
) else (
    color 0C
    echo.
    echo  ERROR: Could not copy file (permission denied).
    echo.
    echo  Solution: Right-click install-simple.bat and select
    echo  "Run as administrator"
    echo.
    pause
    exit /b 1
)

pause
color 07
exit /b 0
