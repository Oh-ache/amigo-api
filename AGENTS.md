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

## ⚠️ 修改 `.api` / `.proto` 后必须重置

**规则**：每次修改 `common/api/<module>.api` 或 `common/proto/<module>.proto` 后，**必须**用 `apireset` / `rpcreset` 重新生成代码，**不可**只跑 `apigen` / `rpcgen`。

### 为什么必须重置

- `apigen` / `rpcgen` 是**增量**的，只会新增 handler / logic，**不会清理**已删除的接口、字段、类型。
- 一旦 `.api` 中**删除**了某个路由或类型、`.proto` 中**删除**了某个 message 或 rpc 方法，旧的生成文件会残留，导致编译通过但运行时出现死代码 / 类型不一致 / 路由冲突。
- `apireset` / `rpcreset` 会先清空 `types/`、`routes.go`、`common/pb/<m>.pb.go` 等生成产物，再基于最新的 `.api` / `.proto` 全量重建，保证生成代码与定义文件严格一致。

### 执行清单

修改 `.api` 或 `.proto` 后，按以下步骤操作：

1. **重置前**：用 `git status` 确认手写的业务逻辑（`logic/*.go`）已提交或已备份，避免被覆盖。
2. **重置**：进入 `script/` 目录，针对修改的每个模块执行：
   ```sh
   cd script

   # 修改了 .api
   sh apireset <module>

   # 修改了 .proto
   sh rpcreset <module>
   ```
   两个文件都改了，两个脚本都要跑。
3. **重置后**：用 `git diff --stat` 核对生成文件的改动范围，确认与预期一致。
4. **还原手写逻辑**：如果被覆盖，从备份 / git 历史中恢复 `app/<m>/api/internal/logic/`、`app/<m>/rpc/internal/logic/` 下的业务代码。
5. **编译验证**：
   ```sh
   go build ./...
   go vet ./app/<m>/...
   ```
6. **重启服务**（如涉及运行中的服务）：
   ```sh
   pm2 reload <module>Rpc <module>Api
   ```

### 多模块批量重置

如果一次改动涉及多个模块，可一次性传入：

```sh
cd script
sh apireset sdk device
sh rpcreset sdk device
```

### 跳过重置的反模式（禁止）

- ❌ 只跑 `apigen` / `rpcgen`，不跑 `apireset` / `rpcreset`
- ❌ 手动编辑 `common/pb/*.pb.go` 或 `app/<m>/rpc/internal/server/*Server.go`
- ❌ 手动编辑 `app/<m>/api/internal/handler/routes.go`、`app/<m>/api/internal/types/`
- ❌ 修改 `.api` / `.proto` 后直接 `git commit`，不重新生成

以上任何一条都会让生成代码与定义文件脱钩，后续修改容易引发编译失败或运行时异常。

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
- 在 `app/gateway/etc/gateway-api.yaml` 的 `Routes` 列表中新增一条路由规则
- 新建 Docker 构建文件

已有模块新增接口以上步骤自动生效，无需额外配置。
