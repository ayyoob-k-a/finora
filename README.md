Finora 🔐

Modern Authentication Service built with Go

Finora is a secure and scalable authentication backend designed for mobile and web apps.
It provides OTP-based email verification & login, built with a clean Go architecture for easy integration and extension.

✨ Features

🔑 OTP-based authentication (email login)

📧 Email delivery via SMTP (Gmail supported)

🗄️ PostgreSQL for persistent storage

⚡ Redis for caching & OTP expiry handling

🔒 Password-less, secure login flow

🧩 Modular architecture using repository pattern

🚀 Production-ready setup with environment configs

🛠️ Tech Stack

Language: Go (Golang)

Database: PostgreSQL

Cache: Redis

Email Service: SMTP (Gmail, extensible)

Architecture: Clean, modular, repository-based

🚀 Getting Started
1️⃣ Clone the Repository
git clone https://github.com/your-username/finora.git
cd finora

2️⃣ Setup Environment Variables

Create a .env file in the root:

PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/finora
REDIS_URL=localhost:6379
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

3️⃣ Run the Server
go run cmd/main.go

📂 Project Structure
finora/
│── cmd/             # Entry point (main.go)
│── internal/
│   ├── auth/        # OTP & authentication logic
│   ├── db/          # Database layer
│   ├── cache/       # Redis integration
│   └── email/       # SMTP email service
│── pkg/             # Utility packages
│── .env.example     # Example environment variables
│── go.mod
│── go.sum

📡 API Endpoints
Method	Endpoint	Description
POST	/auth/send-otp	Send OTP to user’s email
POST	/auth/verify	Verify OTP & login user
🧪 Testing
go test ./...

📜 License

MIT License © 2025 [AYYOOB]

👉 I can also add badges (Go version, CI build status, license, PostgreSQL, Redis) at the top for a more professional open-source look.

Do you want me to include those badges in this README?
