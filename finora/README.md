# Finora - Personal Expense Management API

Finora is a comprehensive Go backend API for personal expense management. It provides features for tracking expenses, managing EMIs, group expense splitting, friend management, and detailed financial reporting.

## 🚀 Features Implemented

### ✅ Core Authentication System
- **OTP-based Authentication**: Email and SMS OTP verification using Twilio
- **JWT Token Management**: Secure authentication with refresh token support
- **User Registration & Profile Management**: Complete user lifecycle management

### ✅ Database & Models
- **PostgreSQL Integration**: Full database setup with GORM
- **UUID Primary Keys**: All models use UUIDs for better security and scalability
- **Auto-Migration**: Automatic database schema creation and updates
- **Comprehensive Models**: User, Transaction, EMI, Friend, Group, Notification, Category, etc.

### ✅ Security & Middleware
- **JWT Authentication Middleware**: Secure API endpoint protection
- **Rate Limiting**: Built-in rate limiting for OTP and general API requests
- **CORS Support**: Configurable cross-origin resource sharing
- **Input Validation**: Comprehensive request validation and sanitization
- **Security Headers**: XSS protection, content type validation, etc.

### ✅ API Endpoints (Implemented)
- **Authentication**: `/api/auth/send-otp`, `/api/auth/verify-otp`, `/api/auth/refresh`
- **User Management**: `/api/user/profile`, `/api/user/dashboard`
- **Health Check**: `/health` endpoint for monitoring

### 📋 API Endpoints (Placeholders Ready)
The following endpoints have routing and placeholder implementations ready:
- **Categories**: CRUD operations for expense categories
- **Transactions**: Complete transaction management
- **EMIs**: EMI tracking and payment recording
- **Friends**: Friend request and management system
- **Groups**: Group expense management with split calculations
- **Reports**: Monthly and category-wise reporting
- **Notifications**: Push notification system

## 🛠️ Technology Stack

- **Language**: Go 1.23.2
- **Framework**: Gin Gonic
- **Database**: PostgreSQL with GORM
- **Cache**: Redis (ready for implementation)
- **Authentication**: JWT with custom claims
- **Email**: SMTP with HTML templates
- **SMS**: Twilio integration
- **Containerization**: Docker & Docker Compose

## 📦 Dependencies

```go
require (
    github.com/gin-gonic/gin v1.10.1
    github.com/golang-jwt/jwt/v5 v5.2.3
    github.com/google/uuid v1.6.0
    gorm.io/gorm v1.25.12
    gorm.io/driver/postgres v1.6.0
    github.com/redis/go-redis/v9 v9.3.1
    github.com/twilio/twilio-go v1.22.3
    github.com/go-playground/validator/v10 v10.27.0
    github.com/joho/godotenv v1.5.1
    github.com/gin-contrib/cors v1.7.2
    golang.org/x/crypto v0.40.0
    github.com/robfig/cron/v3 v3.0.1
    gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
    github.com/stretchr/testify v1.10.0
    go.uber.org/zap v1.27.0
)
```

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone and navigate to the project**:
   ```bash
   cd finora
   ```

2. **Create environment file**:
   ```bash
   # Create .env file with your configurations
   cp .env.example .env
   ```

3. **Configure your environment variables**:
   ```env
   # Database (automatically configured in docker-compose)
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=finora_user
   DB_PASSWORD=finora_password
   DB_NAME=finora_db
   
   # JWT Secret (IMPORTANT: Change this in production!)
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   
   # SMTP Configuration (for email OTP)
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-app-password
   
   # Twilio Configuration (for SMS OTP)
   TWILIO_ACCOUNT_SID=your-twilio-account-sid
   TWILIO_AUTH_TOKEN=your-twilio-auth-token
   TWILIO_PHONE_NUMBER=+1234567890
   ```

4. **Start the services**:
   ```bash
   docker-compose up -d
   ```

5. **Access the services**:
   - **API**: http://localhost:8080
   - **PostgreSQL**: localhost:5432
   - **Redis**: localhost:6379
   - **pgAdmin** (development): http://localhost:5050

### Manual Setup

1. **Install PostgreSQL and Redis**

2. **Install Go dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up environment variables** (create `.env` file)

4. **Run the application**:
   ```bash
   go run main.go
   ```

## 📚 API Documentation

### Authentication Endpoints

#### Send OTP
```http
POST /api/auth/send-otp
Content-Type: application/json

{
  "phone": "+1234567890",
  "email": "user@example.com"
}
```

#### Verify OTP
```http
POST /api/auth/verify-otp
Content-Type: application/json

{
  "phone": "+1234567890",
  "otp": "123456"
}
```

#### Refresh Token
```http
POST /api/auth/refresh
Content-Type: application/json
Authorization: Bearer <refresh_token>

{
  "refresh_token": "<refresh_token>"
}
```

### User Management Endpoints

#### Get Profile
```http
GET /api/user/profile
Authorization: Bearer <jwt_token>
```

#### Update Profile
```http
PUT /api/user/profile
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "John Doe Updated",
  "photo_url": "https://example.com/photo.jpg",
  "default_currency": "USD",
  "monthly_income": 5000.00
}
```

#### Get Dashboard
```http
GET /api/user/dashboard
Authorization: Bearer <jwt_token>
```

## 🏗️ Project Structure

```
finora/
├── configs/              # Configuration management
├── db/                   # Database connection and migrations
├── domain/               # Database models (entities)
├── handler/              # HTTP handlers (controllers)
├── middleware/           # Custom middleware
├── model/
│   └── dto/             # Data Transfer Objects
├── service/             # Business logic layer
├── utils/               # Utility functions (auth, validation, etc.)
├── docker-compose.yml   # Docker services configuration
├── Dockerfile          # Application containerization
├── go.mod              # Go module dependencies
├── main.go            # Application entry point
└── README.md          # This file
```

## 🔧 Configuration

The application supports configuration via environment variables:

### Database Configuration
- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSLMODE`: SSL mode (default: disable)

### JWT Configuration
- `JWT_SECRET`: Secret key for JWT signing (REQUIRED)
- `JWT_EXPIRY`: Token expiration time (default: 24h)
- `REFRESH_TOKEN_EXPIRY`: Refresh token expiration (default: 168h)

### Server Configuration
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release, default: debug)

### Email Configuration (Optional)
- `SMTP_HOST`: SMTP server host
- `SMTP_PORT`: SMTP server port
- `SMTP_USERNAME`: SMTP username
- `SMTP_PASSWORD`: SMTP password/app password

### SMS Configuration (Optional)
- `TWILIO_ACCOUNT_SID`: Twilio Account SID
- `TWILIO_AUTH_TOKEN`: Twilio Auth Token
- `TWILIO_PHONE_NUMBER`: Twilio phone number

## 🔐 Security Features

- **JWT Authentication**: Stateless authentication with configurable expiration
- **Rate Limiting**: Prevents abuse with customizable limits
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: GORM provides built-in protection
- **CORS Configuration**: Configurable cross-origin policies
- **Security Headers**: XSS protection, content type validation
- **Password Security**: Bcrypt hashing for sensitive data

## 📈 What's Next?

The foundation is solid! Here's what needs to be implemented to complete the full API:

### High Priority
1. **Transaction Management System** - CRUD operations, filtering, pagination
2. **EMI Management System** - Payment tracking, reminders
3. **Category Management** - CRUD for expense categories

### Medium Priority
4. **Friend Management System** - Friend requests, acceptance, removal
5. **Group Expense Management** - Split calculations, settlements
6. **Notification System** - Push notifications, EMI reminders

### Low Priority
7. **Reporting System** - Monthly reports, category breakdowns
8. **Background Jobs** - Recurring transactions, automated reminders
9. **Redis Caching** - Performance optimization
10. **Advanced Security** - Enhanced rate limiting, audit logs

## 🐛 Testing

Health check endpoint is available:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the project
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📞 Support

For support, email support@finora.com or create an issue in the repository.

---

**Built with ❤️ using Go and modern best practices**
