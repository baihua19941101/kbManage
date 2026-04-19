# Quickstart: 平台应用目录与扩展市场

## 当前执行状态（2026-04-19）

- 010 已完成 `specify` 与 `plan` 阶段。
- 当前工作分支为 `010-platform-marketplace`。
- 实现前仍需按仓库治理要求完成数据库备份、国内依赖源配置、中文 PR 与用户明确同意后合并。

## 1. 前提

- 仓库位于功能分支 `010-platform-marketplace`
- 后端配置统一使用 `backend/config/config.dev.yaml`
- 前端配置统一使用 `frontend/.env.development`
- Go 依赖通过 `GOPROXY=https://goproxy.cn,direct`
- npm 依赖通过 `https://registry.npmmirror.com`

## 2. 实现前数据库备份

```bash
mkdir -p artifacts/010-platform-marketplace
docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/010-platform-marketplace/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

如果当前环境 `admin/123456` 不可用，必须：

1. 改用容器内实际可用凭据完成备份。
2. 在 `artifacts/010-platform-marketplace/backup-manifest.txt` 记录差异原因、实际命令和产物路径。
3. 使用临时 MySQL 容器做恢复抽样验证。

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

### 目录与扩展相关外部来源

- 如需接入目录来源、模板仓库或扩展制品代理，优先选择已批准的国内镜像或代理。
- 计划与任务中必须记录来源策略，不得省略。

## 4. 最小联调路径

完成实现后，最小联调至少覆盖：

1. 新增一个目录来源并完成一次同步。
2. 导入至少两个模板，其中一个模板包含多个版本和依赖关系。
3. 将模板发布到一个工作空间范围和一个项目范围。
4. 查看目标范围内的模板可见性、安装记录和升级入口。
5. 注册一个扩展包并查看兼容性结论、权限声明和可见范围。
6. 对扩展执行一次启用和一次停用。
7. 检索至少一条模板分发审计记录和一条扩展生命周期审计记录。

## 5. 验收清单

- 目录来源、模板分类、版本、依赖、参数表单和部署约束可在统一目录中心查看
- 模板按工作空间、项目或集群范围受控分发
- 安装记录、升级入口、版本变更说明和下线状态清晰可见
- 扩展包、插件和集成模块支持注册、启停、兼容性和权限声明治理
- 权限模型、审计模型和发布模型已贯通
- README、配置说明和验证材料在实现阶段同步更新

## 6. 规划产物

- 规格：`specs/010-platform-marketplace/spec.md`
- 计划：`specs/010-platform-marketplace/plan.md`
- 研究：`specs/010-platform-marketplace/research.md`
- 数据模型：`specs/010-platform-marketplace/data-model.md`
- 契约：`specs/010-platform-marketplace/contracts/openapi.yaml`
