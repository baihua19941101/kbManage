# PR 摘要（005-security-and-policy）

## 目标
实现多集群安全与策略治理中心首期能力：策略中心、准入模式切换、灰度验证、例外治理、违规闭环与策略审计检索。

## 本次范围
- Phase 0-2：治理门槛、配置与基础设施
- US1：策略定义/分配/层级展示
- US2：模式切换、例外申请审批、例外到期回收
- US3：违规查询、整改状态更新、策略审计查询

## 关键接口
- `/api/v1/security-policies`
- `/api/v1/security-policies/:policyId/assignments`
- `/api/v1/security-policies/:policyId/mode-switch`
- `/api/v1/security-policies/hits`
- `/api/v1/security-policies/hits/:hitId/exceptions`
- `/api/v1/security-policies/hits/:hitId/remediation`
- `/api/v1/security-policies/exceptions`
- `/api/v1/security-policies/exceptions/:exceptionId/review`
- `/api/v1/audit/security-policies/events`

## 前端页面
- `/security-policies`
- `/security-policies/rollout`
- `/security-policies/violations`
- `/audit-events/security-policy`

## 验证
- 后端：`go test -p 1 ./...` PASS
- 前端：005 相关测试 `--maxWorkers=1` PASS
- 前端：005 相关 ESLint PASS

## 风险与后续
- 例外 `startsAt` 当前未单独持久化字段，后续可模型增强。
- Ant Design deprecation warning 存在，不影响功能。
- 性能压测与更细颗粒审计导出可在后续迭代增强。
