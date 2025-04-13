
# Full Stack Microservices Boilerplate

This is a full-stack microservices project setup with a React + TypeScript frontend and a Golang backend using Docker and PostgreSQL. It is production-ready and easy to scale with reusable architecture.

---

## 📁 Folder Structure

- /client → React frontend
- /server → Go backend microservices (auth, broker, etc.)
- /server/docker-compose.yml → Orchestrates services with Postgres

---

## 📦 Tech Stack

### Frontend (client/)
- React
- TypeScript
- Redux Toolkit
- Vite
- Tailwind CSS (optional)

### Backend (server/)
- Go (Golang)
- PostgreSQL
- Docker + Docker Compose
- Gorilla Mux
- sqlc / GORM (if used for DB layer)

---

## 🐳 Docker Compose

```bash
# Start all services
docker-compose up --build

# Stop all services
docker-compose down
```

---

## 🧠 Microservice Architecture

- **Auth Service:** Handles user registration, login, JWT generation
- **Broker Service:** Routes external/internal service communication
- [Add more microservices like Mail, Notification, etc.]

---

## 🌍 Environment Variables

Each microservice should have its own `.env` file. Common variables:
- `PORT`
- `DATABASE_URL`
- `JWT_SECRET`
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`

---

## 🧪 Scripts

In `client/`:

```bash
npm install
npm run dev
```

In each microservice:

```bash
go run cmd/api/main.go
```

---

## 🛠️ Tips

- Use `sql-migrate` or `golang-migrate` for handling migrations
- Store secrets in `.env` and do not push them to GitHub
- Use separate Dockerfiles for each microservice for better scaling
- Modularize Go code with `internal/` packages

---

## 📜 License

MIT
# go-microservices
