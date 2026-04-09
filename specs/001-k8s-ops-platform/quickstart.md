# Quickstart: 多集群 Kubernetes 可视化管理平台

## 1. 前提
- 当前工作分支必须是功能分支，禁止在 `main` 上直接开发。
- 开发前先确认远程仓库和中文 PR 流程仍然有效。
- 任何数据库结构调整前都必须执行备份。

## 2. 实施前数据库备份
在开始实现数据库表结构与迁移前，执行以下命令：

```bash
mkdir -p artifacts

docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/mysql-backup-before-001-k8s-ops-platform-$(date +%Y%m%d-%H%M%S).sql
```

恢复抽样验证建议：

```bash
docker exec -i mysql8 sh -lc 'mysql -hlocalhost -P3306 -uadmin -p123456 -e "CREATE DATABASE IF NOT EXISTS restore_check;"'
docker exec -i mysql8 sh -lc 'mysql -hlocalhost -P3306 -uadmin -p123456 restore_check' < artifacts/<backup-file>.sql
```

## 3. 后端开发环境
建议的目录结构：`backend/`

依赖下载使用国内源：

```bash
export GOPROXY=https://goproxy.cn,direct
```

初始化示例：

```bash
mkdir -p backend && cd backend
go mod init github.com/baihua19941101/kbManage/backend
go env -w GOPROXY=https://goproxy.cn,direct
```

首批后端模块建议优先级：
- `auth`: 用户密码登录、刷新令牌、会话失效
- `cluster`: 集群接入、凭据校验、健康同步
- `workspace` / `project`: 作用域管理与角色绑定
- `operation`: 运维动作入队、确认与执行状态
- `audit`: 审计记录查询与导出

## 4. 前端开发环境
建议的目录结构：`frontend/`

依赖下载使用国内源：

```bash
npm config set registry https://registry.npmmirror.com
```

初始化示例：

```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
npm install antd @tanstack/react-query react-router-dom
```

首批前端页面建议优先级：
- 登录页
- 集群总览页
- 工作空间/项目管理页
- 资源列表与详情页
- 运维操作中心
- 审计查询页

## 5. 设计验证顺序
1. 先落库平台用户、角色、工作空间、项目、集群和审计表。
2. 完成用户名密码登录与平台级 RBAC。
3. 完成工作空间/项目授权与资源列表索引。
4. 接入集群健康同步与资源筛选。
5. 实现受控运维操作与审计导出。

## 6. 最小验收清单
- 能创建平台用户并完成用户名密码登录。
- 能接入至少一个集群并看到基础资源列表。
- 能创建工作空间和项目，并将角色绑定到指定作用域。
- 能发起至少一种运维操作并看到执行结果。
- 能按时间和操作者筛选审计记录并导出。

## 7. PR 交付要求
- 所有提交说明必须专业、中文、可审计。
- 推送远程后必须创建中文 PR。
- 合并前必须取得用户明确同意。
