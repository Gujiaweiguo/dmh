# Backend Module Guide

**Module**: `backend/` — Go + go-zero REST API
**Entry**: `api/dmh.go`
**Port**: 8889

---

## Structure

```
backend/
├── api/
│   ├── dmh.go              # Main entry point
│   ├── dmh.api             # go-zero API spec
│   ├── etc/                # Configs (dmh-api.yaml, *.dev.yaml)
│   └── internal/
│       ├── handler/        # HTTP handlers (thin, parse→call logic)
│       ├── logic/          # Business logic (fat, contains real code)
│       ├── middleware/     # Auth, CORS, rate limiting
│       ├── svc/            # ServiceContext (DB, Redis)
│       └── types/          # Request/response types (generated)
├── model/                  # GORM models (User, Brand, Campaign, etc.)
├── common/                 # Shared utilities
│   ├── poster/             # Poster generation
│   ├── syncadapter/        # External DB sync
│   ├── utils/              # Helpers
│   └── wechatpay/          # WeChat Pay integration
├── migrations/             # SQL migrations
├── scripts/                # Deploy, test scripts
└── test/
    ├── integration/        # Integration tests (live API)
    └── performance/        # Benchmark tests
```

## Where to Look

| Task | Location |
|------|----------|
| Add new API endpoint | `api/dmh.api` → regenerate → `handler/`, `logic/`, `types/` |
| Modify business logic | `api/internal/logic/<module>/` |
| Add middleware | `api/internal/middleware/` |
| Database schema | `model/*.go` + `migrations/*.sql` |
| Fix auth/JWT | `api/internal/middleware/authmiddleware.go`, `logic/auth/` |
| Add test | Co-located `*_test.go` or `test/integration/` |

## Conventions

### Handler vs Logic
- **Handler**: Parse request, call logic, return response. No business logic.
- **Logic**: All business logic goes here. Handlers are thin wrappers.

### Error Handling
```go
// Use fmt.Errorf for errors with context
return nil, fmt.Errorf("failed to get user: %w", err)

// Logic returns errors, handlers use httpx.ErrorCtx
```

### GORM Patterns
```go
// Soft delete: WHERE deleted_at IS NULL (automatic with gorm.Model)
// Queries: Use scopes for reusable filters
db.Where("status = ?", "active").Find(&users)
```

### go-zero Code Generation
```bash
# After editing dmh.api:
cd backend/api
goctl api go -api dmh.api -dir .
```

## Anti-Patterns

| Avoid | Reason |
|-------|--------|
| Business logic in handlers | Handlers should be thin |
| Direct DB in handlers | Go through logic layer |
| `as any` type assertions | Use proper types |
| Ignoring errors | Always handle or wrap |
| Raw SQL in logic | Use GORM, or document why |

## Commands

```bash
# Development
go run api/dmh.go -f api/etc/dmh-api.yaml

# Build
go build -o dmh-api api/dmh.go

# Test all (MUST use -p 1 for database isolation)
go test -p 1 ./...

# Test with coverage
go test -p 1 ./... -coverprofile=coverage.out

# Integration tests (requires running API)
DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
go test ./test/integration/... -v -count=1

# Format
gofmt -w .
```

## Testing

### Test Commands Quick Reference

| Scenario | Command | Duration | Notes |
|----------|---------|----------|-------|
| Fast verify (PR pre-check) | `go test -p 1 ./... -short` | <30s | Skips integration tests |
| Full unit tests | `go test -p 1 ./...` | 2-3min | All unit tests |
| Integration tests | `DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1` | 5-10min | Requires running API |
| Coverage report | `go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out` | 3-5min | Threshold: ≥78% |
| Single test | `go test -v -run TestFunctionName ./path/to/package` | Varies | Debug specific test |

### Layered Test Strategy

| Layer | Responsibility | Mock Strategy | Test Focus |
|-------|----------------|---------------|------------|
| **Handler** | HTTP parse/response | Mock Logic | Request parsing, response format, HTTP codes |
| **Logic** | Business logic | Mock Repository | Business rules, edge cases, error handling |
| **Repository** | Data access | Real MySQL8 | SQL correctness, transactions, constraints |

### Test File Locations

```
backend/
├── api/internal/
│   ├── handler/
│   │   └── *_handler_test.go     # Handler tests (mock logic)
│   └── logic/
│       └── *_logic_test.go       # Logic tests (mock repo or real DB)
├── model/
│   └── *_test.go                 # Model tests (real DB)
└── test/
    ├── integration/              # Integration tests (live API)
    └── performance/              # Benchmark tests
```

### Database Isolation (CRITICAL)

**MUST use `-p 1` flag**:

```bash
go test -p 1 ./...
```

**Why**: Integration tests share `dmh_test` database. Parallel execution causes:
- Primary key conflicts (`Duplicate entry 'X' for key 'users.PRIMARY'`)
- Race conditions in test data setup

### Test Data Setup Pattern

Use "delete-then-create" for idempotent initialization:

```go
// CORRECT: Delete first, then create (idempotent)
func (s *UserRepoSuite) SetupTest() {
    s.db.Exec("DELETE FROM users WHERE id IN (1, 2, 3)")
    for _, user := range testUsers {
        s.db.Create(&user)
    }
}

// WRONG: Check-then-create (race condition)
if err := suite.db.Where("id = ?", user.Id).First(&existing).Error; err != nil {
    suite.db.Create(&user)  // Can fail with duplicate key
}
```

### SKIP Reason Tags

When tests must skip due to environment issues, use standard tags:

| Tag | Trigger |
|-----|---------|
| `API_UNAVAILABLE` | API service unreachable or timeout |
| `MYSQL_UNAVAILABLE` | MySQL connection failed |
| `REDIS_UNAVAILABLE` | Redis connection failed (dependency tests only) |
| `LOGIN_FAILED` | Login endpoint returned non-200 |
| `DATA_PREP_FAILED` | Test data preparation failed |

Example:
```go
if !isAPIAvailable() {
    t.Skip("SKIP: API_UNAVAILABLE - API service not running")
}
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DMH_INTEGRATION_BASE_URL` | `http://localhost:8889` | API base URL |
| `DMH_TEST_ADMIN_USERNAME` | `admin` | Test admin username |
| `DMH_TEST_ADMIN_PASSWORD` | `123456` | Test admin password |

### Coverage Threshold

- **Current**: ~68%
- **Target**: ≥78%

Check coverage:
```bash
go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | tail -1
```

### Performance Test Split

| Gate | Command | Duration | Notes |
|------|---------|----------|-------|
| PR Smoke | `go test -short -bench=. -benchtime=100ms ./test/performance/...` | <1s | No backend needed |
| Nightly Full | `go test -v ./test/performance/...` | ~1min | Requires backend |

### Test Anti-Patterns

| Avoid | Reason | Correct Approach |
|-------|--------|------------------|
| Business logic in handlers | Breaks layering | Put in Logic layer |
| Integration tests with mock data | Loses integration value | Use real database |
| Shared mutable state between tests | Test interference | Independent init/cleanup per test |
| Check-then-create pattern | Race condition | Delete-then-create (idempotent) |
| `t.Skip()` to hide failures | Hides real defects | Fix or mark as known issue |
| Infinite retries for flaky tests | Masks instability | Root cause analysis, fix or isolate |

### Pre-Commit Checklist

Before committing:

- [ ] `go test -p 1 ./... -short` passes
- [ ] Coverage ≥78%
- [ ] No `t.Skip()` hiding failures
- [ ] No business logic in handlers
- [ ] Integration tests use real DB (not mocks)
- [ ] Test data is idempotent (delete-then-create)
## Database Migrations

```bash
# Create migration
echo "-- Migration: Description" > migrations/$(date +%Y%m%d)_description.sql

# Run migration
docker exec -i mysql8 mysql -uroot -p'Admin168' dmh < migrations/xxx.sql
```

## Config Files

| File | Purpose |
|------|---------|
| `api/etc/dmh-api.yaml` | Default config (Docker) |
| `api/etc/dmh-api.dev.yaml` | Local development |
| `api/etc/dmh-api.docker.yaml` | Docker-specific |

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| go-zero | 1.6.0 | REST framework |
| GORM | 1.25.5 | ORM |
| golang-jwt | 4.5.0 | Auth tokens |
| testify | 1.10.0 | Testing |
