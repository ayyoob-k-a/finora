# Start Finora API on Port 8081 - SIMPLE VERSION

Write-Host "🚀 Starting Finora API on Port 8081" -ForegroundColor Green

# Clear and set environment variables
$env:PORT = '8081'
$env:DB_USER = 'finora_user'
$env:DB_PASSWORD = 'finora123'
$env:DB_NAME = 'finora_db'
$env:JWT_SECRET = 'finora-port-8081-fixed'
$env:GIN_MODE = 'debug'

Write-Host "✅ Environment variables set:" -ForegroundColor Cyan
Write-Host "  PORT = $($env:PORT)" -ForegroundColor Yellow
Write-Host "  DB_USER = $($env:DB_USER)" -ForegroundColor White

Write-Host "`n🔄 Starting server..." -ForegroundColor Blue
go run main.go
