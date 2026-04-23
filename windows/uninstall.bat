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
    echo  GitBro from Program Files and system PATH.
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

set "REMOVED=0"
set "INSTALL_DIR=C:\Program Files\GitBro\bin"
set "PATH_KEY=HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment"

call :remove_file "%INSTALL_DIR%\gitbro.exe"
call :remove_file "C:\Program Files\Go\bin\gitbro.exe"
call :remove_file "%USERPROFILE%\GitBro\bin\gitbro.exe"
call :remove_file "%LOCALAPPDATA%\Programs\GitBro\gitbro.exe"
call :remove_file "%LOCALAPPDATA%\Microsoft\WindowsApps\gitbro.exe"

if exist "%INSTALL_DIR%" rmdir "%INSTALL_DIR%" >nul 2>&1
if exist "C:\Program Files\GitBro" rmdir "C:\Program Files\GitBro" >nul 2>&1
if exist "%USERPROFILE%\GitBro\bin" rmdir "%USERPROFILE%\GitBro\bin" >nul 2>&1
if exist "%USERPROFILE%\GitBro" rmdir "%USERPROFILE%\GitBro" >nul 2>&1
if exist "%LOCALAPPDATA%\Programs\GitBro" rmdir "%LOCALAPPDATA%\Programs\GitBro" >nul 2>&1

echo  [*] Removing GitBro from system PATH...
set "MACHINE_PATH="
for /f "tokens=2,*" %%A in ('reg query "%PATH_KEY%" /v Path 2^>nul') do set "MACHINE_PATH=%%B"
if defined MACHINE_PATH (
    call set "NEW_PATH=%%MACHINE_PATH:%INSTALL_DIR%;=%%"
    call set "NEW_PATH=%%NEW_PATH:;%INSTALL_DIR%=%%"
    call set "NEW_PATH=%%NEW_PATH:%INSTALL_DIR%=%%"
    reg add "%PATH_KEY%" /v Path /t REG_EXPAND_SZ /d "%NEW_PATH%" /f >nul 2>&1
)

echo.
echo  Checking PATH for remaining GitBro commands...
echo.
where gitbro >nul 2>&1
if %ERRORLEVEL% equ 0 (
    color 0E
    echo  WARNING: GitBro is still found in PATH at:
    where gitbro
    echo.
    echo  Delete the remaining file manually or remove that folder from PATH.
    echo.
) else (
    color 0A
    echo  [OK] GitBro is no longer found in PATH.
    echo.
)

if "%REMOVED%"=="1" (
    echo  [OK] Uninstallation finished.
) else (
    echo  NOTE: No installed GitBro file was found in common locations.
)

echo.
echo  Next Steps:
echo  1. Close all Command Prompt and PowerShell windows
echo  2. Open a NEW Command Prompt
echo  3. Run: where gitbro
echo.
pause
color 07
exit /b 0

:remove_file
set "TARGET=%~1"
if exist "%TARGET%" (
    echo  [*] Removing %TARGET%
    del /F /Q "%TARGET%" >nul 2>&1
    if exist "%TARGET%" (
        color 0C
        echo  ERROR: Could not remove %TARGET%
    ) else (
        set "REMOVED=1"
        echo  [OK] Removed %TARGET%
    )
)
exit /b 0