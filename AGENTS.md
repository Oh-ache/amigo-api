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

## 新增接口流程

以 `user` 模块为例，新增一个列表查询接口。

### 1. 修改 .api 文件（HTTP 层）

在 `common/api/<module>.api` 中添加请求/响应类型和路由：

```api
type (
    // 请求参数 — GET 用 form 标签，POST 用 json 标签
    ListUserReq {
        Page     int64  `form:"page,default=1,optional"`
        PageSize int64  `form:"page_size,default=10,optional"`
        Username string `form:"username,default=,optional"`
    }
    ListUserResp {
        List  []GetUserResp `json:"list"`
        Total int64         `json:"total"`
    }
)

@server (
    group:  user
    prefix: /api/user
    jwt:    Auth       // 需要鉴权则保留，公开接口去掉此行
)
service user {
    @handler UserList
    get /list (ListUserReq) returns (ListUserResp)
}
```

> 所有 handler 已通过模板自动包裹 `CommonResp{code, msg, data}`，api 文件中的 `returns` 类型最终位于 `data` 字段内。空响应使用 `import "common.api"` 提供的 `EmptyResp`。

### 2. 修改 .proto 文件（RPC 层）

在 `common/proto/<module>.proto` 中添加 message 和 rpc 方法：

```proto
message ListUserReq {
    int64 page = 1;
    int64 page_size = 2;
    string username = 3;
}

message ListUserResp {
    repeated UserResp list = 1;
    int64 total = 2;
}

service User {
    rpc ListUser (ListUserReq) returns (ListUserResp);
}
```

### 3. 一键生成代码

```sh
cd script

# 生成 API（handler、logic、types、routes）
sh apigen user

# 生成 RPC（pb 文件、logic、server）
sh rpcgen user
```

| 命令 | 生成内容 | 是否覆盖已有文件 |
|------|---------|:---:|
| `apigen <m>` | handler/*.go, logic/*.go, types/types.go, handler/routes.go | 仅新增 handler/logic，覆盖 types/routes |
| `apireset <m>` | 同上 | **删除后重新生成**所有 handler/logic/types/routes |
| `rpcgen <m>` | common/pb/<m>.pb.go, common/pb/<m>_grpc.pb.go, logic/*.go, server/*.go | 仅新增 logic，覆盖 pb/server |
| `rpcreset <m>` | 同上 | **删除后重新生成** pb/server |

> 日常迭代用 `apigen`/`rpcgen`（保留已实现的逻辑代码）。大幅重构时用 `apireset`/`rpcreset`（注意备份已写好的 logic）。

### 4. 实现业务逻辑

生成后的模板文件位于：

| 层 | 路径 | 说明 |
|----|------|------|
| API handler | `app/<m>/api/internal/handler/<group>/<name>Handler.go` | 已自动包裹 CommonResp，**一般无需修改** |
| API logic | `app/<m>/api/internal/logic/<group>/<name>Logic.go` | 在此组装参数，调用 RPC client |
| RPC logic | `app/<m>/rpc/internal/logic/<name>Logic.go` | 在此操作 DB/缓存，实现核心业务 |
| RPC server | `app/<m>/rpc/internal/server/<m>Server.go` | 自动生成，**无需修改** |

数据流：**Handler（脱壳 HTTP）→ API Logic（组装 RPC 请求）→ RPC Logic（DB/缓存操作）**

### 5. 新数据表

如果接口涉及新表，先生成 model：

```sh
cd app/<module>/
modelgen <db_name> <table_name>   # 示例: modelgen amigo user
```

生成文件：
- `app/<m>/model/<table>Model_gen.go` — 自动生成的基础 CRUD，**勿手动修改**
- `app/<m>/model/<table>Model.go` — 扩展方法（List、CheckDuplicate 等）

然后在 `app/<m>/rpc/internal/svc/serviceContext.go` 中注入 model 实例。

### 6. 同级调用

如需在 API logic 中调用同一模块下已有 RPC logic（不走 gRPC 重开端口），直接在 svc 中暴露内部方法，API logic 通过 `l.svcCtx` 访问。

### 7. 新模块

如果是**全新模块**（非已有模块新增接口），还需：
- 新建 `etc/<service>.yaml` 配置文件
- 在 `app/gateway/internal/handler/routes.go` 中新增反向代理路由
- 在 `app/gateway/internal/config/config.go` 的 `Upstreams` 中新增字段
- 新建 Docker 构建文件

已有模块新增接口以上步骤自动生效，无需额外配置。
