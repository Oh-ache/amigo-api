# Amigo API 文档

本目录包含后端微服务（`baseCode` / `device` / `sdk` / `user`）的接口文档。

## 目录结构

```
document/
├── README.md                       # 本文件
├── api/                            # Markdown 格式（给人看）
│   ├── baseCode.md
│   ├── device.md
│   ├── sdk.md
│   └── user.md
├── swagger/                        # Swagger 2.0 格式（给工具用）
│   ├── all.swagger.json            # ⭐ 整合版本：所有模块合一份
│   ├── baseCode.swagger.json
│   ├── device.swagger.json
│   ├── sdk.swagger.json
│   └── user.swagger.json
└── sdk/
    └── minimax/                    # MiniMax AI 子模块独立使用文档
```

## 源文件（最权威）

接口定义在 `common/api/<module>.api`，改这里后必须跑：

```sh
cd script
sh apireset <module>      # 重新生成 handler / logic / types / routes
```

更新本文档：

```sh
# Markdown（按 .api 文件批量重生成）
cd common/api
goctl api doc --dir .

# Swagger（按模块分别重生成）
goctl-swagger swagger --filename <output.json> --host <host> --basepath <basepath> --schemes http,https <module>.api
```

## 推荐使用：导入 Apifox / Postman

### ⭐ 整合版本（一键导入所有模块）

直接导入 `document/swagger/all.swagger.json`，可一次性看到 `baseCode` / `device` / `sdk` / `user` 全部接口。

- 默认 host：`192.168.31.62:8888`（gateway）
- 默认 schemes：`http`、`https`
- 跨模块同名类型（如 `CommonResp` / `EmptyResp`）已合并；其它类型自动加模块前缀（如 `sdk.WeatherReq`、`user.UserLoginReq`）

### 分模块版本

适合只需要某个模块的场景：

| 文件 | host（建议） | 备注 |
|------|--------------|------|
| `swagger/baseCode.swagger.json` | `192.168.31.62:9090` | 基础代码管理 |
| `swagger/sdk.swagger.json` | `192.168.31.62:9091` | SDK 服务（含天气/IP/OSS/验证码/AI 等） |
| `swagger/user.swagger.json` | `192.168.31.62:9092` | 用户管理 |
| `swagger/device.swagger.json` | `192.168.31.62:9093` | 设备管理 |

## 注意事项

- `goctl-swagger` 官方插件（20220621 版本）已过期，与当前 `goctl 1.8.2` 不兼容，会报 `unexpected end of JSON input`。
  - 临时方案：手工编辑或在 `apireset` 流程里挂自定义脚本
  - 长期方案：升级 `goctl-swagger` 或迁移到 OpenAPI v3 工具链
- `api/<module>.md` 是 `goctl api doc` 自动生成的，仅展示路由 + 类型结构，不展示嵌套类型定义。
  整合 markdown 可以考虑后续用工具二次处理。
- `document/README.md`（旧版）原本只讲 baseCode 模块，现已重写为本总览。

## 服务地址

| 模块 | 直连端口 | gateway 路由前缀 |
|------|---------|------------------|
| gateway | `8888` | `/api/<module>/` |
| baseCode | `9090` | `/api/base_code/`、`/api/base_code_item/`、`/api/base_code_sort/` |
| sdk | `9091` | `/api/sdk/` |
| user | `9092` | `/api/user/` |
| device | `9093` | `/api/device/`、`/api/firmware/`、`/api/firmware_task/`、`/api/work_order/`、`/api/debug/` |
| ai | `9094` | `/api/ai/` |
| mqueue | `7001` | `/api/mqueue/` |