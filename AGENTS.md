## 语言设置
- 你必须始终使用**简体中文**回复。
- 代码中的注释、解释说明、思考过程均需使用中文。
- 仅代码本身的语法关键字（如 function, var, class 等）保留英文。

# Agent Guidance

## Stack
- Go 1.25 with go-zero framework
- Monorepo: multiple modules under `app/`
- PM2 for process management

## Modules
| Module | Path | Notes |
|--------|------|-------|
| gateway | `app/gateway/` | REST entrypoint |
| baseCode, device, user, sdk | `app/<name>/api/` + `app/<name>/rpc/` | API + gRPC pairs |
| queue | `app/job/queue/` | Async job worker |
| mqueue | `app/job/mqueue/` | Message queue handler |

## Code Generation
```sh
./script/apigen <module>   # Generate API from common/api/<module>.api
./script/rpcgen <module>   # Generate RPC from common/proto/<module>.proto
```
Generated code lives under `app/<module>/api/` and `app/<module>/rpc/`. Common proto/pb lives in `common/pb/`.

## Running Services
```sh
pm2 start ecosystem.config.js    # All services
pm2 start <name>                 # Single service (e.g., baseCodeRpc)
```

## Dependencies
- Redis (go-redis), MySQL (xorm), gRPC, Asynq for queues
- Config files: `etc/<service>.yaml` per module

## Script Commands

### API Module
```sh
cd script
# Generate user module API
sh apigen user
# Reset user module API
sh apireset user
```

### RPC Module
```sh
cd script
# Generate user module RPC
sh rpcgen user
# Reset user module RPC
sh rpcreset user
```

### Model Generation
```sh
cd app/{module}/
# amigo=db name, user=table name
modelgen amigo user
```
