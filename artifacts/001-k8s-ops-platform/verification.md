# Verification Baseline - 001-k8s-ops-platform

## 验证时间
- 2026-04-10（Asia/Shanghai）

## 后端
- 命令：`cd backend && go test ./...`
- 结果：通过
- 说明：`tests/contract` 与 `tests/integration` 全部通过；其余目录为 no test files。

## 前端
- 命令：`cd frontend && npm run lint`
- 结果：通过

- 命令：`cd frontend && npm run test -- --run`
- 结果：通过（`4` 个测试文件，`7` 个测试用例）
- 说明：存在 Ant Design deprecation warning 与 `Could not parse CSS stylesheet` 警告，不阻塞测试通过。

- 命令：`cd frontend && npm run build`
- 结果：通过
- 说明：已完成路由级懒加载，构建产物从单一大包拆分为多 chunk（如 `LoginPage`、`ClusterOverviewPage`、`ResourceListPage`、`AuditEventPage` 等）；当前仍有一个 `client` chunk（约 606k）高于默认告警阈值。

## 结论
- 当前代码基线满足本轮交付验证要求，可进入 PR 汇总与风险评审阶段。

## 性能证据附件
- 测试环境报告：`artifacts/001-k8s-ops-platform/test-environment-load-report.md`
- 可复现实验脚本：`artifacts/001-k8s-ops-platform/repro-perf-smoke.sh`
