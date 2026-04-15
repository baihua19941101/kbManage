# PR Readiness（005-security-and-policy）

## 治理检查
- [x] 功能分支开发（005-security-and-policy）
- [x] 数据库备份与恢复抽样验证
- [x] 国内依赖源配置说明
- [x] 中文交付文档与 PR 材料
- [x] 子代理模型固定为 gpt-5.3-codex（medium）

## 测试检查
- [x] 后端全量测试通过（`go test -p 1 ./...`）
- [x] 前端 005 相关测试通过（`--maxWorkers=1`）
- [x] 前端 005 相关 lint 通过

## 交付检查
- [x] spec/plan/tasks 已同步最新执行状态
- [x] quickstart 验证记录已补齐
- [x] smoke 脚本已提供

## 合并门槛
- [ ] 已推送远程分支
- [ ] 已创建/更新中文 PR
- [ ] 已获得用户明确同意合并
