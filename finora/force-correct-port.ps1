# FORCE CORRECT PORT - Finora API Server
# This script FORCES the server to run on port 8081

Write-Host "🔧 FORCING CORRECT PORT CONFIGURATION" -ForegroundColor Red

# STEP 1: Clear ALL environment variables that might interfere
Write-Host "`n🧹 Clearing conflicting environment variables..." -ForegroundColor Yellow
Remove-Item Env:PORT -ErrorAction SilentlyContinue
Remove-Item Env:DB_* -ErrorAction SilentlyContinue
Remove-Item Env:GIN_MODE -ErrorAction SilentlyContinue
Remove-Item Env:JWT_SECRET -ErrorAction SilentlyContinue

# STEP 2: FORCE SET the correct environment variables
Write-Host "🎯 FORCING correct configuration..." -ForegroundColor Green
$env:PORT = '8081'              # FORCE: API server port
$env:GIN_MODE = 'debug'         
$env:DB_HOST = 'localhost'
$env:DB_PORT = '5432'           # Database port (different from API port)
$env:DB_USER = 'finora_user'    
$env:DB_PASSWORD = 'finora123'  
$env:DB_NAME = 'finora_db'
$env:DB_SSLMODE = 'disable'
$env:JWT_SECRET = 'finora-force-correct-port-2024'

# STEP 3: Verify the environment variables
Write-Host "`n✅ FORCED Configuration Verification:" -ForegroundColor Cyan
Write-Host "  🌐 API Server Port: $($env:PORT) (MUST be 8081)" -ForegroundColor Green
Write-Host "  🗄️  Database Port: $($env:DB_PORT) (5432 for PostgreSQL)" -ForegroundColor White
Write-Host "  👤 Database User: $($env:DB_USER)" -ForegroundColor White
Write-Host "  📊 Database Name: $($env:DB_NAME)" -ForegroundColor White

# STEP 4: Change to correct directory
if (!(Test-Path "main.go")) {
    Write-Host "❌ main.go not found in current directory" -ForegroundColor Red
    Write-Host "Current directory: $(Get-Location)" -ForegroundColor Yellow
    exit 1
}

Write-Host "`n🎯 EXPECTED RESULT:" -ForegroundColor Yellow
Write-Host "  ✅ Server WILL start on http://localhost:8081" -ForegroundColor Green
Write-Host "  ✅ You WILL be able to connect to port 8081" -ForegroundColor Green
Write-Host "  ✅ No more connection refused errors" -ForegroundColor Green

Write-Host "`n🚀 STARTING SERVER WITH FORCED PORT 8081..." -ForegroundColor Green
Write-Host "=" * 60 -ForegroundColor Cyan

try {
    # Force compilation to ensure latest changes
    go build main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Compilation failed" -ForegroundColor Red
        exit 1
    }
    
    # Start with forced environment
    go run main.go
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
    exit 1
}

Write-Host "`n🎉 Server should now be running on PORT 8081!" -ForegroundColor Green
