# Regression Baseline Summary

Date: 2026-02-23
Scope: Post-fix regression baseline after auth/test/docs updates

## 1) Basic Information

- Run ID: `2026-02-23-run001`
- Execution time: 2026-02-23 10:10 - 10:16 (UTC+8)
- Executor: AI Agent
- Branch: `main`
- Commit: `61a3fe6ba02903101b73816f17aa1ac53891de98`

## 2) Execution Notes

- `make up` could not be used in this environment due to Docker permission (`make: docker: Permission denied`).
- Fallback method: start backend and frontend services with local process mode for regression execution.

## 3) Commands and Results

### 3.1 Integration Suite

Command:

```bash
make test-integration
```

Result:

- Status: PASS
- Key marker: `PASS` for `go test ./test/integration/... -v -count=1`
- Coverage note: this run validates integration behavior, not coverage gate.

### 3.2 Order Regression

Command:

```bash
make test-order-regression
```

Result:

- Status: PASS
- Key marker: `[order-mysql8-regression] PASS`
- Included checks: `TestOrderVerifyRoutesAuthGuard`, `TestOrderCreateDuplicateMessage`

### 3.3 Frontend E2E (Headless)

Command:

```bash
make test-e2e-headless
```

Result:

- Admin E2E: PASS (`24/24`)
- H5 E2E: PASS (`7/7`)

## 4) Related Change Scope

Validated against the following merged changes:

- `backend/api/internal/handler/routes.go`: restore JWT guard for distributor routes
- `backend/api/dmh.go`: remove temporary debug prints in poster static route
- `backend/api/internal/handler/brand/handler_test.go`: enable and stabilize body-parse coverage tests
- `docs/H5_PENDING_FEATURES.md`: add H5 pending features tracking list

## 5) Evidence Pointers

- Backend local run log: `/tmp/dmh-api-regression.log`
- Session command outputs: current CLI session logs (integration/order/e2e)

## 6) Conclusion

This regression baseline run is GREEN:

- Integration suite passed
- Order regression passed
- Admin/H5 E2E passed

Release risk related to this change set is currently low under the validated matrix above.
