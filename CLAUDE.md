# CLAUDE.md

## Project: KZ Salary Calculator

Full-stack web application for salary calculation in Kazakhstan.

Stack:

- Go
- PostgreSQL
- HTML + Tailwind + JS
- Docker Compose
- REST API

Claude Code MUST follow rules below.

---

## Architecture

Use layered architecture:

cmd/server/main.go

internal/
    handlers/calculate_handler.go
            └──history_handler.go
            └──mci_handler.go
    services/salary_service.go
    repository/postgres.go
            └──calculation_repo.go
    models/calculation.go
    calculator/calculator.go
    middleware/
    config/config.go
web/
│   ├── index.html
│   ├── app.js
│   └── style.css
├── migrations/
│   └── 001_create_table.sql
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
├── CLAUDE.md
Rules:

- handlers → HTTP only
- services → business logic
- repository → DB logic
- calculator → formulas
- middleware → logging / recovery
- config → env / constants

Do NOT put logic in handlers.

Service layer is required.

Repository must use interfaces.

main.go must stay clean.

Do not put initialization logic in main.

All startup logic must be placed in separate file (app.go / server.go).

main.go should only call Run().

---

## Logging

Use logging for:

- every HTTP request
- errors
- DB errors
- server start/stop

Create middleware:

internal/middleware/logging.go

Log:

method
path
status
duration

---

## HTTP rules

Use correct status codes:

POST → 201
GET → 200
DELETE → 204
NOT FOUND → 404
BAD REQUEST → 400

Must check sql.ErrNoRows → return 404

Never return request JSON as response.
Return DB result.

---

## Graceful shutdown

Server must support graceful shutdown.

Use context + signal.

Required.

---

## Validation

Validate input:

salary > 0

mode = gross | net

Return 400 if invalid.

---

## Salary rules

OPV = 10%
VOSMS = 2%
IPN = 10%

IPN base:

gross - OPV - VOSMS - 14*MCI

Employer:

SO = 3.5%
OOSMS = 3%
SN = 9.5%

MCI configurable.

Must support:

gross → net
net → gross

---

## Repository rules

Repository must use interface.

Example:

type CalculationRepository interface {
    Save(...)
    GetHistory(...)
}

Handler must NOT use SQL.

Service must use repository.

---

## Database rules

Use PostgreSQL.

Use DATE or TIMESTAMP.

NOT TEXT.

Table: calculations

Add indexes.

Example:

CREATE INDEX idx_created_at
ON calculations(created_at);

---

## History

Return last 10.

Support limit / offset.

Pagination required.

---

## SQL rules

Do calculations in SQL if possible.

Do not aggregate in Go.

Use correct WHERE.

No LIKE when exact match required.

---

## Config

Use .env

Must have .env.example

Fields:

DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
MCI

Use config package.

---

## Docker

Must start with:

docker compose up

Services:

app
postgres

---

## Frontend

Use:

Tailwind
Vanilla JS

Must show:

salary
taxes
employer
history

Responsive.

---

## Frontend UX Improvements

- Use Tailwind CSS for modern UI
- Responsive design with grid/flex layout
- Card-style sections for:
  - Form input
  - Employee Deductions
  - Employer Contributions
  - History
- Highlight Net Salary in blue
- Highlight Employer Contributions in green
- Hover and transition effects on cards
- Responsive tables for history
- Input fields with focus styles and validation
- Vanilla JS for API calls
- Show loading spinner while fetching
- Show error messages if API fails

---

## Git

Atomic commits.

Examples:

feat: add calculator
feat: add repository
feat: add logging middleware
fix: correct IPN formula
feat: graceful shutdown

---

## When generating code

Claude MUST:

- use service layer
- use interfaces
- use middleware logging
- check sql.ErrNoRows
- use correct HTTP codes
- use graceful shutdown
- use env config
- not put SQL in handler
- not put logic in handler
- not use TEXT for dates
- add indexes
- support pagination