# 🤖 Development Guidelines: Using Claude Code

## Project: KZ Salary Calculator

A production-ready full-stack salary calculation application demonstrating solid Go backend architecture, PostgreSQL database design, and modern frontend development.

**Stack:**
- Go 1.25 backend
- PostgreSQL 16 database  
- HTML5 + Vanilla JavaScript + Tailwind CSS 3.3 frontend
- Docker Compose for local development and deployment
- REST API with proper error handling and validation

---

## 🎯 When to Use Claude Code

###Use Claude for:
- **Adding new API endpoints** (new handlers, services, DB functions)
- **Fixing bugs** in calculation logic, database queries, or API behavior
- **Improving performance** (adding indexes, optimizing queries, connection pooling)
- **Adding tests** for new features or fixing broken tests
- **Frontend enhancements** (new UI sections, validation, error handling)
- **Configuration issues** (environment variables, docker-compose setup)
- **Documentation** (README, code comments, inline docs)

### Example prompts:

```
"Add a PUT endpoint to /api/calculation/:id that updates an existing calculation record"

"The CalculateFromNet function sometimes doesn't converge. 
Review the algorithm and add proper convergence checking with max iterations"

"Add unit tests for the calculator package covering both gross and net mode calculations"

"Optimize the /api/history endpoint - add indexes and pagination benchmarking"
```

---

## 🏗 Architecture Rules (MUST FOLLOW)

### Directory Structure

```
cmd/server/
├── main.go              # Entry point - only calls Run() from app.go
└── app.go               # Initialization, server setup, graceful shutdown

internal/
├── handlers/            # HTTP layer - request parsing and response formatting
│   ├── calculate_handler.go
│   ├── history_handler.go
│   ├── mci_handler.go
│   └── helpers.go       # Shared handler utilities
├── services/            # Business logic layer - calculations, validation
│   └── salary_service.go
├── repository/          # Data access layer - database operations
│   ├── calculation_repo.go  # Interface definition
│   └── postgres.go          # PostgreSQL implementation
├── models/              # Domain models
│   └── calculation.go
├── calculator/          # Tax calculation formulas (pure functions)
│   ├── calculator.go
│   └── calculator_test.go
├── middleware/          # HTTP middleware
│   └── logging.go
└── config/              # Configuration and validation
    └── config.go

web/                    # Frontend assets (served as static files)
├── index.html
├── app.js
└── style.css

migrations/             # Database schema and version control
└── 001_create_table.sql

```

### Layer Responsibilities

| Layer | Responsibility | Examples |
|-------|---|---|
| **handlers/** | HTTP only | Parse request JSON, call service, return HTTP response with proper status codes |
| **services/** | Business logic | Validate inputs, perform calculations, orchestrate repository calls |
| **repository/** | Database access | Execute SQL queries, map rows to models, handle db errors |
| **calculator/** | Pure functions | Tax formulas, math operations, no side effects |
| **config/** | Configuration | Load environment variables, validate ranges, return errors |

### Critical Rules

1. **NO LOGIC IN HANDLERS**
   - ❌ Don't calculate taxes in handlers
   - ❌ Don't validate salary range in handlers  
   - ✅ Parse request, call service, return response

2. **SERVICE LAYER IS REQUIRED**
   - ALL business logic must be in `internal/services/`
   - Handlers only call services
   - Repository only called by services

3. **REPOSITORY MUST USE INTERFACES**
   - Define `CalculationRepository` interface
   - Handlers must NOT use SQL directly
   - Easy to swap PostgreSQL for another database

4. **main.go MUST STAY CLEAN**
   - Only imports and a call to Run()
   - All initialization in app.go
   - Example:
     ```go
     func main() {
         if err := Run(); err != nil {
             log.Fatal(err)
         }
     }
     ```

5. **No initialization logic in main**
   - Database connections → app.go
   - Route registration → app.go
   - Middleware setup → app.go

---

## 📝 HTTP Rules (MUST FOLLOW)

### Status Codes

Always use correct HTTP status codes:

| Method | Success | Examples |
|--------|---------|----------|
| **POST** | `201` | Create calculation |
| **GET** | `200` | Get history, get MCI |
| **PUT** | `200` | Update calculation |
| **DELETE** | `204` | Delete calculation |
| **Error: Not Found** | `404` | Calculation by ID not found |
| **Error: Bad Request** | `400` | Invalid salary, invalid mode |
| **Error: Request Too Large** | `413` | Payload exceeds 1MB |
| **Error: Server Error** | `500` | Database error, unexpected panic |

### Error Handling

```go
// Handle sql.ErrNoRows - return 404
if err == sql.ErrNoRows {
    return c.JSON(404, map[string]string{"error": "calculation not found"})
}

// Validation errors - return 400
if salary <= 0 {
    return c.JSON(400, map[string]string{"error": "salary must be > 0"})
}

// Database errors - return 500
if err != nil {
    log.Printf("database error: %v", err)
    return c.JSON(500, map[string]string{"error": "internal server error"})
}
```

### Response Format

```go
// CORRECT - Return database result
calculation := repo.GetByID(id)
return c.JSON(200, calculation)

// WRONG - Don't return request input
input := request.JSON()
return c.JSON(201, input)  // ❌ BAD

// CORRECT - Return calculated result
result := service.Calculate(input)
return c.JSON(201, result)  // ✅ GOOD
```

---

## 🏃 Graceful Shutdown (REQUIRED)

All changes to `app.go` must preserve graceful shutdown:

```go
func Run() error {
    // ... setup ...
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    
    go startServer(server, addr)
    
    <-quit  // Wait for shutdown signal
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    return server.Shutdown(ctx)  // Graceful shutdown
}
```

**Key points:**
- Listen for OS signals (SIGINT, SIGTERM)
- Use context.WithTimeout for graceful shutdown (max 30s)
- Allow requests to complete before shutting down

---

## ✅ Input Validation (REQUIRED)

Validate all user inputs in the service layer:

```go
// In services/salary_service.go
func (s *SalaryService) Calculate(req CalculateRequest) (*Calculation, error) {
    // Validate salary amount
    if req.Salary <= 0 {
        return nil, errors.New("salary must be > 0")
    }
    
    // Validate mode
    if req.Mode != "gross" && req.Mode != "net" {
        return nil, errors.New("mode must be 'gross' or 'net'")
    }
    
    // Proceed with calculation
    return s.calculator.CalculateFromGross(req.Salary)
}
```

### Validation Rules

| Field | Rule | Error Message |
|-------|------|---|
| `salary` | Must be > 0 | `"salary must be > 0"` |
| `salary` | Max 100,000,000 | `"salary exceeds maximum"` |
| `mode` | Must be "gross" or "net" | `"mode must be 'gross' or 'net'"` |

---

## 💰 Kazakhstan Tax System Rules

### Employee Deductions

```
ОПВ (Pension):      gross * 0.10
ВОСМС (Health):     gross * 0.02
ИПН (Income Tax):   (gross - OПВ - ВОСМС - 14*MCI) * 0.10

Net Salary = gross - ОПВ - ВОСМС - ИПН
```

### Employer Contributions

```
СО (Social):        (gross - ОПВ) * 0.035
ООСМС (Employer):   gross * 0.03
СН (Social Tax):    (gross - ОПВ - ВОСМС) * 0.095 - СО

Employer Total = gross + СО + ООСМС + СН
```

### MCI (Minimum Pension Amount)

- Current: 3,932 KZT (2025)
- Configurable via `MCI` environment variable
- Used in ИПН calculation: maximum IPN base is 14 × MCI
- **Update annually** when government announces new rate

### Key Note: CalculateFromNet

The reverse calculation (NET → GROSS) uses iterative approach:

```go
const MaxIterations = 20
const Epsilon = 0.01  // 1 tenge tolerance

// Converges when: |calculated_net - target_net| < Epsilon
// Must return error if convergence not found after MaxIterations
```

---

## 📊 Database Rules

### Schema Requirements

Use `NUMERIC(15,2)` for all monetary values:

```sql
CREATE TABLE calculations (
    id BIGSERIAL PRIMARY KEY,
    gross_salary NUMERIC(15,2) NOT NULL,
    net_salary NUMERIC(15,2) NOT NULL,
    opv NUMERIC(15,2) NOT NULL,
    vosms NUMERIC(15,2) NOT NULL,
    ipn NUMERIC(15,2) NOT NULL,
    so NUMERIC(15,2) NOT NULL,
    oosms NUMERIC(15,2) NOT NULL,
    sn NUMERIC(15,2) NOT NULL,
    employer_total NUMERIC(15,2) NOT NULL,
    mode VARCHAR(5) NOT NULL,  -- 'gross' or 'net'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Indexed Queries

```sql
-- History queries are slow without this
CREATE INDEX idx_created_at ON calculations(created_at DESC);

-- Filtering by mode
CREATE INDEX idx_mode ON calculations(mode);

-- Combined filter + sort
CREATE INDEX idx_mode_created_at ON calculations(mode, created_at DESC);
```

### SQL Best Practices

- Use DATE or TIMESTAMP, NEVER TEXT for dates
- Use prepared statements (prevents SQL injection)
- Use `sql.ErrNoRows` to check "not found"
- Do calculations in SQL when possible
- Use exact matches, not LIKE for exact matching

### When Adding New Queries

1. **Always use parameterized queries** (prevents SQL injection)
   ```go
   rows, err := db.Query("SELECT * FROM calculations WHERE id = ?", id)
   ```

2. **Check sql.ErrNoRows for 404s**
   ```go
   if err == sql.ErrNoRows {
       return nil, 404  // Not found
   }
   ```

3. **Add indexes for frequently queried fields**
   ```go
   // If you add a filter, add an index
   CREATE INDEX idx_field_name ON table_name(field_name);
   ```

---

## 🧪 Testing

### Test Organization

Tests live alongside code they test:

- `internal/calculator/calculator.go` → `internal/calculator/calculator_test.go`
- `internal/services/salary_service.go` → `internal/services/salary_service_test.go`

### What to Test

1. **Calculator**: Math accuracy, convergence, edge cases
   ```go
   func TestCalculateFromGross(t *testing.T) {
       gross := 500000.0
       calc, _ := CalculateFromGross(gross)
       if calc.NetSalary < 0 {
           t.Error("net salary cannot be negative")
       }
   }
   ```

2. **Services**: Input validation, business logic
   ```go
   func TestValidationError(t *testing.T) {
       _, err := service.Calculate(CalculateRequest{Salary: -100})
       if err == nil {
           t.Error("should reject negative salary")
       }
   }
   ```

3. **Integration**: Full request → response flow (use docker compose)

### Running Tests

```bash
# All tests
go test ./...

# Single package
go test ./internal/calculator

# With coverage
go test -cover ./...

# Verbose output
go test -v ./internal/calculator
```

---

## 🔧 Configuration Management (REQUIRED)

### Environment Variables

Load via `internal/config/config.go`:

```go
type Config struct {
    DBHost     string
    DBPort     int      // Must validate: 1-65535
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string   // require, disable, prefer
    ServerPort int      // Must validate: 1-65535
    MCI        float64
}

func Load() (*Config, error) {
    // Must validate all fields
    // Must return error if invalid
}
```

### Required Variables

| Variable | Default | Validation |
|----------|---------|---|
| `DB_HOST` | `localhost` | Non-empty string |
| `DB_PORT` | `5432` | Integer 1-65535 |
| `DB_USER` | `postgres` | Non-empty string |
| `DB_NAME` | `salary_db` | Non-empty string |
| `DB_SSL_MODE` | `require` | One of: require, disable, prefer |
| `SERVER_PORT` | `8080` | Integer 1-65535 |
| `MCI` | `3932` | Positive number |

### Docker Compose Configuration

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_SSL_MODE: disable  # OK for local, require for production
    depends_on:
      - postgres
```

---

## 📋 Logging (REQUIRED)

Log all important events in `internal/middleware/logging.go`:

```go
// HTTP Request Logging
LOG: method path status duration client_ip
POST /api/calculate 201 125ms 127.0.0.1

// Error Logging
LOG: error_type error_message context
ERROR: database "connection refused" host=postgres port=5432
ERROR: validation "salary must be > 0" salary=-100
```

### Where to Log

- Every HTTP request (method, path, status, duration)
- All errors (with context)
- Database errors (with query context)
- Server start/stop events

### Logging NOT to do

- Don't log passwords
- Don't log sensitive data
- Don't log entire request bodies in production
- Don't log with fmt.Println (use proper logging)

---

## 🐳 Docker & Docker Compose

### Multi-stage Build

The Dockerfile uses multi-stage build to keep images small:

```dockerfile
# Stage 1: Build
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

# Stage 2: Runtime
FROM debian:bookworm-slim
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
```

### Docker Compose Services

```yaml
version: '3'
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: salary_db
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_SSL_MODE: disable
    depends_on:
      - postgres
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3
```

### Starting for Development

```bash
docker compose up --build
```

This starts both PostgreSQL and the application, with automatic migration.

---

## ✨ Frontend Guidelines

Use vanilla JavaScript (no frameworks):

```javascript
// Good: Direct API calls with proper error handling
async function calculate() {
    try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 10000);
        
        const response = await fetch('/api/calculate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ salary, mode }),
            signal: controller.signal
        });
        
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        return await response.json();
    } catch (error) {
        if (error.name === 'AbortError') {
            showError('Request timeout');
        } else {
            showError('Network error');
        }
    }
}

// Good: Use events for interactivity
document.getElementById('calculate-btn').addEventListener('click', calculate);

// Good: Update DOM based on response
function showResult(data) {
    document.getElementById('net-salary').textContent = 
        data.net_salary.toFixed(2);
}
```

---

## 🚀 When Adding New Features

### Step 1: Plan the architecture
- Where does it fit? (handler, service, repository, calculator?)
- What database changes needed?
- What validation is required?

### Step 2: Database first
- Write migration: `migrations/XXX_add_feature.sql`
- Test with: `psql -d salary_db -f migrations/XXX_add_feature.sql`

### Step 3: Backend in layers
1. Update `models/` if needed
2. Update `repository/` interface and postgres implementation
3. Add `services/` business logic
4. Add `handlers/` HTTP endpoint
5. Add request validation

### Step 4: Tests
- Add unit tests for calculator/services
- Test with curl before frontend

### Step 5: Frontend
- Update HTML/CSS
- Add JavaScript functions
- Test all error cases

### Example: Add rating feature

```
1. Migration: ALTER TABLE calculations ADD COLUMN rating INT;
2. Model: Add Rating field to Calculation
3. Repository: Update Save, GetByID, GetHistory to handle rating
4. Service: Validate rating is 1-5
5. Handler: Accept rating in GET request
6. Tests: Test rating validation
7. Frontend: Add star rating UI
```

---

## 🔍 Common Issues & Solutions

### Issue: CalculateFromNet doesn't converge
**Cause:** Too few iterations or epsilon too large
**Solution:** Increase MaxIterations, check formula logic
```go
const MaxIterations = 20  // Increased from 10
const Epsilon = 0.01       // 1 tenge tolerance
```

### Issue: Slow history queries
**Cause:** Missing index on created_at
**Solution:** Add index
```sql
CREATE INDEX idx_created_at ON calculations(created_at DESC);
```

### Issue: Database connection errors in Docker
**Cause:** App starts before PostgreSQL is ready
**Solution:** Use depends_on + healthcheck in docker-compose.yml

### Issue: Configuration validation fails
**Cause:** Invalid port or SSL mode
**Solution:** Check .env file, valid ports are 1-65535, valid SSL modes are: require, disable, prefer

---

## 📚 Code Review Checklist

When reviewing code changes, ensure:

- [ ] Business logic is in services, NOT handlers
- [ ] SQL queries are in repository, NOT handlers
- [ ] Input validation is in services
- [ ] Proper HTTP status codes (201 for POST, 400 for validation, 404 for not found)
- [ ] Error handling with context
- [ ] No hardcoded values (use environment variables)
- [ ] Database migrations for schema changes
- [ ] Tests for new logic
- [ ] No SQL injection (using parameterized queries)
- [ ] Graceful shutdown NOT broken
- [ ] Configuration changes documented in README

---

## 📖 Git Commit Messages

Use atomic commits with clear messages:

```
feat: add calculation rating feature
fix: correct IPN formula calculation
feat: add timeout handling to frontend
test: add calculator convergence tests
docs: update README with API examples
chore: update dependencies
refactor: extract tax calculation helper function
```

**Format:** `<type>: <description>`

Types: `feat`, `fix`, `test`, `docs`, `chore`, `refactor`, `perf`

---

## ❓ Getting Help

When stuck, provide Claude with:

1. **What you want to do:** "Add a new endpoint to update calculations"
2. **Where it fits:** "It should be in handlers/calculate_handler.go and services/salary_service.go"
3. **What you've tried:** "I created the handler but don't know how to add the service method"
4. **Current error:** If there's an error, paste it completely

Example good prompt:
```
I need to add a PUT /api/calculations/:id endpoint to update an existing calculation.
Following the layered architecture, I need to:
1. Update the repository to include an Update method
2. Add a service method that validates the update
3. Create a handler that parses the ID and calls the service

Can you help me start with the repository interface?
```

---

## 🎓 Learning from This Project

This codebase demonstrates:
- ✅ Clean architecture with separation of concerns
- ✅ Interface-based design for flexibility
- ✅ Proper error handling at each layer
- ✅ Configuration management
-✅ Database design with proper indexing
- ✅ Graceful shutdown
- ✅ Comprehensive HTTP error handling
- ✅ Test coverage for critical logic
- ✅ Docker containerization

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