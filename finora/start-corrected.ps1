# Finora API Server - CORRECTED Configuration Script
# This script fixes all configuration issues and starts the server properly

Write-Host "🔧 Starting Finora API Server with CORRECTED Configuration" -ForegroundColor Green

# ===== CORRECTED CONFIGURATION =====
# Clear any conflicting environment variables first
Remove-Item Env:PORT -ErrorAction SilentlyContinue
Remove-Item Env:DB_USER -ErrorAction SilentlyContinue
Remove-Item Env:DB_PASSWORD -ErrorAction SilentlyContinue

# Set CORRECTED environment variables
$env:PORT = '8081'              # CORRECTED: Use proper API port (not 5432!)
$env:GIN_MODE = 'debug'         
$env:DB_HOST = 'localhost'
$env:DB_PORT = '5432'           # PostgreSQL port (for database connection)
$env:DB_USER = 'finora_user'    # CORRECTED: Use proper user (not postgres!)
$env:DB_PASSWORD = 'finora123'  # CORRECTED: Use proper password
$env:DB_NAME = 'finora_db'
$env:DB_SSLMODE = 'disable'
$env:JWT_SECRET = 'finora-corrected-secret-key-2024'

Write-Host "`n📋 CORRECTED Configuration:" -ForegroundColor Cyan
Write-Host "  🌐 API Server Port: $($env:PORT) (CORRECTED from 5432)" -ForegroundColor Yellow
Write-Host "  🗄️  Database User: $($env:DB_USER) (CORRECTED from postgres)" -ForegroundColor Yellow
Write-Host "  🔐 Database Password: $($env:DB_PASSWORD) (SET PROPERLY)" -ForegroundColor Yellow
Write-Host "  📊 Database Connection: $($env:DB_HOST):$($env:DB_PORT)/$($env:DB_NAME)" -ForegroundColor White

Write-Host "`n⚠️  Configuration Issues FIXED:" -ForegroundColor Red
Write-Host "  ✅ Server port changed from 5432 to 8081" -ForegroundColor Green
Write-Host "  ✅ Database user changed from 'postgres' to 'finora_user'" -ForegroundColor Green
Write-Host "  ✅ Database password properly set" -ForegroundColor Green
Write-Host "  ✅ All environment variables cleared and reset" -ForegroundColor Green

try {
    # Change to correct directory
    Set-Location -Path $PSScriptRoot
    
    # Check if Go is available
    if (!(Get-Command "go" -ErrorAction SilentlyContinue)) {
        Write-Host "❌ Go is not installed or not in PATH" -ForegroundColor Red
        Write-Host "Please install Go from: https://golang.org/dl/" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "`n🔍 Verifying configuration..." -ForegroundColor Blue
    
    # Quick compile test
    $compileResult = go build -o finora-corrected.exe main.go 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Compilation failed:" -ForegroundColor Red
        Write-Host $compileResult -ForegroundColor Yellow
        exit 1
    } else {
        Write-Host "✅ Compilation successful!" -ForegroundColor Green
    }
    
    Write-Host "`n📦 Ensuring dependencies are up to date..." -ForegroundColor Blue
    go mod tidy
    
    Write-Host "`n🚀 Starting server with CORRECTED configuration..." -ForegroundColor Green
    Write-Host "`n🎯 Expected Behavior:" -ForegroundColor Cyan
    Write-Host "  ✅ Server will start on http://localhost:8081 (NOT 5432)" -ForegroundColor White
    Write-Host "  ✅ Database will try to connect as 'finora_user' (NOT postgres)" -ForegroundColor White
    Write-Host "  ✅ If database is not running, API will work in 'API-only mode'" -ForegroundColor White
    Write-Host "  ✅ All 36 endpoints will be available for testing" -ForegroundColor White
    
    Write-Host "`n📚 Quick Test Guide:" -ForegroundColor Yellow
    Write-Host "  1. Open browser: http://localhost:8081/health" -ForegroundColor White
    Write-Host "  2. Import Postman collection: Finora_API_v3_FIXED.postman_collection.json" -ForegroundColor White
    Write-Host "  3. Test authentication: Send OTP → Verify OTP" -ForegroundColor White
    Write-Host "  4. Test any endpoint with the JWT token" -ForegroundColor White
    
    Write-Host "`n🎉 STARTING SERVER NOW (Press Ctrl+C to stop)..." -ForegroundColor Green
    Write-Host "=" -repeat 60 -ForegroundColor Cyan
    
    # Clean up test executable
    if (Test-Path "finora-corrected.exe") {
        Remove-Item "finora-corrected.exe"
    }
    
    # Start the server with corrected configuration
    go run main.go
    
} catch {
    Write-Host "`n❌ Error starting server: $_" -ForegroundColor Red
    Write-Host "`n🔧 Troubleshooting:" -ForegroundColor Yellow
    Write-Host "  1. Make sure you're in the finora directory" -ForegroundColor White
    Write-Host "  2. Run: go mod tidy" -ForegroundColor White
    Write-Host "  3. Check if PostgreSQL is running (optional)" -ForegroundColor White
    Write-Host "  4. Verify Go installation: go version" -ForegroundColor White
    exit 1
}

Write-Host "`n👋 Server stopped. Configuration has been corrected!" -ForegroundColor Green
