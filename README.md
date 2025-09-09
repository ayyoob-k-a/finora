# Finora API v3 - Postman Collection

[![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)](https://www.postman.com/)
[![API](https://img.shields.io/badge/API-v3-blue?style=for-the-badge)](https://github.com/your-username/finora-api-v3)
[![Status](https://img.shields.io/badge/Status-Active-green?style=for-the-badge)](https://github.com/your-username/finora-api-v3)

A comprehensive Postman collection for testing the Finora API v3 - a personal finance management system with transaction tracking, EMI management, social features, and reporting capabilities.

## 🚀 Features

- **Health Check**: Monitor API status and database connectivity
- **Authentication**: OTP-based login system with JWT tokens
- **Transaction Management**: Create, retrieve, and manage financial transactions
- **EMI Tracking**: Loan management with payment scheduling
- **Social Features**: Friend requests and expense groups
- **Reporting**: Monthly spending reports and analytics
- **Notifications**: Real-time notification system
- **Categories**: Expense categorization for better tracking

## 📋 Prerequisites

- [Postman Desktop App](https://www.postman.com/downloads/) or Postman Web
- Finora API server running (default: `http://localhost:8081`)
- Basic understanding of REST APIs and authentication

## 🔧 Installation & Setup

### 1. Import the Collection

**Option A: Direct Import**
1. Download the `finora-api-v3-fixed.postman_collection.json` file
2. Open Postman
3. Click **Import** → **Upload Files**
4. Select the downloaded JSON file

**Option B: Import from URL**
1. Copy the raw GitHub URL of the collection file
2. In Postman, click **Import** → **Link**
3. Paste the URL and import

### 2. Environment Setup

The collection includes auto-configuration, but you can customize these variables:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `base_url` | `http://localhost:8081` | API server URL |
| `test_phone` | `+1234567890` | Test phone number |
| `test_email` | `test@finora.app` | Test email address |

**To customize:**
1. Create a new Environment in Postman
2. Add the variables above with your preferred values
3. Select the environment before running requests

## 🏃‍♂️ Quick Start

### Step 1: Health Check
Run the **🏥 Health Check** request to verify:
- API server is running
- Database connection status
- System version information

### Step 2: Authentication Flow
1. **🔐 Send OTP**: Request OTP to your phone/email
2. **🔐 Verify OTP**: Verify with OTP `123456` (test environment)

> **Note**: The collection automatically saves authentication tokens for subsequent requests.

### Step 3: Basic Operations
1. **📊 Get Categories**: Fetch available expense categories
2. **💸 Create Transaction**: Add your first transaction
3. **💸 Get Transactions**: Retrieve transaction history

## 📚 API Endpoints Overview

### Authentication
- `POST /api/auth/send-otp` - Send OTP for login
- `POST /api/auth/verify-otp` - Verify OTP and get JWT tokens

### Transactions
- `POST /api/transactions` - Create new transaction
- `GET /api/transactions` - Get transactions with pagination and filters

### EMI Management  
- `POST /api/emis` - Create EMI record
- `POST /api/emis/{id}/payment` - Record EMI payment

### Social Features
- `POST /api/friends/request` - Send friend request
- `POST /api/groups` - Create expense group

### Reports & Analytics
- `GET /api/reports/monthly` - Monthly spending reports
- `GET /api/notifications` - User notifications
- `GET /api/categories` - Expense categories

## 🧪 Testing Features

### Automated Tests
Each request includes comprehensive test scripts that verify:
- Response status codes
- Required response fields
- Data structure validation
- Authentication token management

### Database-Aware Testing
The collection intelligently handles both:
- **Full Database Mode**: Complete functionality testing
- **API-Only Mode**: Graceful degradation when database is unavailable

### Environment Variables
Tests automatically manage:
- Authentication tokens (`access_token`, `refresh_token`)
- Entity IDs (`user_id`, `transaction_id`, `emi_id`, etc.)
- Test data persistence across requests

## 🔍 Request Examples

### Create Transaction
```json
{
  "type": "expense",
  "amount": 25.50,
  "category_id": "cat_123",
  "description": "Coffee and pastry",
  "transaction_date": "2024-01-15T10:30:00Z",
  "is_recurring": false
}
```

### Create EMI
```json
{
  "title": "Car Loan",
  "total_amount": 25000.00,
  "monthly_amount": 450.00,
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2029-01-01T00:00:00Z",
  "due_date": 5,
  "description": "Monthly car loan payment"
}
```

## 🐛 Troubleshooting

### Common Issues

**Authentication Errors (401)**
- Ensure you've completed the OTP flow
- Check if tokens are saved in environment variables
- Verify the `Authorization` header format: `Bearer {token}`

**Database Connection Issues (503)**
- The API runs in API-only mode without database
- Some features may return placeholder data
- Check server logs for database connectivity

**Request Timeouts**
- Default timeout is set to 5 seconds
- Increase timeout in Postman settings if needed
- Check server performance and network connectivity

### Debug Mode
Enable Postman Console (View → Show Postman Console) to see:
- Detailed request/response logs
- Test execution results
- Environment variable changes

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-endpoint`
3. Update the collection with new requests
4. Add comprehensive tests for new endpoints
5. Update this README with new features
6. Submit a pull request

### Adding New Endpoints
When adding new requests:
1. Follow the naming convention: `🔥 Endpoint Name`
2. Add proper descriptions and examples
3. Include test scripts for validation
4. Update environment variables if needed

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Links

- [Finora API Documentation](https://docs.finora.app)
- [Postman Documentation](https://learning.postman.com/)
- [JWT.io](https://jwt.io/) - JWT token debugging

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/finora-api-v3/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/finora-api-v3/discussions)
- **Email**: finoraisme@gmail.com

---

**Made with ❤️ for the Finora community**

> This collection is actively maintained and updated with the latest API changes. Star ⭐ this repository to stay updated!
