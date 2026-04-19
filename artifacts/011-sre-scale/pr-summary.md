# 011-sre-scale PR 摘要

## 变更范围

- 新增平台 SRE 与规模化治理后端领域模型、迁移、适配层、仓储、服务与 API 路由。
- 新增前端 SRE 工作台页面、组件、hooks、服务层、菜单和审计页。
- 接入 011 权限语义、SRE 审计查询链路、配置样例与 README 联调说明。
- 补充 011 数据库备份证据、验证记录和 smoke 脚本。

## 核心能力

- 高可用策略与维护窗口
- 平台健康总览
- 升级前检查与升级计划
- 回退验证
- 容量基线、规模化证据与运行手册
- SRE 审计查询

## 验证

- `cd backend && go test -run TestNonExistent -count=0 ./...`
- `cd backend && go test ./tests/contract -run TestPlatformSREContract -count=1 -p 1`
- `cd backend && go test ./tests/integration -run TestPlatformSREIntegration -count=1 -p 1`
- `cd frontend && npm run lint`
- `cd frontend && npm run build`

## 风险

- 前端 Vitest runner 存在项目级挂起问题，011 页面测试文件已补齐但未能在本地纳入通过项。
- 当前后端和前端主干可编译可联调，剩余风险集中在前端 Vitest runner 稳定性。
