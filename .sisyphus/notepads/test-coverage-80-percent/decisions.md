# Test Coverage 80% - Decisions

## 架构决策

### 1. 覆盖率目标
- 后端: 80% (当前 45.6%)
- 前端 Admin: 70% (考虑 UI 测试难度)
- 前端 H5: 维持 80%+

### 2. 测试优先级
1. auth (认证核心)
2. order (订单核心)
3. campaign (活动核心)
4. reward/distributor (分销核心)

### 3. 测试模式
- Handler: sqlite in-memory + httptest
- Views: Vitest + jsdom + mock API

## Guardrails

- ❌ 禁止重构生产代码
- ❌ 禁止为自动生成代码写测试
- ❌ 禁止创建新测试框架
