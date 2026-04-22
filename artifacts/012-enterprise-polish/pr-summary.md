# 012-enterprise-polish PR 摘要

## 变更范围

- 新增企业级治理报表与产品化交付收尾后端领域模型、迁移、适配层、仓储、服务与 API 路由。
- 新增前端企业治理工作台页面、组件、hooks、服务层、菜单和企业审计页。
- 接入 012 权限语义、企业审计查询链路、配置样例与 README 联调说明。
- 补充 012 数据库备份证据、验证记录和 smoke 脚本。

## 核心能力

- 深度权限审计与关键操作追踪
- 治理覆盖率与统一治理待办
- 管理汇报、审计复核和客户交付三类治理报表
- 报表导出留痕
- 交付材料目录、交付就绪包与交付检查清单
- 企业治理审计查询

## 验证

- `cd backend && go test -run TestNonExistent -count=0 ./...`
- `cd backend && go test ./tests/contract -run TestEnterprisePolishContract -count=1 -p 1`
- `cd backend && go test ./tests/integration -run TestEnterprisePolishIntegration -count=1 -p 1`
- `cd frontend && npm run lint`
- `cd frontend && npm run build`

## 风险

- 前端 Vitest runner 存在项目级挂起问题，012 页面测试文件已补齐但未能在本地纳入通过项。
- 012 当前实现以企业治理收尾主线为主，报表模板与交付材料内容仍属于首期标准模板，不包含客户专属深度定制。
