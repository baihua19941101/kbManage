# PR Readiness Checklist - 002

## 基本信息

- Feature Branch: `002-observability-center`
- Target Branch: `main`
- 状态: Ready for review（待用户批准合并）

## 宪章与流程检查

- [x] 专用分支开发（未在 main 直接开发）
- [x] 数据库备份完成并有恢复抽样记录
- [x] 使用国内源/代理策略已记录
- [x] 交付说明与文档使用中文
- [x] 明确记录“合并需用户批准”

## 质量检查

- [x] 后端契约测试覆盖 US1/US2/US3
- [x] 后端集成测试覆盖 US1/US2/US3
- [x] 前端页面测试覆盖 US1/US2/US3
- [x] 前端 lint 通过
- [x] 关键权限与审计路径覆盖

## 风险清单（已知）

1. Ant Design 在 Vitest 环境出现 deprecation warning（不影响通过），建议后续统一升级组件 API。
2. 当前 Alertmanager/Loki/Prometheus 仍以兼容接口为主，生产环境需进一步做真实后端联调与性能压测。
3. SQLite 测试环境下禁用 observability sync worker 以避免锁冲突，生产 MySQL 不受该限制。

## 发布前建议

1. 在目标环境执行一次真实 Prometheus/Loki/Alertmanager 联调 smoke。
2. 补充 30 天范围告警检索的压测数据与阈值报告。
3. 评估告警治理动作的审计字段是否满足合规侧检索要求。

## 合并门禁

- [ ] 用户确认同意合并到 `main`
