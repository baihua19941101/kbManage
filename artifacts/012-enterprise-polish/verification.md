# 012-enterprise-polish 验证记录

日期：2026-04-20

## 已完成验证

- 后端静态编译：`cd backend && go test -run TestNonExistent -count=0 ./...`
- 012 后端契约测试：`cd backend && go test ./tests/contract -run TestEnterprisePolishContract -count=1 -p 1`
- 012 后端集成测试：`cd backend && go test ./tests/integration -run TestEnterprisePolishIntegration -count=1 -p 1`
- 前端 Lint：`cd frontend && npm run lint`
- 前端构建：`cd frontend && npm run build`
- 前端页面测试文件已补齐：`frontend/src/features/enterprise-polish/pages/*.test.tsx`、`frontend/src/features/audit/pages/EnterpriseAuditPage.test.tsx`

## 结果摘要

- 012 新增后端企业治理领域、路由、权限、审计与 Redis key 接线后编译通过。
- 012 专属后端契约测试与集成测试已在低并发模式下通过。
- 012 前端企业治理页面、菜单、路由与企业审计页构建通过。
- 012 的数据库备份与恢复抽样验证已完成，运行时使用 `root/123456`，`admin/123456` 不可用的偏差已记录在 `backup-manifest.txt`。

## 未完成项

- 前端 Vitest runner 在当前仓库仍存在项目级挂起问题；单文件执行 `PermissionAuditPage.test.tsx` 停在 `RUN` 阶段，012 页面测试文件因此未纳入通过项。
