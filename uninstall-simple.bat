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
    echo  This uninstaller needs administrator rights to remove
    echo  GitBro from the Go bin directory.
    echo.
    echo  A new window will open with administrator privileges.
    echo.
    timeout /t 2 /nobreak
    
    REM Re-run as administrator
    powershell -NoProfile -Command "Start-Process '%~0' -Verb RunAs"
    exit /b 0
)

color 0C
cls
setlocal enabledelayedexpansion

title GitBro Uninstaller

echo.
echo  ========================================================
echo  .                                                      .
echo  .         GitBro Uninstaller                          .
echo  .                                                      .
echo  ========================================================
echo.
echo  This will remove GitBro from your system.
echo.
set /p confirm="  Are you sure? (Y/N): "
if /i not "%confirm%"=="Y" (
    color 0A
    echo.
    echo  Uninstallation cancelled.
    echo.
    timeout /t 2 /nobreak
    exit /b 0
)

color 0B
cls
echo.
echo  Uninstalling GitBro...
echo.

set "INSTALL_DIR=C:\Program Files\Go\bin"

if not exist "%INSTALL_DIR%\gitbro.exe" (
    color 0C
    echo.
    echo  ERROR: GitBro is not installed.
    echo.
    timeout /t 3 /nobreak
    exit /b 1
)

echo  [*] Removing gitbro.exe...
del /F /Q "%INSTALL_DIR%\gitbro.exe" >nul 2>&1

if %ERRORLEVEL% equ 0 (
    color 0A
    cls
    echo.
    echo.
    echo  ========================================================
    echo  .                                                      .
    echo  .         Uninstallation Complete!                    .
    echo  .                                                      .
    echo  ========================================================
    echo.
    echo  GitBro has been successfully removed.
    echo.
    echo  Next Steps:
    echo  1. Close all Command Prompt and PowerShell windows
    echo  2. Open a NEW Command Prompt
    echo  3. Verify: 'gitbro' command will no longer work
    echo.
) else (
    color 0C
    echo.
    echo  ERROR: Could not delete file (permission denied).
    echo.
    echo  Solution: Right-click uninstall-simple.bat and select
    echo  "Run as administrator"
    echo.
)

pause
color 07
exit /b %ERRORLEVEL%
