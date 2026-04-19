# 010-platform-marketplace PR 摘要

## 变更范围

- 新增平台应用目录与扩展市场后端领域模型、迁移、仓储、服务与 API。
- 新增前端平台市场菜单、路由、页面、抽屉表单、hooks 与审计页面。
- 接入 010 权限语义与市场审计查询链路。
- 补充 010 配置样例、README 联调说明、验证记录与 smoke 脚本。

## 核心能力

- 目录来源创建与同步
- 模板中心与模板详情
- 模板范围分发与安装记录
- 扩展注册、启停与兼容性查看
- 市场审计事件查询

## 验证

- `cd backend && go test -run TestNonExistent -count=0 ./...`
- `cd backend && go test ./tests/contract -run TestPlatformMarketplaceContract -count=1 -p 1`
- `cd backend && go test ./tests/integration -run TestPlatformMarketplaceIntegration -count=1 -p 1`
- `cd frontend && npm run lint`
- `cd frontend && npm run build`

## 备注

- 未经用户明确同意，不得合并到 `main`。
