## 1. Documentation Updates

- [x] 1.1 Fix `distributor_integration_test.go` data cleanup logic (先删除再创建模式)
- [x] 1.2 Update `backend/AGENTS.md` with test execution guidelines
- [x] 1.3 Update root `AGENTS.md` with test execution command reference

## 2. CI/CD Configuration

- [x] 2.1 Review existing GitHub workflows for test commands
- [x] 2.2 Update workflows to use `go test -p 1 ./...` for backend tests
  - coverage-gate.yml ✅
  - stability-checks.yml ✅
  - feedback-guard.yml ✅
- [ ] 2.3 Verify CI pipeline runs successfully with sequential tests

## 3. Spec Update

- [x] 3.1 Add backend test parallelism requirement to `system-test-execution` spec
- [x] 3.2 Validate spec changes with `openspec validate`

## 4. Verification

- [x] 4.1 Run `go test -p 1 ./...` and verify all tests pass
- [ ] 4.2 Confirm CI workflow passes with new configuration
