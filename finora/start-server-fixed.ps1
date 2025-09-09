# Finora API Server - Fixed Configuration Script
# This script sets correct environment variables and starts the server

Write-Host "🚀 Starting Finora API Server with Fixed Configuration..." -ForegroundColor Green

# ===== SERVER CONFIGURATION (FIXED) =====
$env:PORT = '8081'              # FIXED: Was incorrectly set to 5432 (PostgreSQL port)
$env:GIN_MODE = 'debug'         # Use 'release' for production

# ===== DATABASE CONFIGURATION (FIXED) =====
$env:DB_HOST = 'localhost'
$env:DB_PORT = '5432'
$env:DB_USER = 'finora_user'    # FIXED: Was incorrectly set to 'postgres' 
$env:DB_PASSWORD = 'finora123'  # FIXED: Match the user we created
$env:DB_NAME = 'finora_db'
$env:DB_SSLMODE = 'disable'

# ===== SECURITY & JWT =====
$env:JWT_SECRET = 'development-secret-key-finora-2024-very-long-and-secure'
$env:JWT_EXPIRY = '24h'

# ===== EMAIL & SMS CONFIGURATION =====
$env:SMTP_HOST = 'smtp.gmail.com'
$env:SMTP_PORT = '587'
$env:SMTP_USERNAME = 'your-email@gmail.com'
$env:SMTP_PASSWORD = 'your-app-password'

$env:TWILIO_ACCOUNT_SID = 'your-twilio-sid'
$env:TWILIO_AUTH_TOKEN = 'your-twilio-token'
$env:TWILIO_PHONE_NUMBER = '+1234567890'

# ===== REDIS CONFIGURATION =====
$env:REDIS_HOST = 'localhost'
$env:REDIS_PORT = '6379'
$env:REDIS_PASSWORD = ''

Write-Host "📋 Configuration Summary:" -ForegroundColor Cyan
Write-Host "  🌐 Server Port: $($env:PORT)" -ForegroundColor White
Write-Host "  🗄️  Database: $($env:DB_USER)@$($env:DB_HOST):$($env:DB_PORT)/$($env:DB_NAME)" -ForegroundColor White
Write-Host "  🔐 JWT Secret: Configured ✅" -ForegroundColor White
Write-Host "  ⚙️  Mode: $($env:GIN_MODE)" -ForegroundColor White

Write-Host "`n🔧 Starting server..." -ForegroundColor Yellow

try {
    # Change to the correct directory
    Set-Location -Path $PSScriptRoot
    
    # Check if Go is available
    if (!(Get-Command "go" -ErrorAction SilentlyContinue)) {
        Write-Host "❌ Go is not installed or not in PATH" -ForegroundColor Red
        Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
        exit 1
    }
    
    # Tidy dependencies
    Write-Host "📦 Tidying Go modules..." -ForegroundColor Blue
    go mod tidy
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to tidy Go modules" -ForegroundColor Red
        exit 1
    }
    
    # Start the server
    Write-Host "🚀 Starting Finora API on http://localhost:$($env:PORT)" -ForegroundColor Green
    Write-Host "📚 API Documentation: See COMPLETE_API_DOCUMENTATION.md" -ForegroundColor Cyan
    Write-Host "🧪 Test with Postman: Import Finora_API_v2.postman_collection.json" -ForegroundColor Cyan
    Write-Host "⏹️  Press Ctrl+C to stop the server" -ForegroundColor Yellow
    Write-Host "`n" -ForegroundColor White
    
    go run main.go
    
} catch {
    Write-Host "❌ Error starting server: $_" -ForegroundColor Red
    exit 1
}

Write-Host "`n🛑 Server stopped." -ForegroundColor Yellow
