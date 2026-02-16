# Run Summary - 2026-02-16-run001

## Scope

- Environment: local runtime (`http://localhost:8889`) + workspace test commands
- Test groups:
  - Backend unit: `cd backend && go test ./... -v`
  - Backend integration: `cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 DMH_TEST_ADMIN_USERNAME=admin DMH_TEST_ADMIN_PASSWORD=123456 go test ./test/integration/... -v -count=1`
  - Frontend Admin unit: `cd frontend-admin && npm run test`
  - Frontend H5 unit: `cd frontend-h5 && npm run test`
  - OpenSpec strict validation: `openspec validate --all --strict --no-interactive`

## Result

- Backend unit: PASS
- Backend integration: PASS（`ok dmh/test/integration 4.130s`，`SKIP=0`）
- Frontend Admin unit: PASS（`34 passed / 314 tests`）
- Frontend H5 unit: PASS（`54 passed / 985 tests`）
- OpenSpec strict validation: PASS（`7 passed, 0 failed`）

## Notes

- 初次集成测试出现登录失败（`用户名或密码错误`）导致大量跳过。
- 处置动作：在 `backend/` 重编译 `dmh-api`，替换 `deploy/dmh-api` 并重启 `dmh-api` 容器后重跑。
- 重跑后结果：集成测试 `PASS` 且 `SKIP=0`。
- Release gate impact: **PASS**。

## Evidence Files

- `backend-unit.log`
- `backend-integration.log`
- `frontend-admin-unit.log`
- `frontend-h5-unit.log`
- `openspec-validate.log`
