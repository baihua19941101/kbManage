# kbManage Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-04-09

## Active Technologies
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query (001-k8s-ops-platform)
- MySQL 8.4（平台元数据、权限、审计、资源索引）；Redis 8.x（会话、短时缓存、任务协同） (001-k8s-ops-platform)

## Project Structure

```text
backend/
frontend/
specs/
```

## Commands

npm test && npm run lint

## Code Style

Go 1.25；TypeScript 5.x；React 19.2: Follow standard conventions

## Recent Changes
- 001-k8s-ops-platform: Added Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query

<!-- MANUAL ADDITIONS START -->
## Configuration Conventions

- 后端所有可配置项必须集中在 `backend/config/*.yaml`，默认使用 `backend/config/config.dev.yaml`。
- 后端配置文件必须带注释，新增配置项时同步更新 `backend/config/config.example.yaml`。
- 后端仅允许通过 `CONFIG_FILE` 指定配置文件路径；如需临时调试可使用同名环境变量覆盖单项。
- 前端配置统一放在 `frontend/.env.development`（开发）与 `frontend/.env.example`（模板）。
- 前端启动端口必须通过 `VITE_PORT` 配置，不允许在脚本中硬编码固定端口。
- README 必须同步维护“配置说明 + 启动命令 + 端口配置方式”。
<!-- MANUAL ADDITIONS END -->


## 协作代理规则
- 必须使用中文与用户交流
- 每次新会话都必须读取spec了解项目功能和计划
- 每完成一个任务，必须更新spec
- 如果有新的需求或改动，必须先和用户进行讨论，经过用户允许之后必须立刻更新spec
- 完成新的需求或改动，必须立刻更新spec
- 所有语言框架安装依赖时必须使用国内源或代理下载

## 开发流程规则
- 新功能必须使用分支开发模式，从远程的main分支拉取代码，禁止在主分支进行开发，开发完成后，提交github
- 合并分支必须经过用户同意后才能进行合并
- 提交代码必须包含专业的说明描述

## 数据库备份流程
- 开发新功能前必须保证数据库备份完成，才能进行新功能开发
- 数据库备份操作可以登录到docker容器内进行