@echo off
setlocal

title GitBro Installer
set "LOG=%~dp0install.log"
set "INSTALL_DIR=C:\Program Files\GitBro\bin"
set "TARGET=%INSTALL_DIR%\gitbro.exe"

echo GitBro installer started at %DATE% %TIME% > "%LOG%"

REM Check if running as administrator.
net session >nul 2>&1
if errorlevel 1 (
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
    echo  GitBro to Program Files and add it to system PATH.
    echo.
    echo  A new window will open with administrator privileges.
    echo.
    echo Requesting admin privileges... >> "%LOG%"
    timeout /t 2 /nobreak >nul
    powershell -NoProfile -ExecutionPolicy Bypass -Command "Start-Process -FilePath '%~f0' -Verb RunAs -WorkingDirectory '%~dp0'"
    exit /b 0
)

color 0A
echo.
echo  ========================================================
echo  .                                                      .
echo  .         GitBro Installation                         .
echo  .                                                      .
echo  ========================================================
echo.
echo  This will install GitBro to:
echo  %INSTALL_DIR%
echo.
echo  Go does NOT need to be installed on this computer.
echo.
pause

echo.
echo  Installing GitBro...
echo.
echo Running as administrator. >> "%LOG%"

if not exist "%~dp0gitbro.exe" (
    color 0C
    echo.
    echo  ERROR: gitbro.exe not found next to install.bat
    echo.
    echo  Please extract the zip first, then run install.bat again.
    echo.
    echo Missing binary: %~dp0gitbro.exe >> "%LOG%"
    pause
    exit /b 1
)

if not exist "%INSTALL_DIR%" (
    echo  [*] Creating install directory...
    mkdir "%INSTALL_DIR%" >> "%LOG%" 2>&1
    if errorlevel 1 goto create_failed
)

echo  [*] Copying gitbro.exe...
copy /Y "%~dp0gitbro.exe" "%TARGET%" >> "%LOG%" 2>&1
if errorlevel 1 goto copy_failed
if not exist "%TARGET%" goto copy_failed

echo  [*] Adding GitBro to system PATH...
powershell -NoProfile -ExecutionPolicy Bypass -Command "$dir='C:\Program Files\GitBro\bin'; $path=[Environment]::GetEnvironmentVariable('Path','Machine'); if ([string]::IsNullOrWhiteSpace($path)) { $new=$dir } elseif (($path -split ';') -contains $dir) { $new=$path } else { $new=$path.TrimEnd(';') + ';' + $dir }; [Environment]::SetEnvironmentVariable('Path',$new,'Machine')" >> "%LOG%" 2>&1
if errorlevel 1 goto path_failed

echo.
echo  ========================================================
echo  .                                                      .
echo  .         Installation Complete!                      .
echo  .                                                      .
echo  ========================================================
echo.
echo  GitBro is now installed and ready to use!
echo.
echo  Location: %TARGET%
echo.
echo  IMPORTANT - Next Steps:
echo  1. Close all Command Prompt and PowerShell windows
echo  2. Open a FRESH Command Prompt
echo  3. Type: gitbro from any folder
echo.
echo Installation complete. >> "%LOG%"
pause
color 07
exit /b 0

:create_failed
color 0C
echo.
echo  ERROR: Could not create install directory:
echo  %INSTALL_DIR%
echo.
echo  Check install.log for details:
echo  %LOG%
echo.
pause
exit /b 1

:copy_failed
color 0C
echo.
echo  ERROR: Could not copy gitbro.exe to:
echo  %TARGET%
echo.
echo  Check install.log for details:
echo  %LOG%
echo.
pause
exit /b 1

:path_failed
color 0E
echo.
echo  WARNING: GitBro was copied, but PATH update failed.
echo.
echo  GitBro location:
echo  %TARGET%
echo.
echo  Add this folder manually to system PATH:
echo  %INSTALL_DIR%
echo.
echo  Check install.log for details:
echo  %LOG%
echo.
pause
exit /b 1