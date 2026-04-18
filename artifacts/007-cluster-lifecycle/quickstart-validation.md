quickstart_validation_for: 007-cluster-lifecycle
validated_at: 2026-04-17

backend:
  config_file: backend/config/config.dev.yaml
  result: passed
  notes:
    - 007 路由已注册到主路由树
    - Go 编译与 007 contract/integration 包测试通过

frontend:
  env_file: frontend/.env.development
  result: passed
  notes:
    - 007 页面、菜单与权限门禁已加入主应用路由
    - lint 与生产构建通过

manual_smoke_focus:
  - /cluster-lifecycle
  - /cluster-lifecycle/register
  - /cluster-lifecycle/provision
  - /cluster-lifecycle/upgrades
  - /cluster-lifecycle/node-pools
  - /cluster-lifecycle/retirement
  - /cluster-lifecycle/drivers
  - /cluster-lifecycle/templates
  - /cluster-lifecycle/capabilities
  - /audit-events/cluster-lifecycle

follow_up:
  - 收敛 vitest 退出缓慢问题后再做 007 全量页面回归
