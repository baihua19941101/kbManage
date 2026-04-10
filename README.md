# kbManage

多集群 Kubernetes 可视化管理平台（开发中）。

## 目录结构

```text
backend/   # Go + Gin 后端
frontend/  # React + Vite 前端
specs/     # 规格、计划与任务
artifacts/ # 交付与治理证据
```

## 环境要求

- Go 1.25+
- Node.js 20+
- MySQL 8+
- Redis 8+

## 依赖镜像（中国网络）

后端（Go）：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

前端（npm）：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

## 后端配置（单配置文件）

后端所有可配置项统一在配置文件管理：

- 默认配置文件：`backend/config/config.dev.yaml`
- 模板文件：`backend/config/config.example.yaml`
- 可通过环境变量指定文件：`CONFIG_FILE=/path/to/config.yaml`

配置项包括：

- `server.http_addr`
- `mysql.host / port / user / password / database / parse_time`
- `redis.addr / password / db`
- `security.jwt_secret / access_token_ttl / refresh_token_ttl`

说明：

- 配置文件已包含中文注释。
- 支持少量环境变量临时覆盖（如 `MYSQL_HOST`、`REDIS_ADDR`、`HTTP_ADDR`），用于 CI 或调试。

## 前端配置（env 文件）

前端统一使用 `env` 文件配置：

- 开发默认：`frontend/.env.development`
- 模板：`frontend/.env.example`

关键项：

- `VITE_API_BASE_URL`：后端 API 前缀
- `VITE_HOST`：前端 dev server host
- `VITE_PORT`：前端 dev server 端口（可自定义）

## 启动方式

### 1) 启动后端

```bash
cd backend
# 使用默认配置文件 config/config.dev.yaml
go run ./cmd/server
```

自定义配置文件：

```bash
cd backend
CONFIG_FILE=./config/config.dev.yaml go run ./cmd/server
```

### 2) 启动前端

```bash
cd frontend
npm install
npm run dev
```

如果需要改端口，修改 `frontend/.env.development` 中的 `VITE_PORT` 即可，例如：

```env
VITE_PORT=5180
```

## 常用命令

```bash
# 仓库根目录
make backend-test
make frontend-test
make lint
make test
```

## 当前状态说明

- 已完成基础骨架、配置体系与 US1（多集群接入与资源总览）基础实现。
- 登录与刷新会话接口仍在后续故事中完善，当前可先用于页面与接口骨架联调。
