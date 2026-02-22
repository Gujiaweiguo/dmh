## ADDED Requirements

### Requirement: Backend test sequential execution for database isolation
后端测试 SHALL 使用顺序执行模式（`-p 1`）以避免共享测试数据库的竞争条件。

#### Scenario: Run backend tests sequentially
- **WHEN** 执行后端单元测试和集成测试
- **THEN** 系统 SHALL 使用 `go test -p 1 ./...` 命令
- **AND** 测试 SHALL 按包顺序执行以避免数据库竞争

#### Scenario: CI workflow enforces sequential test execution
- **WHEN** CI/CD 流水线运行后端测试
- **THEN** 工作流 SHALL 配置为使用 `-p 1` 参数
- **AND** 工作流 SHALL 明确注释使用顺序执行的原因

### Requirement: Idempotent test data setup
测试数据初始化 SHALL 使用幂等操作（先删除再创建），避免重复运行时的主键冲突。

#### Scenario: Test data cleanup before creation
- **WHEN** 测试套件的 `SetupTest` 或 `createTestData` 执行
- **THEN** 数据初始化 SHALL 先执行 DELETE 操作清除旧数据
- **AND** 然后 SHALL 执行 CREATE 操作创建新数据

#### Scenario: Avoid check-then-create race condition
- **WHEN** 多个测试并发初始化共享数据库
- **THEN** 代码 SHALL NOT 使用 "检查再创建" 模式（SELECT + INSERT）
- **AND** 代码 SHALL 使用 "先删除再创建" 模式（DELETE + INSERT）
