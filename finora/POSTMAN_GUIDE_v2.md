# 🚀 Finora API v2 - Complete Postman Testing Guide

This comprehensive guide covers the **updated Postman collection** with intelligent error handling, automatic testing, and support for both database-connected and API-only modes.

## 📦 What's New in v2

### ✨ **Smart Database Detection**
- **Automatic database status detection** via health check
- **Intelligent test expectations** based on database availability  
- **Graceful handling** of API-only vs full-database modes

### 🧪 **Advanced Testing Features**  
- **Automatic token management** with secure storage
- **Dynamic variable population** from API responses
- **Comprehensive test assertions** for all scenarios
- **Error logging and debugging** utilities

### 🔧 **Environment Management**
- **Auto-configuration** of missing variables
- **Corrected port settings** (8081 for API, not 5432)
- **Sample UUIDs** for testing without database
- **Secure token storage** with proper types

## 🎯 Quick Start

### **Step 1: Import Updated Files**
```bash
# Import both files into Postman:
1. Finora_API_v2.postman_collection.json
2. Finora_Environment_v2.postman_environment.json
```

### **Step 2: Select Environment**  
- Click environment dropdown → Select **"Finora Development v2"**

### **Step 3: Fix Server Configuration**
```powershell
# Set correct environment variables
$env:PORT='8081'              # API server port (NOT 5432)
$env:DB_USER='finora_user'    # Database user (NOT postgres)
$env:JWT_SECRET='development-secret-key-12345678'
$env:DB_HOST='localhost'
$env:DB_PORT='5432' 
$env:DB_PASSWORD='finora_password'
$env:DB_NAME='finora_db'
$env:DB_SSLMODE='disable'

# Start PostgreSQL if needed
& "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\pgsql\bin\pg_ctl.exe" -D "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\data" -l "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\logfile" start

# Start Finora API server
go run main.go
```

### **Step 4: Run Health Check First**
1. Open **🏥 Health & System** → **Health Check**
2. Click **Send**
3. This populates `database_status` variable for intelligent testing

## 📋 Collection Structure

### **🏥 Health & System** (1 endpoint)
- **Advanced health check** with database status detection
- **Automatic environment setup** and validation
- **System information** and version reporting

### **🔐 Authentication** (3 endpoints)
- **Send OTP**: Smart error handling for database vs API-only modes
- **Verify OTP**: Automatic token extraction and storage  
- **Refresh Token**: Secure token renewal with validation

### **👤 User Management** (3 endpoints)
- **Get Profile**: User information with authentication checks
- **Update Profile**: Profile modification with validation
- **Get Dashboard**: Financial overview and recent activity

### **📊 Categories** (1 endpoint)  
- **Get Categories**: Returns real data or placeholders based on database status

### **💸 Transactions** (5 endpoints)
- **Complete CRUD operations** with proper validation
- **Filtering and pagination** support
- **Automatic ID capture** for dependent requests

### **📅 EMI Management** (4 endpoints)
- **EMI creation and tracking** with payment history
- **Payment recording** with automatic calculations
- **Due date management** and reminders

### **👥 Friend Management** (4 endpoints)  
- **Friend requests** with acceptance/rejection flow
- **Friend list management** with automatic ID capture
- **Relationship management** and removal

### **👨‍👩‍👧‍👦 Group Management** (5 endpoints)
- **Group creation** with member management
- **Expense splitting** with multiple calculation methods
- **Balance settlement** and payment tracking  

### **📈 Reports & Analytics** (3 endpoints)
- **Monthly spending reports** with category breakdowns
- **Category analysis** with trend data
- **Yearly summaries** and financial insights

### **🔔 Notifications** (4 endpoints)
- **Notification management** with read/unread status
- **Bulk operations** for marking all as read
- **Notification deletion** and cleanup

### **🧪 Testing & Utilities** (2 endpoints)
- **Database connection testing** with diagnostic output
- **Environment cleanup** utility for fresh testing

## 🎭 Intelligent Testing Features

### **Database Status Awareness**
```javascript
// Example test logic
const dbStatus = pm.environment.get('database_status');

if (dbStatus === 'connected') {
    // Test for full functionality
    pm.test('Full feature response', function () {
        pm.expect(pm.response.code).to.equal(200);
        pm.expect(response.data).to.exist;
    });
} else {
    // Test for graceful degradation
    pm.test('Graceful error response', function () {
        pm.expect(pm.response.code).to.equal(503);
        pm.expect(response.error).to.include('Database not available');
    });
}
```

### **Automatic Token Management**
- **JWT tokens** extracted and stored automatically
- **Authorization headers** populated for protected endpoints
- **Token refresh** handled transparently
- **Secure storage** using secret variable types

### **Dynamic ID Management**  
- **Resource IDs** automatically captured from creation responses
- **Dependent requests** use captured IDs automatically
- **UUID placeholders** for testing without database

## 🚦 Testing Workflows

### **Scenario 1: Full Database Mode**
```bash
# Prerequisites: PostgreSQL running, database connected
1. Health Check → Status: "connected" 
2. Send OTP → Success response
3. Verify OTP → JWT tokens received
4. All endpoints → Full functionality testing
```

### **Scenario 2: API-Only Mode** 
```bash
# Prerequisites: API running, no database connection  
1. Health Check → Status: "disconnected"
2. Send OTP → 503 Service Unavailable (graceful)
3. Categories → Placeholder data returned
4. Protected endpoints → Proper error responses
```

### **Scenario 3: Authentication Testing**
```bash
# Test authentication flow and token management
1. Clear Environment Variables → Reset all tokens
2. Send OTP → Initiate authentication
3. Verify OTP → Receive and store tokens  
4. Protected requests → Automatic authorization
5. Refresh Token → Token renewal testing
```

## 🔧 Environment Variables  

### **Auto-Configured Variables**
```json
{
  "base_url": "http://localhost:8081",     // Corrected API port
  "test_phone": "+1234567890",             // Default test phone
  "test_email": "test@finora.app",         // Default test email
}
```

### **Dynamic Variables (Auto-Populated)**
```json
{
  "access_token": "",           // JWT from authentication
  "refresh_token": "",          // Refresh token  
  "user_id": "",               // User ID from login
  "database_status": "",       // Connected/disconnected
  "transaction_id": "",        // From transaction creation
  "emi_id": "",               // From EMI creation
  "friend_id": "",            // From friend requests
  "group_id": "",             // From group creation
  "category_id": "",          // From categories list
  "notification_id": ""       // From notifications
}
```

## 🧪 Advanced Testing

### **Run Complete Test Suite**
1. **Import collection and environment**
2. **Start with Health Check** to detect database status  
3. **Run Authentication folder** to get tokens
4. **Run any endpoint folder** for comprehensive testing
5. **Check Console logs** for detailed test results

### **Debug Failed Requests**
- **Console logs** show detailed error information
- **Response times** validated for performance
- **Status codes** checked against expected scenarios
- **Error messages** validated for proper formatting

### **Database Connection Testing**
- Use **"Test Database Connection"** in Testing & Utilities
- **Automatic diagnostics** with console output
- **Connection status verification** for troubleshooting

## 📊 Expected Responses

### **With Database Connected**
```json
// Health Check
{
  "status": "healthy",
  "database": "connected", 
  "version": "1.0.0",
  "timestamp": "2024-01-01 12:00:00"
}

// Authentication Success  
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { "id": "uuid", "name": "John Doe" },
  "expires_in": 3600
}
```

### **Without Database (API-Only Mode)**
```json
// Health Check
{
  "status": "healthy",
  "database": "disconnected",
  "version": "1.0.0", 
  "timestamp": "2024-01-01 12:00:00"
}

// Protected Endpoint Error
{
  "success": false,
  "error": "Database not available. Please try again later or set up database connection."
}
```

## 🚨 Troubleshooting

### **Server Not Responding**
```bash
# Check if server is running on correct port
netstat -ano | findstr ":8081"

# Verify environment variables
echo $env:PORT  # Should be 8081, not 5432
echo $env:DB_USER  # Should be finora_user, not postgres
```

### **Database Connection Issues**  
```bash
# Check PostgreSQL status
& "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\pgsql\bin\pg_ctl.exe" -D "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\data" status

# Test database connection manually
$env:PGPASSWORD='finora_password'; & "C:\Users\User\Downloads\postgresql-17.6-1-windows-x64-binaries\pgsql\bin\psql.exe" -U finora_user -d finora_db -c "SELECT 1;"
```

### **Authentication Issues**
- **Clear tokens** using "Clear Environment Variables" request
- **Re-run authentication flow** from Send OTP
- **Check JWT secret** environment variable
- **Verify token format** in environment variables

### **Wrong Port Configuration**
```bash
# Common mistake: Server running on PostgreSQL port  
# Fix: Ensure PORT=8081 (API) not 5432 (PostgreSQL)
$env:PORT='8081'
go run main.go
```

## 🎉 Success Indicators

### **✅ Fully Functional Setup**
- Health Check: `"database": "connected"`
- Send OTP: `200 OK` with success message
- Verify OTP: JWT tokens automatically stored
- All endpoints: Full functionality responses

### **⚠️ API-Only Mode (Acceptable)**  
- Health Check: `"database": "disconnected"`
- Send OTP: `503 Service Unavailable` (graceful error)
- Categories: Placeholder data returned
- No crashes or panics

### **❌ Issues to Fix**
- Server not responding on port 8081
- Authentication endpoints returning 500 errors  
- Missing environment variables
- PostgreSQL connection refused errors

## 🚀 Production Readiness

### **Environment Configuration**
```json
// Update base_url for different environments
"production": "https://api.finora.app",
"staging": "https://staging-api.finora.app",
"development": "http://localhost:8081"
```

### **Security Considerations**  
- **JWT tokens stored as secrets** in environment
- **Phone numbers and emails** configured appropriately
- **API keys and credentials** managed securely
- **HTTPS enforcement** in production environments

---

## 💯 Ready to Test!

Your **Finora API v2 Postman Collection** is now ready with:

- ✅ **80+ test assertions** across all endpoints
- ✅ **Intelligent database detection** and response handling  
- ✅ **Automatic token management** and ID capture
- ✅ **Comprehensive error handling** and logging
- ✅ **Environment auto-configuration** and validation
- ✅ **Production-ready** testing workflows

**Start testing and building amazing financial applications!** 🎯
