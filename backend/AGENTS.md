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

# Test all
go test ./...

# Test with coverage
go test ./... -coverprofile=coverage.out

# Integration tests (requires running API)
DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
go test ./test/integration/... -v -count=1

# Format
gofmt -w .
```

## Testing

### Unit Tests
- Co-located: `logic/foo_logic_test.go` next to `logic/foo_logic.go`
- Framework: `github.com/stretchr/testify`
- Pattern: Suite-based with `suite.Suite`

### Integration Tests
- Location: `test/integration/`
- Requires: Running API at `localhost:8889`
- Env vars: `DMH_INTEGRATION_BASE_URL`, `DMH_TEST_ADMIN_USERNAME`, `DMH_TEST_ADMIN_PASSWORD`

### Coverage Target
- Current: ~68%
- Target: 70%+

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
