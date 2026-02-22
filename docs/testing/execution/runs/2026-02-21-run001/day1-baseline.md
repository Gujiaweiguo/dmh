# Day 1 Baseline Report

Date: 2026-02-21
Scope: DMH Day 1 baseline checks (environment, unit tests, lint/format)

## 1) Environment Status

- dmh-api: Up
- dmh-nginx: Up
- mysql8: Up (healthy)
- redis-dmh: Up (healthy)
- Admin frontend (`http://localhost:3000`): HTTP 200
- H5 frontend (`http://localhost:3100`): HTTP 200

Notes:
- `http://localhost:8889/api/health` did not return success in this run, but API container stayed healthy and frontend endpoints were reachable via nginx.

## 2) Test Results

### 2.1 Backend

Command baseline used:

```bash
cd backend && go test ./...
```

Observed:
- Full backend suite did not pass due to failures in some handler/integration-related packages.
- Re-check of selected packages:
  - PASS: `dmh/api/internal/handler/admin`
  - PASS: `dmh/api/internal/handler/auth`
  - FAIL: `dmh/api/internal/handler/brand`
  - FAIL: `dmh/api/internal/handler/distributor`
  - PASS: `dmh/api/internal/handler/feedback`
  - PASS: `dmh/api/internal/handler/role`

Representative failures:
- `TestGetBrandsHandler_Pagination` count assertion mismatch in `backend/api/internal/handler/brand/handler_test.go`.
- Multiple FK constraint failures in distributor tests (missing/invalid referenced user/brand/campaign rows) in `backend/api/internal/handler/distributor/handler_coverage_test.go` and `backend/api/internal/handler/distributor/handler_test.go`.

### 2.2 Frontend Admin

Command:

```bash
cd frontend-admin && npm run test
```

Result:
- PASS: 41/41 test files, 357/357 tests.

### 2.3 Frontend H5

Command:

```bash
cd frontend-h5 && npm run test
```

Result:
- PASS: 54/54 test files, 988/988 tests.

## 3) Code Quality (Lint/Format)

Command:

```bash
make check
```

Observed:
- `frontend-admin` lint: PASS
- `frontend-h5` lint: PASS
- `backend` `gofmt -d .` produced diffs (format drift exists, especially in test files and `backend/model/member.go`).

## 4) Baseline Conclusion

Day 1 baseline is partially complete:
- Environment and both frontend unit suites are healthy.
- Backend suite is not green due to existing failing tests in brand/distributor handlers.
- Lint is clean; Go formatting drift exists and should be normalized before release gating.

## 5) Next-Day Focus (Day 2)

1. Stabilize backend handler tests:
   - Fix pagination expectation/data isolation for brand handler tests.
   - Fix distributor test fixture setup order for FK dependencies.
2. Re-run:

```bash
cd backend && go test ./api/internal/handler/brand ./api/internal/handler/distributor
```

3. Then run Day 2 core path checks:

```bash
make test-integration
make test-order-regression
```
