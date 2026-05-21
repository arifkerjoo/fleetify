# Fleetify Backend

Fleetify adalah sistem internal manajemen pemeliharaan armada kendaraan berbasis REST API menggunakan Golang + Fiber + MySQL

Backend ini mendukung:

* Authentication berbasis JWT HttpOnly Cookie
* Role-based access control
* Workflow maintenance report
* Upload foto kendaraan
* Approval management
* Pagination & filtering
* Clean architecture structure
* Dockerized environment

---

# Features

## Authentication

* Register
* Login
* Logout
* Session validation
* JWT HttpOnly Cookie

---

## Vehicle Management

* Create vehicle
* Update vehicle
* Delete vehicle
* List vehicles

---

## Master Items

Digunakan untuk part/jasa estimasi maintenance.

* Create item
* Update item
* Delete item
* List items

---

## Maintenance Report Workflow

### Service Advisor (SA)

* Create maintenance report
* Upload initial vehicle photo
* Add estimated parts/services
* Complete maintenance report
* Upload proof photo

### Approval / Management

* Review incoming reports
* Approve maintenance report

---

# Tech Stack

| Layer        | Technology         |
| ------------ | ------------------ |
| Language     | Go                 |
| Framework    | Fiber v3           |
| Database     | MySQL 8            |
| Auth         | JWT                |
| Container    | Docker             |
| ORM          | GORM               |
| Validation   | Validator          |
| Architecture | Clean Architecture |

---

# Project Structure

```bash
backend/
│
├── cmd/
├── internal/
│   ├── domain/
│   ├── interfaces/
│   ├── repository/
│   ├── usecase/
│   └── infrastructure/
│
├── routes/
├── uploads/
├── utils/
├── .env
├── Dockerfile
├── go.mod
└── main.go
```

---

# API Base URL

```bash
http://localhost:8080/api/v1
```

---

# Roles

| Role     | Description          |
| -------- | -------------------- |
| SA       | Service Advisor      |
| APPROVAL | Management Approval  |



# Environment Variables

Buat file `.env`:

```env
APP_NAME=Fleetify

APP_PORT=8080

DB_HOST=mysql
DB_PORT=3306
DB_NAME=fleetify
DB_USER=fleetify
DB_PASSWORD=fleetify123

JWT_SECRET=super-secret-key
JWT_EXPIRE_HOURS=24

COOKIE_SECURE=false
COOKIE_HTTP_ONLY=true
COOKIE_SAME_SITE=Lax
```

---

# Running with Docker

## Start Services

```bash
docker compose up --build
```

---

## Services

| Service     | Port |
| ----------- | ---- |
| Backend API | 8080 |
| MySQL       | 3306 |

# Seed Data

Backend mendukung seed awal:

```bash
./main seed
```

Seed biasanya berisi:

* dummy users
* vehicles
* master items

---

# Authentication

Backend menggunakan:

* JWT
* HttpOnly Cookie
* Secure Cookie support
* SameSite policy

# API Endpoints

# Auth

| Method | Endpoint       |
| ------ | -------------- |
| POST   | /auth/register |
| POST   | /auth/login    |
| POST   | /auth/logout   |
| GET    | /auth/me       |

---

# Vehicles

| Method | Endpoint      |
| ------ | ------------- |
| GET    | /vehicles     |
| GET    | /vehicles/:id |
| POST   | /vehicles     |
| PUT    | /vehicles/:id |
| DELETE | /vehicles/:id |

---

# Master Items

| Method | Endpoint          |
| ------ | ----------------- |
| GET    | /master-items     |
| GET    | /master-items/:id |
| POST   | /master-items     |
| PUT    | /master-items/:id |
| DELETE | /master-items/:id |

---

# Reports

| Method | Endpoint                   |
| ------ | -------------------------- |
| GET    | /reports                   |
| GET    | /reports/:id               |
| POST   | /reports                   |
| POST   | /reports/:id/initial-photo |
| PATCH  | /reports/:id/approve       |
| PATCH  | /reports/:id/complete      |


# Author

Fleetify Backend
Built with Go + Fiber
