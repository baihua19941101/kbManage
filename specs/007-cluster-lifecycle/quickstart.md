# Quickstart: 多集群 Kubernetes 集群生命周期中心

## 当前执行状态（2026-04-17）

- 007 已完成 `specify` 与 `plan` 阶段，当前已生成 `spec.md`、`plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 和质量清单。
- 当前工作分支为 `007-cluster-lifecycle`。
- 后续进入 `/speckit.tasks` 或实现前，仍需按仓库治理要求落实数据库备份、国内依赖源配置、中文 PR 和用户明确同意后再合并。

## 1. 前提

- 仓库位于功能分支 `007-cluster-lifecycle`
- 后端配置文件统一使用 `backend/config/config.dev.yaml`
- 前端配置文件统一使用 `frontend/.env.development`
- Go 依赖通过 `GOPROXY=https://goproxy.cn,direct`
- npm 依赖通过 `https://registry.npmmirror.com`

## 2. 实施前数据库备份

在开始实现 007 之前，先执行：

```bash
mkdir -p artifacts/007-cluster-lifecycle
docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/007-cluster-lifecycle/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

如果当前环境 `admin/123456` 不可用，必须：

1. 改用容器实际可用凭据执行备份
2. 在 `artifacts/007-cluster-lifecycle/backup-manifest.txt` 中记录差异原因、实际命令和产物路径
3. 使用临时 MySQL 容器做恢复抽样验证

恢复抽样示例：

```bash
docker run -d --rm --name mysql8-restore-check-007 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
cat artifacts/007-cluster-lifecycle/mysql-backup-<timestamp>.sql \
  | docker exec -i mysql8-restore-check-007 sh -lc 'MYSQL_PWD=123456 mysql -h127.0.0.1 -uroot'
docker exec mysql8-restore-check-007 sh -lc 'MYSQL_PWD=123456 mysql -h127.0.0.1 -N -B -uroot -e "SHOW DATABASES;"'
```

## 3. 国内依赖源配置

### Go

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### npm

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

### 其他联调依赖

- 若需要驱动联调镜像、Kubernetes 创建依赖或本地模拟环境，优先使用阿里云、DaoCloud 或已批准的国内代理。
- 禁止在计划和任务中省略镜像源说明。

## 4. 最小联调路径

完成实现后，最小联调至少覆盖：

1. 导入一个已有集群，确认列表中出现接入状态、版本和健康摘要
2. 生成一个注册指引并完成注册状态确认
3. 以一个驱动模板发起一次创建前校验，确认阻断项与通过项都可见
4. 为一个已纳管集群创建升级计划并查看状态变更
5. 对一个节点池执行扩缩动作并确认互斥动作限制
6. 对一个停用/退役流程生成审计记录并可检索
7. 查看至少两个驱动或集群类型的能力矩阵差异

## 5. 完整验收清单

- 导入、注册、创建、升级、节点池调整、停用、退役均有明确状态流
- 创建前校验能区分阻断项与风险提示项
- 驱动版本、模板和能力矩阵之间存在清晰关联
- 权限回收后关键动作立即锁定
- 生命周期审计可按时间、操作者、动作类型、集群状态和结果筛选
- README、配置说明和端口配置方式在实现阶段同步更新

## 6. 规划产物

- 规格：`specs/007-cluster-lifecycle/spec.md`
- 计划：`specs/007-cluster-lifecycle/plan.md`
- 研究：`specs/007-cluster-lifecycle/research.md`
- 数据模型：`specs/007-cluster-lifecycle/data-model.md`
- 契约：`specs/007-cluster-lifecycle/contracts/openapi.yaml`
