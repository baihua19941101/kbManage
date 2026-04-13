# PR Readiness Checklist（003）

## 分支与流程

- [x] 当前工作分支为 `003-workload-operations-control-plane`
- [x] 明确遵循：未获用户批准不得合并主干
- [ ] 已推送远程分支并创建/更新 PR（待执行）

## 治理证据

- [x] 分支/远程/合并门槛证据：`branch-check.txt`
- [x] 依赖镜像与远程仓库证据：`mirror-and-remote-check.txt`
- [x] 数据库备份与恢复抽样证据：`backup-manifest.txt`

## 功能完成度

- [x] US1 已完成
- [x] US2 已完成
- [x] US3 已完成
- [x] Final Phase（T049-T052）全部完成并复核

## 测试与质量

- [x] 后端 Contract 通过
- [x] 后端 Integration 通过
- [x] 前端 Vitest 关键用例通过
- [x] 前端变更文件 ESLint 通过
- [ ] 根目录全量 `make test` / `make lint`（待按发布窗口执行）

## 已知风险

- Ant Design 6 运行时 deprecation warning 仍存在（`Space direction`、`Drawer width`、`Alert message`），不阻塞当前功能，但建议后续统一迁移。
- 批量动作的中间件预校验使用首个 target 的 clusterId，服务层会逐目标再校验；建议追加“多目标跨 cluster 预校验”测试。
