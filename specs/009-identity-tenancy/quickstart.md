# Quickstart: 身份与多租户治理中心

## 当前执行状态（2026-04-19）

- 009 已完成 `specify` 与 `plan` 阶段，当前已生成 `spec.md`、`plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 和质量清单。
- 当前工作分支为 `009-identity-tenancy`。
- 后续进入 `/speckit.tasks` 或实现前，仍需按仓库治理要求落实数据库备份、国内依赖源配置、中文 PR 和用户明确同意后再合并。

## 1. 前提

- 仓库位于功能分支 `009-identity-tenancy`
- 后端配置文件统一使用 `backend/config/config.dev.yaml`
- 前端配置文件统一使用 `frontend/.env.development`
- Go 依赖通过 `GOPROXY=https://goproxy.cn,direct`
- npm 依赖通过 `https://registry.npmmirror.com`

## 2. 实施前数据库备份

在开始实现 009 之前，先执行：

```bash
mkdir -p artifacts/009-identity-tenancy
docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/009-identity-tenancy/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

如果当前环境 `admin/123456` 不可用，必须：

1. 改用容器实际可用凭据执行备份
2. 在 `artifacts/009-identity-tenancy/backup-manifest.txt` 中记录差异原因、实际命令和产物路径
3. 使用临时 MySQL 容器做恢复抽样验证

恢复抽样示例：

```bash
docker run -d --rm --name mysql8-restore-check-009 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
cat artifacts/009-identity-tenancy/mysql-backup-<timestamp>.sql \
  | docker exec -i mysql8-restore-check-009 sh -lc 'MYSQL_PWD=123456 mysql -h127.0.0.1 -uroot'
docker exec mysql8-restore-check-009 sh -lc 'MYSQL_PWD=123456 mysql -h127.0.0.1 -N -B -uroot -e "SHOW DATABASES;"'
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

- 若需要身份目录 SDK、SSO 联调镜像、LDAP 连接组件或目录同步工具，优先使用阿里云、DaoCloud 或已批准的国内代理。
- 禁止在计划和任务中省略镜像源说明。

## 4. 最小联调路径

完成实现后，最小联调至少覆盖：

1. 新增一个外部身份源并保留本地管理员登录能力
2. 创建一个组织、两个团队和至少一个用户组
3. 将组织对象映射到工作空间和项目范围
4. 创建平台级、组织级和项目级角色定义
5. 为用户或用户组分配直接授权、委派授权和临时授权
6. 查看会话列表并验证权限变更后会话影响可见
7. 检索至少一条身份与租户治理域审计记录和一条访问风险摘要

## 5. 完整验收清单

- 身份源、本地账号和统一登录入口可清晰区分并共存
- 组织、团队、用户组、工作空间和项目映射关系可见且边界清晰
- 平台级到资源级角色都有明确适用范围和继承边界
- 委派、临时授权和回收链路可审计、可回看
- 会话治理能够识别权限变化后的受影响访问
- 风险视图能够展示越权、残留权限、继承扩散和异常会话提示
- README、配置说明和端口配置方式在实现阶段同步更新

## 6. 规划产物

- 规格：`specs/009-identity-tenancy/spec.md`
- 计划：`specs/009-identity-tenancy/plan.md`
- 研究：`specs/009-identity-tenancy/research.md`
- 数据模型：`specs/009-identity-tenancy/data-model.md`
- 契约：`specs/009-identity-tenancy/contracts/openapi.yaml`
