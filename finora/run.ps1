# Set environment variables
$env:JWT_SECRET="development-secret-key-12345678"
$env:DB_HOST="localhost"
$env:DB_USER="postgres" 
$env:DB_PASSWORD="password"
$env:DB_NAME="finora_db"
$env:DB_PORT="5432"
$env:PORT="8080"

# Run the application
Write-Host "Starting Finora API..."
go run main.go
