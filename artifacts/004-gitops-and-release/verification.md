# 004 验证基线（2026-04-13）

## 分支与范围
- 分支：`004-gitops-and-release`
- 范围：US1（建模+状态）/US2（发布生命周期）/US3（权限与审计）

## 后端验证
- 命令：`cd backend && go test ./...`
- 结果：通过
- 覆盖点：
  - GitOps contract：sources/target-groups/delivery-units/actions/diff/releases/access-control/audit
  - GitOps integration：modeling/sync/promotion/rollback/scope-authorization/audit

## 前端验证
- 命令：`cd frontend && npm run test -- --run src/features/gitops`
- 结果：通过（7 files / 11 tests）
- 命令：`cd frontend && npm run lint`
- 结果：通过

## 关键验收点
- GitOps 细粒度鉴权：来源管理、目标组、交付单元、动作、operation 查询均按 scope + permission 生效。
- 审计闭环：`gitops.source.verify` 与 `gitops.*.submit` 动作写入审计，支持审计查询过滤。
- 权限回收处理：前端在 403 场景下展示“权限已变更/已回收”并锁定动作按钮。

## 复现脚本验证
- 命令：`bash artifacts/004-gitops-and-release/repro-gitops-smoke.sh`
- 结果：
  - backend tests ✅
  - frontend gitops tests ✅
  - frontend lint ✅
  - runtime probe（`/healthz`）为非阻塞检查，当前记录为 skipped/failed

## 已知非阻塞项
- Ant Design 6.3 相关 deprecation warning 仍存在（不影响测试通过与功能正确性）。
