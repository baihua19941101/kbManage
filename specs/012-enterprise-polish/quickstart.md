# Quickstart: 企业级治理报表与产品化交付收尾

## 1. 实施前准备

1. 确认当前分支为 `012-enterprise-polish`，禁止在 `main` 或 `master` 直接开发。
2. 先完成数据库备份，并在 `artifacts/012-enterprise-polish/backup-manifest.txt` 记录命令、时间戳、产物路径和恢复抽样验证结果。
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

1. 建立后端 `enterprise` 服务域，优先落地权限变更链路、关键操作追踪、治理覆盖率和统一待办主干。
2. 建立治理报表包、导出记录、交付材料目录、交付就绪包和交付检查清单相关模型与接口。
3. 接入企业治理报表的审计动作与权限语义。
4. 在前端新增企业交付收尾工作台，包括深度审计、治理报表、交付材料、交付检查清单和企业审计入口。
5. 补齐契约测试、集成测试、页面测试与 smoke 验证脚本。

## 3. 联调建议

1. 先验证 `/enterprise/audit/permission-trails` 和 `/enterprise/audit/key-operations` 是否能返回统一深度审计视图。
2. 再验证 `/enterprise/governance/coverage` 与 `/enterprise/governance/action-items` 的覆盖率和待办汇总链路。
3. 再验证 `/enterprise/reports` 与 `/enterprise/reports/{reportId}/exports` 的报表生成与导出闭环。
4. 最后验证 `/enterprise/delivery/artifacts`、`/enterprise/delivery/bundles`、`/enterprise/delivery/bundles/{bundleId}/checklists` 和 `/audit/enterprise/events`。

## 4. 验收清单

- 能查看完整权限变更链路。
- 能查看关键操作轨迹和高风险访问分类。
- 能查看治理覆盖率、状态分布与可信度说明。
- 能生成管理汇报、审计复核和客户交付三类治理报表。
- 能记录并查询报表导出与交付材料导出行为。
- 能查看交付材料目录及其适用版本、环境和责任角色。
- 能查看交付检查清单、缺失项和交付就绪结论。
- 能在企业审计视图中查询报表生成、导出和交付确认相关记录。

## 5. 建议验证命令

```bash
cd backend && go test -run TestNonExistent -count=0 ./...
cd backend && go test ./tests/contract -run TestEnterprisePolishContract -count=1 -p 1
cd backend && go test ./tests/integration -run TestEnterprisePolishIntegration -count=1 -p 1
cd frontend && npm run lint
cd frontend && npm run build
```
