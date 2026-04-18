## 007 集群生命周期中心

本次实现新增多集群集群生命周期中心，覆盖：
- 已有集群导入与新集群注册
- 模板化创建、创建前校验、升级计划与执行
- 节点池容量调整、停用与退役流程
- 驱动版本、模板资产、能力矩阵与生命周期审计页面

后端已补齐：
- 007 领域模型、迁移、仓储、路由和核心服务入口
- 007 权限语义与范围授权接入
- `clusterlifecycle.*` 审计写入和查询路由

前端已补齐：
- 生命周期中心主页面与详情页
- 注册、创建、升级、节点池、停用/退役页面
- 驱动、模板、能力矩阵与生命周期审计页面
- 菜单、路由与动作级权限门控
- 007 页面用户可见占位文案与共享命名收敛

验证结果：
- `cd backend && go test -run TestNonExistent -count=0 ./...` 通过
- `cd backend && go test ./tests/contract -count=1 -p 1` 通过
- `cd backend && go test ./tests/integration -count=1 -p 1` 通过
- `cd frontend && npm run lint` 通过
- `cd frontend && npm run build` 通过

已知遗留：
- 007 定向 `vitest` 在单 worker 模式下退出缓慢，需后续专项处理
