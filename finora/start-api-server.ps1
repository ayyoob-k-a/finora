# ================================
# Finora API Server Startup Script  
# ================================

Write-Host "🚀 STARTING FINORA API SERVER" -ForegroundColor Green
Write-Host "==============================" -ForegroundColor Green
Write-Host ""

# Function to check if process is running on port
function Test-Port($port) {
    $connection = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
    return $connection -ne $null
}

# Step 1: Set correct environment variables
Write-Host "Step 1: Setting environment variables..." -ForegroundColor Cyan

$env:PORT = '8081'                                          # API server port (NOT 5432!)
$env:DB_HOST = 'localhost'
$env:DB_PORT = '5432'                                       # PostgreSQL port  
$env:DB_USER = 'finora_user'                              # Database user (NOT postgres!)
$env:DB_PASSWORD = 'finora_password'
$env:DB_NAME = 'finora_db'
$env:DB_SSLMODE = 'disable'
$env:JWT_SECRET = 'development-secret-key-12345678'

Write-Host "✅ Environment variables set:" -ForegroundColor Green
Write-Host "   API Port: $env:PORT" -ForegroundColor Yellow
Write-Host "   DB User: $env:DB_USER" -ForegroundColor Yellow  
Write-Host "   DB Name: $env:DB_NAME" -ForegroundColor Yellow

# Step 2: Check PostgreSQL status
Write-Host ""
Write-Host "Step 2: Checking PostgreSQL status..." -ForegroundColor Cyan

$pgCtlPath = "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\pgsql\bin\pg_ctl.exe"
$pgDataPath = "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\data"
$pgLogPath = "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\logfile"

if (Test-Path $pgCtlPath) {
    try {
        $pgStatus = & $pgCtlPath -D $pgDataPath status 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ PostgreSQL is running" -ForegroundColor Green
        } else {
            Write-Host "⚠️  PostgreSQL not running. Starting..." -ForegroundColor Yellow
            & $pgCtlPath -D $pgDataPath -l $pgLogPath start
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ PostgreSQL started successfully" -ForegroundColor Green
            } else {
                Write-Host "❌ Failed to start PostgreSQL" -ForegroundColor Red
                Write-Host "⚠️  API will run in API-only mode" -ForegroundColor Yellow
            }
        }
    } catch {
        Write-Host "❌ Error checking PostgreSQL: $_" -ForegroundColor Red
        Write-Host "⚠️  API will run in API-only mode" -ForegroundColor Yellow
    }
} else {
    Write-Host "❌ PostgreSQL not found at expected path" -ForegroundColor Red
    Write-Host "⚠️  API will run in API-only mode" -ForegroundColor Yellow
}

# Step 3: Check if API port is available
Write-Host ""
Write-Host "Step 3: Checking API port availability..." -ForegroundColor Cyan

if (Test-Port 8081) {
    Write-Host "⚠️  Port 8081 is already in use" -ForegroundColor Yellow
    Write-Host "Attempting to stop existing processes..." -ForegroundColor Yellow
    
    # Stop any existing finora processes
    Get-Process -Name "*finora*" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
    
    Start-Sleep -Seconds 2
    
    if (Test-Port 8081) {
        Write-Host "❌ Port 8081 still in use. Please check manually." -ForegroundColor Red
        exit 1
    } else {
        Write-Host "✅ Port 8081 is now available" -ForegroundColor Green  
    }
} else {
    Write-Host "✅ Port 8081 is available" -ForegroundColor Green
}

# Step 4: Build and start API server
Write-Host ""
Write-Host "Step 4: Starting Finora API server..." -ForegroundColor Cyan

if (!(Test-Path "main.go")) {
    Write-Host "❌ main.go not found. Please run from finora directory." -ForegroundColor Red
    exit 1
}

Write-Host "Building application..." -ForegroundColor Yellow
$buildResult = go build -o finora.exe . 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed:" -ForegroundColor Red
    Write-Host $buildResult -ForegroundColor Red
    exit 1
}

Write-Host "✅ Build successful. Starting server..." -ForegroundColor Green
Write-Host ""

# Start the server
Write-Host "🚀 FINORA API SERVER STARTING..." -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green
Write-Host "API URL: http://localhost:8081" -ForegroundColor Cyan
Write-Host "Health Check: http://localhost:8081/health" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Yellow
Write-Host ""

try {
    .\finora.exe
} catch {
    Write-Host "❌ Server failed to start: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "🏁 Server stopped." -ForegroundColor Yellow
