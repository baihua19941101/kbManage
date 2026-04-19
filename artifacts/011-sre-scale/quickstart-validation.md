# 011-sre-scale Quickstart 验证

## 前置条件

- 使用 `backend/config/config.dev.yaml`
- 使用 `frontend/.env.development`
- MySQL、Redis 已启动
- 已完成数据库备份，见 `backup-manifest.txt`

## 最小联调步骤

1. 启动后端：`cd backend && go run ./cmd/server`
2. 启动前端：`cd frontend && npm run dev`
3. 使用具备 `sre:read` 的账号登录
4. 打开 `/sre-scale` 查看平台健康总览
5. 打开 `/sre-scale/ha` 创建高可用策略
6. 打开 `/sre-scale/upgrades` 创建升级计划
7. 打开 `/sre-scale/rollback` 登记回退验证
8. 打开 `/sre-scale/capacity` 查看容量基线与规模化证据
9. 打开 `/sre-scale/runbooks` 查看运行手册
10. 打开 `/audit-events/sre` 查看 SRE 审计记录

## 验证结论

- 011 主干页面和 API 路由已具备联调基础。
- 当前可演示高可用、健康、升级、容量、运行手册与审计基本闭环。
