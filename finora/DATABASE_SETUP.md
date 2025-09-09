# 🗄️ Finora Database Setup Guide

This guide will help you set up PostgreSQL database for the Finora API from scratch.

## 🚀 Quick Setup Options

### Option A: PostgreSQL Installation (Windows)
### Option B: Using Docker (if available)

---

## 📦 Option A: Manual PostgreSQL Setup (Windows)

### Step 1: Install PostgreSQL

1. **Download PostgreSQL**:
   - Go to: https://www.postgresql.org/download/windows/
   - Download PostgreSQL 15 or 16 (latest stable version)
   - Choose the Windows x86-64 installer

2. **Run the installer**:
   - Install with default settings
   - **IMPORTANT**: Remember the password you set for the `postgres` user
   - Default port: `5432` (keep this)
   - Install pgAdmin 4 (recommended for database management)

3. **Add to PATH** (if not done automatically):
   - Add `C:\Program Files\PostgreSQL\15\bin` to your Windows PATH

### Step 2: Verify Installation

Open a new PowerShell window and run:
```powershell
psql --version
```

You should see something like: `psql (PostgreSQL) 15.x`

### Step 3: Create Database and User

1. **Connect as superuser**:
```powershell
psql -U postgres -h localhost
```
*Enter the password you set during installation*

2. **Create database and user**:
```sql
-- Create database
CREATE DATABASE finora_db;

-- Create user with password
CREATE USER finora_user WITH ENCRYPTED PASSWORD 'finora_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE finora_db TO finora_user;
GRANT USAGE ON SCHEMA public TO finora_user;
GRANT CREATE ON SCHEMA public TO finora_user;

-- Connect to finora_db and grant table privileges
\c finora_db
GRANT ALL ON ALL TABLES IN SCHEMA public TO finora_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO finora_user;
GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO finora_user;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO finora_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO finora_user;

-- Exit
\q
```

### Step 4: Test Connection

```powershell
psql -U finora_user -d finora_db -h localhost
```
*Password: `finora_password`*

If successful, you should see the database prompt: `finora_db=>`

---

## 🐳 Option B: Docker Setup (Alternative)

If you prefer Docker (and have it installed):

### Step 1: Install Docker Desktop
- Download from: https://www.docker.com/products/docker-desktop/
- Install and restart your computer
- Start Docker Desktop

### Step 2: Start Database
```powershell
# Navigate to project directory
cd C:\Users\User\Documents\ayyoob\finora\finora

# Start only PostgreSQL from docker-compose
docker-compose up postgres -d

# Check if running
docker-compose ps
```

### Step 3: Access Database
```powershell
# Connect to database
docker-compose exec postgres psql -U finora_user -d finora_db
```

---

## ⚙️ Environment Variables

After setting up the database, update your environment variables:

### For Manual PostgreSQL:
```powershell
$env:DB_HOST='localhost'
$env:DB_PORT='5432'
$env:DB_USER='finora_user'
$env:DB_PASSWORD='finora_password'
$env:DB_NAME='finora_db'
$env:DB_SSLMODE='disable'
```

### For Docker:
```powershell
$env:DB_HOST='localhost'
$env:DB_PORT='5432'
$env:DB_USER='finora_user'
$env:DB_PASSWORD='finora_password'
$env:DB_NAME='finora_db'
$env:DB_SSLMODE='disable'
```

---

## 🧪 Test Database Connection

After setup, test the connection with our Go application:

```powershell
# Set all environment variables
$env:JWT_SECRET='development-secret-key-12345678'
$env:PORT='8081'
$env:DB_HOST='localhost'
$env:DB_PORT='5432'
$env:DB_USER='finora_user'
$env:DB_PASSWORD='finora_password'
$env:DB_NAME='finora_db'
$env:DB_SSLMODE='disable'

# Start the API server
go run main.go
```

You should see:
- ✅ `Database connection successful!`
- ✅ `Running database migrations...`
- ✅ `Seeding default categories...`
- ✅ `[GIN-debug] Listening and serving HTTP on :8081`

---

## 🔧 Database Management Tools

### pgAdmin 4 (Recommended)
- **URL**: http://localhost:5050 (if using Docker)
- **Manual**: Start pgAdmin from Windows Start Menu
- **Connection Details**:
  - Host: `localhost`
  - Port: `5432`
  - Database: `finora_db`
  - Username: `finora_user`
  - Password: `finora_password`

### Command Line
```powershell
# Connect to database
psql -U finora_user -d finora_db -h localhost

# Common commands
\dt          -- List tables
\d users     -- Describe users table
\q           -- Quit
```

---

## 🚨 Troubleshooting

### "Connection refused" Error
- ✅ PostgreSQL service is running: `services.msc` → Find "postgresql" service
- ✅ Correct port (5432) is open
- ✅ Firewall allows PostgreSQL

### "Authentication failed" Error  
- ✅ Username/password are correct
- ✅ User has proper permissions on database
- ✅ Try connecting with `postgres` superuser first

### "Database does not exist" Error
- ✅ Database `finora_db` was created
- ✅ Connect as `postgres` user and recreate database

### Permission Errors
- ✅ Run the GRANT commands from Step 3 again
- ✅ Make sure you're connected to the right database when granting permissions

---

## 🎯 What's Next?

After successful database setup:

1. **✅ Database Connected** - PostgreSQL running
2. **✅ Tables Created** - Auto-migration will create all tables
3. **✅ Default Data** - Categories and sample data seeded
4. **✅ API Ready** - All endpoints work with real database
5. **✅ Test with Postman** - Full functionality testing

---

## 📋 Quick Reference

### Database Details
- **Host**: `localhost`
- **Port**: `5432`
- **Database**: `finora_db`
- **User**: `finora_user`
- **Password**: `finora_password`

### Environment Variables
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=finora_user
DB_PASSWORD=finora_password
DB_NAME=finora_db
DB_SSLMODE=disable
JWT_SECRET=development-secret-key-12345678
PORT=8081
```

Ready to proceed with database setup! Choose Option A (Manual) or Option B (Docker) and follow the steps.
