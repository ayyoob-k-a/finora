# ================================
# Finora Database Setup Script
# ================================

Write-Host "🚀 FINORA DATABASE SETUP" -ForegroundColor Green
Write-Host "=========================" -ForegroundColor Green
Write-Host ""

# Function to test command availability
function Test-Command($command) {
    try {
        Get-Command $command -ErrorAction Stop | Out-Null
        return $true
    } catch {
        return $false
    }
}

# Function to test database connection
function Test-DatabaseConnection($host, $port, $user, $password, $database) {
    Write-Host "🔌 Testing database connection..." -ForegroundColor Yellow
    
    $env:PGPASSWORD = $password
    
    try {
        $result = & psql -h $host -p $port -U $user -d $database -c "SELECT 1;" -t -A 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Database connection successful!" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ Database connection failed: $result" -ForegroundColor Red
            return $false
        }
    } catch {
        Write-Host "❌ Database connection failed: $_" -ForegroundColor Red
        return $false
    } finally {
        Remove-Item env:PGPASSWORD -ErrorAction SilentlyContinue
    }
}

# Function to setup database
function Setup-Database($host, $port, $superuser, $superpass, $database, $user, $password) {
    Write-Host "🔧 Setting up database and user..." -ForegroundColor Yellow
    
    $env:PGPASSWORD = $superpass
    
    $sqlCommands = @"
-- Create database
CREATE DATABASE $database;

-- Create user with password
CREATE USER $user WITH ENCRYPTED PASSWORD '$password';

-- Grant privileges  
GRANT ALL PRIVILEGES ON DATABASE $database TO $user;
GRANT USAGE ON SCHEMA public TO $user;
GRANT CREATE ON SCHEMA public TO $user;
"@

    $dbSqlCommands = @"
-- Grant table privileges
GRANT ALL ON ALL TABLES IN SCHEMA public TO $user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO $user;
GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO $user;

-- Set default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO $user;
"@

    try {
        # Create database and user
        Write-Host "  📝 Creating database and user..." -ForegroundColor Cyan
        $result1 = $sqlCommands | & psql -h $host -p $port -U $superuser -d postgres 2>&1
        
        if ($LASTEXITCODE -ne 0) {
            Write-Host "⚠️  Database/user might already exist. Continuing..." -ForegroundColor Yellow
        }
        
        # Grant permissions on the database
        Write-Host "  🔐 Setting up permissions..." -ForegroundColor Cyan
        $result2 = $dbSqlCommands | & psql -h $host -p $port -U $superuser -d $database 2>&1
        
        Write-Host "✅ Database setup completed!" -ForegroundColor Green
        return $true
        
    } catch {
        Write-Host "❌ Database setup failed: $_" -ForegroundColor Red
        return $false
    } finally {
        Remove-Item env:PGPASSWORD -ErrorAction SilentlyContinue
    }
}

# Main Setup Process
Write-Host "Step 1: Checking PostgreSQL installation..." -ForegroundColor Cyan

if (-not (Test-Command "psql")) {
    Write-Host "❌ PostgreSQL is not installed or not in PATH" -ForegroundColor Red
    Write-Host ""
    Write-Host "📋 TO INSTALL POSTGRESQL:" -ForegroundColor Yellow
    Write-Host "1. Visit: https://www.postgresql.org/download/windows/"
    Write-Host "2. Download PostgreSQL 15+ installer"
    Write-Host "3. Install with default settings"
    Write-Host "4. Remember the 'postgres' user password!"
    Write-Host "5. Add to PATH: C:\Program Files\PostgreSQL\15\bin"
    Write-Host "6. Restart PowerShell and run this script again"
    Write-Host ""
    Write-Host "Press any key to exit..."
    Read-Host
    exit 1
}

Write-Host "✅ PostgreSQL found!" -ForegroundColor Green

# Database configuration
$DB_HOST = "localhost"
$DB_PORT = "5432"  
$DB_NAME = "finora_db"
$DB_USER = "finora_user"
$DB_PASSWORD = "finora_password"

Write-Host ""
Write-Host "Step 2: Database Configuration" -ForegroundColor Cyan
Write-Host "Host: $DB_HOST"
Write-Host "Port: $DB_PORT"
Write-Host "Database: $DB_NAME"
Write-Host "User: $DB_USER"
Write-Host "Password: $DB_PASSWORD"
Write-Host ""

# Test if database already exists
if (Test-DatabaseConnection $DB_HOST $DB_PORT $DB_USER $DB_PASSWORD $DB_NAME) {
    Write-Host "✅ Database is already set up and working!" -ForegroundColor Green
    $useExisting = Read-Host "Continue with existing database? (Y/n)"
    if ($useExisting -eq "" -or $useExisting.ToLower() -eq "y") {
        Write-Host "🎯 Using existing database..." -ForegroundColor Green
    } else {
        Write-Host "❌ Setup cancelled." -ForegroundColor Red
        exit 0
    }
} else {
    Write-Host "Step 3: Setting up new database..." -ForegroundColor Cyan
    Write-Host ""
    
    # Get superuser password
    Write-Host "🔐 Enter password for PostgreSQL superuser 'postgres':" -ForegroundColor Yellow
    $POSTGRES_PASSWORD = Read-Host -AsSecureString
    $POSTGRES_PASSWORD_TEXT = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($POSTGRES_PASSWORD))
    
    # Setup database
    if (-not (Setup-Database $DB_HOST $DB_PORT "postgres" $POSTGRES_PASSWORD_TEXT $DB_NAME $DB_USER $DB_PASSWORD)) {
        Write-Host "❌ Database setup failed. Please check the error messages above." -ForegroundColor Red
        Write-Host "Press any key to exit..."
        Read-Host
        exit 1
    }
    
    # Test the new setup
    if (-not (Test-DatabaseConnection $DB_HOST $DB_PORT $DB_USER $DB_PASSWORD $DB_NAME)) {
        Write-Host "❌ Database setup verification failed." -ForegroundColor Red
        Write-Host "Press any key to exit..."
        Read-Host  
        exit 1
    }
}

Write-Host ""
Write-Host "Step 4: Setting environment variables..." -ForegroundColor Cyan

# Set environment variables
$env:DB_HOST = $DB_HOST
$env:DB_PORT = $DB_PORT
$env:DB_USER = $DB_USER
$env:DB_PASSWORD = $DB_PASSWORD
$env:DB_NAME = $DB_NAME
$env:DB_SSLMODE = "disable"
$env:JWT_SECRET = "development-secret-key-12345678"
$env:PORT = "8081"

Write-Host "✅ Environment variables set!" -ForegroundColor Green

Write-Host ""
Write-Host "Step 5: Testing Finora API with database..." -ForegroundColor Cyan
Write-Host "Starting API server..." -ForegroundColor Yellow

# Start the Go application
try {
    Write-Host "🚀 Starting Finora API..." -ForegroundColor Green
    Write-Host "Press Ctrl+C to stop the server when done testing." -ForegroundColor Yellow
    Write-Host ""
    
    & go run main.go
    
} catch {
    Write-Host "❌ Failed to start API server: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎉 SETUP COMPLETE!" -ForegroundColor Green
Write-Host "================" -ForegroundColor Green
Write-Host ""
Write-Host "📋 CONFIGURATION:" -ForegroundColor Cyan
Write-Host "Database: $DB_NAME"
Write-Host "User: $DB_USER"  
Write-Host "Host: $DB_HOST:$DB_PORT"
Write-Host "API Server: http://localhost:8081"
Write-Host ""
Write-Host "🧪 NEXT STEPS:" -ForegroundColor Cyan
Write-Host "1. Import Postman collection: Finora_API.postman_collection.json"
Write-Host "2. Test endpoints starting with: GET /health"
Write-Host "3. Try authentication: POST /api/auth/send-otp"
Write-Host ""
Write-Host "Happy coding! 🚀" -ForegroundColor Green
