# 009 Quickstart 验证

根据 `specs/009-identity-tenancy/quickstart.md`，本轮至少完成了以下验证：

1. 已创建外部身份源并保留本地兜底身份，`/identity/sources` 与 `/identity/sessions` contract/integration 通过。
2. 已创建组织单元并完成租户边界映射，`/identity/organizations` 与 `/identity/organizations/{unitId}/mappings` contract/integration 通过。
3. 已创建角色定义、授权分配与委派链路，`/identity/roles`、`/identity/assignments`、`/identity/delegations` contract/integration 通过。
4. 已生成会话治理与访问风险视图，`/identity/sessions` 与 `/identity/access-risks` contract/integration 通过。
5. 已暴露身份治理审计查询入口 `/audit/identity/events`，前端审计页已接入路由与菜单。

补充说明：

- 登录方式切换由 `/identity/login-mode` 补充实现，用于支撑统一登录方式切换体验。
- 会话回收由 `/identity/sessions/{sessionId}/revoke` 补充实现，用于支撑前端会话治理操作。
