# Quickstart: 平台 SRE 与规模化治理

## 1. 实施前准备

1. 确认当前分支为 `011-sre-scale`，禁止在 `main` 或 `master` 直接开发。
2. 先完成数据库备份，并在 `artifacts/011-sre-scale/backup-manifest.txt` 记录命令、时间戳、产物路径和恢复抽样验证结果。
3. 如需安装 Go 依赖，先执行：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

4. 如需安装前端依赖，先执行：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

## 2. 最小实现路径

1. 建立后端 `sre` 服务域，优先落地高可用策略、健康总览、升级前检查与升级计划主干。
2. 建立 `maintenance window`、`capacity baseline`、`scale evidence` 和 `runbook` 相关模型与查询接口。
3. 接入 SRE 审计动作与权限语义。
4. 在前端新增平台 SRE 工作台，包括总览、高可用、升级治理、容量趋势、运行手册与审计入口。
5. 补齐契约测试、集成测试、页面测试与 smoke 验证脚本。

## 3. 联调建议

1. 先验证 `/sre/health/overview` 是否能返回统一健康总览。
2. 再验证 `/sre/ha-policies` 和 `/sre/maintenance-windows` 的配置与读取链路。
3. 再验证 `/sre/upgrades/prechecks`、`/sre/upgrades` 与 `/sre/upgrades/{upgradeId}/rollback-validations` 的升级闭环。
4. 最后验证 `/sre/capacity/baselines`、`/sre/scale-evidence`、`/sre/runbooks` 和 `/audit/sre/events`。

## 4. 验收清单

- 能创建并查看高可用策略。
- 能查看平台组件健康、依赖状态、任务积压和容量风险。
- 能执行升级前检查并获得放行/阻断结论。
- 能发起升级计划并查看滚动阶段进度。
- 能登记回退验证结果并查看剩余风险。
- 能查看容量基线、趋势、压测证据和可信度说明。
- 能将异常场景关联到运行手册和告警基线。
- 能在审计视图中查询高可用、升级、回退和维护窗口相关记录。

## 5. 建议验证命令

```bash
cd backend && go test -run TestNonExistent -count=0 ./...
cd backend && go test ./tests/contract -run TestPlatformSREContract -count=1 -p 1
cd backend && go test ./tests/integration -run TestPlatformSREIntegration -count=1 -p 1
cd frontend && npm run lint
cd frontend && npm run build
```
