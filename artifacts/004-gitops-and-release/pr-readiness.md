# PR 就绪清单（004）

## 代码与功能
- [x] US1 功能完成
- [x] US2 功能完成
- [x] US3 功能完成
- [x] GitOps 审计页完成并接入路由

## 测试与质量
- [x] backend `go test ./...`
- [x] frontend `npm run test -- --run src/features/gitops`
- [x] frontend `npm run lint`
- [x] 关键错误路径（403 权限回收）具备前端锁定态

## 文档与交付物
- [x] `specs/004-gitops-and-release/spec.md` 状态已更新
- [x] `specs/004-gitops-and-release/tasks.md` 任务勾选已更新
- [x] `artifacts/004-gitops-and-release/verification.md`
- [x] `artifacts/004-gitops-and-release/quickstart-validation.md`
- [x] `artifacts/004-gitops-and-release/repro-gitops-smoke.sh`
- [x] `artifacts/004-gitops-and-release/pr-summary.md`

## 合并前确认
- [ ] 用户明确批准合并
- [ ] 推送远程并创建中文 PR
