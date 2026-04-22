# 012-enterprise-polish Quickstart 验证

## 前置条件

- 使用 `backend/config/config.dev.yaml`
- 使用 `frontend/.env.development`
- MySQL、Redis 已启动
- 已完成数据库备份，见 `backup-manifest.txt`

## 最小联调步骤

1. 启动后端：`cd backend && go run ./cmd/server`
2. 启动前端：`cd frontend && npm run dev`
3. 使用具备 `enterprise:read` 的账号登录
4. 打开 `/enterprise-polish` 查看权限变更审计
5. 打开 `/enterprise-polish/risks` 查看关键操作与风险趋势
6. 打开 `/enterprise-polish/reports` 生成治理报表
7. 打开 `/enterprise-polish/exports` 查看治理待办和导出背景数据
8. 打开 `/enterprise-polish/artifacts` 查看交付材料目录
9. 打开 `/enterprise-polish/delivery` 查看交付清单和交付就绪状态
10. 打开 `/audit-events/enterprise` 查看企业治理审计记录

## 验证结论

- 012 主干页面和 API 路由已具备联调基础。
- 当前可演示深度审计、治理报表、导出留痕、交付材料目录与交付清单基本闭环。
