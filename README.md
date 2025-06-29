# Go Microservices Boilerplate

A production-ready microservices architecture built with Go, featuring a React + TypeScript frontend and multiple backend services orchestrated with Docker. RabbitMQ handles asynchronous communication with the broker service as the single point of entry.

---

## 🏗️ Architecture Overview

```
┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐
│   React Client  │      │  Broker Service  │      │   Auth Service  │
│  (Frontend UI)  │────▶│  (API Gateway)   │────▶ │   PostgreSQL    │
└─────────────────┘      └──────────────────┘      └─────────────────┘
                                 │                           │
                                 │                           │ (publishes events)
                                 ▼                           ▼
                        ┌─────────────────────────────────────────┐
                        │            RabbitMQ                     │
                        │      (Message Queue/Event Bus)          │
                        └─────────────────────────────────────────┘
                                 │                   │
                                 │ (consumes)        │ (consumes)
                                 ▼                   ▼
                        ┌─────────────────┐    ┌─────────────────┐
                        │  Logger Service │    │  Mailer Service │
                        │    MongoDB      │    │   (SMTP Send)   │
                        └─────────────────┘    └─────────────────┘
```

**Broker Service**: Single point of entry routing all external requests to appropriate microservices  
**RabbitMQ**: Handles asynchronous communication between services (listener functionality built into broker)

Event Flow:
1. Auth Service publishes events → RabbitMQ (login attempts, registrations, etc.)
2. Broker Service publishes events → RabbitMQ (API calls, errors, etc.)
3. Logger Service consumes events ← RabbitMQ (stores logs in MongoDB)
4. Mailer Service consumes events ← RabbitMQ (sends emails via SMTP)

Example Events:
- user.registered → triggers welcome email + logging
- user.login.failed → triggers security email + error logging  
- api.error → triggers admin notification + error logging

---

## 📁 Project Structure

```
go-microservices/
├── client/                     # React + TypeScript Frontend
├── server/
│   ├── broker/                 # API Gateway + Message Listener
│   ├── auth/                   # Authentication (PostgreSQL)
│   ├── logger/                 # Logging Service (MongoDB)
│   ├── mailer/                 # Email Service (SMTP)
│   ├── docker-compose.yml
│   └── .env
```

---

## 🛠️ Tech Stack

**Backend**: Go, Chi/Gorilla Mux, PostgreSQL, MongoDB, RabbitMQ, Docker  
**Frontend**: React, TypeScript, Tailwind CSS, Vite

---

## 🔧 Services

### 🚪 **Broker Service** (Port: 8084)
Single point of entry for all external requests with built-in RabbitMQ listener

**Endpoints:**
- `POST /handle` - Routes actions: `auth`, `log`, `logdirect`, `mail`

### 🔐 **Auth Service** (Port: 8081) 
User authentication with PostgreSQL storage
- `POST /register` - User registration
- `POST /login` - User authentication

### 📝 **Logger Service** (Port: 8082)
Centralized logging with MongoDB storage
- `GET /api/v1/logs` - Get all logs
- `POST /api/v1/logs` - Create log entry
- `GET /api/v1/logs/stats` - Log statistics

### 📧 **Mailer Service** (Port: 8083)
Email service with SMTP support
- `POST /api/v1/send` - Send single email
- `POST /api/v1/send/batch` - Send batch emails

---

## 🚀 Quick Start

### 1. Environment Setup
Create `.env` in `server/` directory:

```env
# Databases
POSTGRES_URL=postgres://postgres:password@postgres:5432/users?sslmode=disable
MONGO_URL=mongodb://mongo:27017/logs

# RabbitMQ
RABBITMQ_HOST=rabbitmq
RABBITMQ_USER=admin
RABBITMQ_PASS=password

# Service Ports
BROKER_PORT=8084
AUTH_PORT=8081
LOGGER_PORT=8082
MAILER_PORT=8083

# Email Configuration
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password
FROM_ADDRESS=noreply@yourapp.com
```

### 2. Start Services
```bash
cd server
docker-compose build
docker-compose up -d
```

### 3. Start Frontend
```bash
cd client
npm install && npm run dev
```

---

## 📡 API Usage

All requests go through the broker service at `http://localhost:8084/handle`

### Authentication
```json
{
  "action": "auth",
  "auth": {
    "email": "user@example.com",
    "password": "password123"
  }
}
```

### Direct Logging
```json
{
  "action": "logdirect",
  "log": {
    "name": "user-action",
    "data": "User logged in"
  }
}
```

### Async Logging (via RabbitMQ)
```json
{
  "action": "log",
  "log": {
    "name": "async-event",
    "data": "Processed asynchronously"
  }
}
```

### Send Email
```json
{
  "action": "mail",
  "mail": {
    "to": "user@example.com",
    "subject": "Welcome!",
    "message": "Welcome to our platform!"
  }
}
```

---

## 🔍 Service URLs

- **Frontend**: http://localhost:3000
- **Broker (API Gateway)**: http://localhost:8084
- **RabbitMQ Management**: http://localhost:15672 (admin/password)

---

## 🚀 Production Features

- **Docker containerization** for easy deployment
- **Environment-based configuration** 
- **Health checks** and monitoring
- **CORS** and security middleware
- **Rate limiting** on email service
- **Async processing** with RabbitMQ
- **Independent service scaling**

---

## 📜 License

MIT License - Ready for production use with enterprise-grade patterns.