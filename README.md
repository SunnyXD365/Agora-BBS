# Agora-BBS

Agora-BBS 是一套把阅读感知、渐进权限、冷静期、语境反馈、匿名盲审和 LLM 辅助治理嵌入论坛流程的全栈项目。

## 快速启动

```powershell
Copy-Item backend/.env.example backend/.env
docker-compose -f docker-compose.dev.yml up -d
Get-Content backend/seed.sql -Raw | docker exec -i agora-postgres psql -U agora_user -d agora_db
```

访问：

- 论坛：<http://localhost:8080>
- API 健康检查：<http://localhost:8080/api/ping>
- Temporal UI：<http://localhost:8233>

普通用户和管理员共用 `/login`。管理员密码正确后还需完成邮箱验证码，随后自动进入独立的 `/admin` 管理端；开发环境未配置 SMTP 时登录页会提供本地演示码，生产环境必须配置 SMTP。

开发 Seed 管理员 `admin` 的密码为 `admin123`，其余演示账号密码均为 `password`：`demo_admin`、`demo_l0`、`demo_l1`、`demo_l2`、`demo_l3_a`、`demo_l3_b`、`demo_l3_c`。Seed 仅用于本地，生产环境使用 `backend` 镜像内的 `/app/admin` 命令提升已有用户。

## 架构

```text
Browser → Nginx → Next.js / Go API → PostgreSQL / Redis
                                  → Temporal → Go Worker → OpenAI-compatible LLM
```

开发、测试和生产说明见：

- [开发环境指南](docs/20_开发环境指南.md)
- [API 文档](docs/30_API文档.md)
- [Git 开发规范](docs/40_Git%20开发规范.md)
- [项目计划与架构说明](docs/50_项目开发计划与架构说明.md)
- [测试、部署与运维](docs/60_测试部署与运维.md)

## 质量检查

```powershell
docker-compose -f docker-compose.dev.yml exec -T agora-backend-api sh -c "go test ./... && go vet ./..."
docker-compose -f docker-compose.dev.yml exec -T agora-frontend sh -c "npm run lint && npm run build"
docker-compose -f docker-compose.dev.yml --profile test run --rm agora-e2e
docker-compose -f docker-compose.dev.yml --profile test run --rm agora-k6
```

仓库只提交 `.env.example`。真实数据库密码、JWT Secret 和 LLM API Key 不得进入 Git。
