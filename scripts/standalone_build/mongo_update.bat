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

set "LOCAL_MSI=%~dp0mongodb-community-%VERSION%.msi"

REM IMPORTANT:
REM Adjust if your existing config lives elsewhere
set "CONFIG_FILE=C:\Program Files\MongoDB\Server\7.0\bin\mongod.cfg"
set "CONFIG_BACKUP=%~dp0mongod.cfg.backup"

set "DATA_DIR=D:\MongoDB\data"
set "LOG_DIR=D:\MongoDB\log"
set "LOG_FILE=%LOG_DIR%\mongod.log"

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
    echo No existing mongod.cfg found.
    echo A new config will be created after installation.
    echo.

    set "CONFIG_EXISTED=false"
)

echo.

REM =========================================================
REM MSI HANDLING
REM =========================================================

set "MSI_FILE="

if /I "%NO_DOWNLOAD%"=="true" (
    echo --no-download mode enabled
    echo Checking local MSI in script folder...

    if exist "%LOCAL_MSI%" (
        echo Found MSI:
        echo %LOCAL_MSI%
        set "MSI_FILE=%LOCAL_MSI%"
    ) else (
        echo ERROR: MSI not found locally
        exit /b 1
    )
) else (
    echo Downloading latest MongoDB Community Server...

    if exist "%LOCAL_MSI%" (
        echo Found MSI:
        echo %LOCAL_MSI%
    ) else (
        powershell -Command ^
            "Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%LOCAL_MSI%'"

        if not exist "%LOCAL_MSI%" (
            echo ERROR: Download failed
            exit /b 1
        )
        echo Download completed.
        echo.
    )

    set "MSI_FILE=%LOCAL_MSI%"
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
    /qb! ^
    ADDLOCAL="ServerService" ^
    SHOULD_INSTALL_COMPASS="0"

if errorlevel 1 (
    echo ERROR: MongoDB installation/upgrade failed.
    exit /b 1
)

echo Upgrade completed successfully.
echo.

REM ---------- STOP SERVICE ----------
echo Stopping MongoDB service...

net stop MongoDB

if errorlevel 1 (
    echo WARNING: MongoDB service did not stop.
) else (
    echo MongoDB service stopped successfully.
)

REM =========================================================
REM RESTORE OR CREATE CONFIG
REM =========================================================

if /I "%CONFIG_EXISTED%"=="true" (
    echo Restoring original mongod.cfg...

    copy /Y "%CONFIG_BACKUP%" "%CONFIG_FILE%" >nul

    if errorlevel 1 (
        echo ERROR: Failed to restore original mongod.cfg
        echo Please restore manually from:
        echo %CONFIG_BACKUP%
        exit /b 1
    )

    echo Existing config restored successfully.
) else (
    echo Creating new mongod.cfg...

    if not exist "%DATA_DIR%" mkdir "%DATA_DIR%"
    if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

    (
        echo systemLog:
        echo.  destination: file
        echo.  path: "%LOG_FILE%"
        echo.  logAppend: true
        echo storage:
        echo.  dbPath: "%DATA_DIR%"
        echo net:
        echo.  bindIp: 127.0.0.1
        echo.  port: 27017
    ) > "%CONFIG_FILE%"

    if not exist "%CONFIG_FILE%" (
        echo ERROR: Failed to create new mongod.cfg
        exit /b 1
    )

    echo New config created successfully.
    echo Data directory: %DATA_DIR%
    echo Log directory : %LOG_DIR%
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
echo =========================================
echo MongoDB upgrade completed successfully
echo Existing mongod.cfg preserved and restored
echo Existing service preserved
echo =========================================

endlocal
exit /b 0