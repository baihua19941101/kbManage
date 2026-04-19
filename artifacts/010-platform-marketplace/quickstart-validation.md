# 010-platform-marketplace Quickstart 验证

## 前置条件

- 使用 `backend/config/config.dev.yaml`
- 使用 `frontend/.env.development`
- MySQL、Redis 已启动
- 已完成数据库备份，见 `backup-manifest.txt`

## 验证步骤

1. 启动后端：`cd backend && go run ./cmd/server`
2. 启动前端：`cd frontend && npm run dev`
3. 登录具备 `marketplace:read` 与相应管理权限的账号
4. 打开 `/platform-marketplace/catalog-sources`，创建目录来源
5. 执行来源同步，确认模板列表出现
6. 打开 `/platform-marketplace/distribution`，选择模板并发布到工作空间/项目/集群 ID
7. 打开 `/platform-marketplace/installations`，确认安装记录与升级入口
8. 打开 `/platform-marketplace/extensions`，注册扩展并执行启停
9. 打开 `/platform-marketplace/compatibility`，查看兼容性结论
10. 打开 `/audit-events/platform-marketplace`，确认市场动作审计记录

## 验证结论

- 页面路由可访问。
- 后端 API 已提供目录、模板、分发、安装、扩展与审计查询能力。
- 主流程符合 010 规格定义的首期范围。
