
# Full Stack Microservices Boilerplate

This is a full-stack microservices project setup with a React + TypeScript frontend and a Golang backend using Docker and PostgreSQL. It is production-ready and easy to scale with reusable architecture.

---

## 📁 Folder Structure

- /client → React frontend
- /server → Go backend microservices (auth, broker, etc.)
- /server/docker-compose.yml → Orchestrates services with Postgres

---

## 📦 Tech Stack

### Backend (server/)
- Go (Golang)
- PostgreSQL
- Docker + Docker Compose
- Chi - Router

### Frontend (client/)
- React
- TypeScript
- Tailwind CSS

---

## 🐳 Docker Compose

```bash
# Build all services
docker compose build

# Start all services
docker compose up

# Stop all services
docker compose down
```

---

## 🧠 Microservice Architecture

- **Auth Service:** Handles user registration, login, JWT generation
- **Broker Service:** Routes external/internal service communication
- [More microservices like Mail, Notification, etc. to be added]

---

## 🌍 Environment Variables

There is only a single `.env` file in server directory.

```bash

MONGO_URL=mongodb://mongo:27017
RABBITMQ_URL=amqp://rabbitmq:5672
POSTGRES_URL=postgres://postgres:<password>@localhost:5432/db

# Specific ports for each service
AUTH_PORT=8081
MAIL_PORT=8082
LISTENER_PORT=8083
BROKER_PORT=8084

```

---

## 🧪 Scripts

In `client/`:

```bash
npm install
npm run dev
```

In `server/`:

```bash
docker compose build
```

```bash
docker compose up
```

---

## 🛠️ Tips

- Use  `go-migrate` for handling migrations
- Use separate Dockerfiles for each microservice for better scaling
- Modularize Go code with `internal/` packages

---

## 📜 License

MIT# go-microservices
