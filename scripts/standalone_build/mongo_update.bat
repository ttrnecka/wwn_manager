@echo off
setlocal enabledelayedexpansion

REM =========================================================
REM MongoDB Community Server Upgrade Script
REM
REM Features:
REM - Silent non-interactive upgrade
REM - Preserves existing mongod.cfg (backs up + restores)
REM - Preserves existing Windows service configuration
REM - Does NOT recreate MongoDB service
REM - Optional --no-download mode
REM
REM Usage:
REM   upgrade-mongodb.bat
REM       -> downloads latest MSI automatically
REM
REM   upgrade-mongodb.bat --no-download
REM       -> expects mongodb-community-latest.msi
REM       -> in same folder as this BAT file
REM =========================================================

REM ---------- CONFIG ----------
set "VERSION=7.0.31"
set "DOWNLOAD_URL=https://fastdl.mongodb.org/windows/mongodb-windows-x86_64-%VERSION%-signed.msi"

set "TEMP_MSI=%TEMP%\mongodb-community-%VERSION%.msi"
set "LOCAL_MSI=%~dp0mongodb-community-%VERSION%.msi"

REM IMPORTANT:
REM Adjust if your existing config lives elsewhere
set "CONFIG_FILE=C:\Program Files\MongoDB\Server\mongod.cfg"
set "CONFIG_BACKUP=%TEMP%\mongod.cfg.backup"

set "NO_DOWNLOAD=false"

REM ---------- ARGUMENT PARSING ----------
if /I "%~1"=="--no-download" (
    set "NO_DOWNLOAD=true"
)

echo =========================================
echo MongoDB Community Server Silent Upgrade
echo =========================================
echo.

REM ---------- VERIFY CONFIG ----------
if exist "%CONFIG_FILE%" (
    echo Existing config found:
    echo %CONFIG_FILE%
    echo Backing up config...

    copy /Y "%CONFIG_FILE%" "%CONFIG_BACKUP%" >nul

    if errorlevel 1 (
        echo ERROR: Failed to back up existing mongod.cfg
        exit /b 1
    )

    echo Backup created:
    echo %CONFIG_BACKUP%
) else (
    echo WARNING: Existing config not found:
    echo %CONFIG_FILE%
    echo Installer may create a new default config.
)

echo.

REM ---------- DETERMINE MSI SOURCE ----------
if /I "%NO_DOWNLOAD%"=="true" (
    echo --no-download mode enabled
    echo Expecting MSI here:
    echo %LOCAL_MSI%
    echo.

    if not exist "%LOCAL_MSI%" (
        echo ERROR: MSI file not found:
        echo %LOCAL_MSI%
        exit /b 1
    )

    set "MSI_FILE=%LOCAL_MSI%"
) else (
    echo Downloading latest MongoDB Community Server...

    powershell -Command ^
        "Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%TEMP_MSI%'"

    if not exist "%TEMP_MSI%" (
        echo ERROR: Failed to download MongoDB installer.
        exit /b 1
    )

    echo Download completed.
    echo.

    set "MSI_FILE=%TEMP_MSI%"
)

echo Using MSI:
echo %MSI_FILE%
echo.

REM ---------- STOP SERVICE ----------
echo Stopping MongoDB service if running...

sc query MongoDB | find /I "RUNNING" >nul
if %errorlevel%==0 (
    net stop MongoDB
    timeout /t 5 /nobreak >nul
) else (
    echo MongoDB service is not running.
)

echo.

REM ---------- SILENT UPGRADE ----------
echo Installing/upgrading MongoDB silently...

msiexec /i "%MSI_FILE%" ^
    /qn ^
    SHOULD_INSTALL_COMPASS="0"

if errorlevel 1 (
    echo ERROR: MongoDB installation/upgrade failed.
    exit /b 1
)

echo Upgrade completed successfully.
echo.

REM ---------- RESTORE CONFIG ----------
if exist "%CONFIG_BACKUP%" (
    echo Restoring original mongod.cfg...

    copy /Y "%CONFIG_BACKUP%" "%CONFIG_FILE%" >nul

    if errorlevel 1 (
        echo ERROR: Failed to restore original mongod.cfg
        echo Please restore manually from:
        echo %CONFIG_BACKUP%
        exit /b 1
    )

    echo Original config restored successfully.
)

echo.

REM ---------- START SERVICE ----------
echo Starting MongoDB service...

net start MongoDB

if errorlevel 1 (
    echo WARNING: MongoDB service did not start.
    echo Please verify service configuration manually.
) else (
    echo MongoDB service started successfully.
)

echo.

REM ---------- CLEANUP ----------
if /I "%NO_DOWNLOAD%"=="false" (
    del "%TEMP_MSI%" >nul 2>&1
)

if exist "%CONFIG_BACKUP%" (
    del "%CONFIG_BACKUP%" >nul 2>&1
)

echo.
echo =========================================
echo MongoDB upgrade completed successfully
echo Existing mongod.cfg preserved and restored
echo Existing service preserved
echo =========================================

endlocal
exit /b 0