# BaseCode API 文档

本目录包含 BaseCode 模块的 OpenAPI 接口文档，可直接导入到 Apifox 中使用。

## 文件说明

- `baseCode-openapi.yaml` - OpenAPI v3 格式的 YAML 文档
- `baseCode-openapi.json` - OpenAPI v3 格式的 JSON 文档

## 如何导入到 Apifox

1. 打开 Apifox 应用
2. 进入项目设置或创建新项目
3. 选择「导入/导出」功能
4. 选择「OpenAPI / Swagger」格式导入
5. 选择 `baseCode-openapi.json` 或 `baseCode-openapi.yaml` 文件
6. 完成导入

## API 模块说明

文档包含以下三个主要模块：

### 1. BaseCode - 基础代码管理
- `GET /api/base_code/get` - 获取基础代码
- `GET /api/base_code/list` - 获取基础代码列表
- `POST /api/base_code/add` - 添加基础代码
- `POST /api/base_code/update` - 更新基础代码
- `POST /api/base_code/delete` - 删除基础代码

### 2. BaseCodeItem - 基础代码项管理
- `GET /api/base_code_item/get` - 获取基础代码项
- `GET /api/base_code_item/list` - 获取基础代码项列表
- `POST /api/base_code_item/add` - 添加基础代码项
- `POST /api/base_code_item/update` - 更新基础代码项
- `POST /api/base_code_item/delete` - 删除基础代码项

### 3. BaseCodeSort - 基础代码分类管理
- `GET /api/base_code_sort/get` - 获取基础代码分类
- `GET /api/base_code_sort/list` - 获取基础代码分类列表
- `POST /api/base_code_sort/add` - 添加基础代码分类
- `POST /api/base_code_sort/update` - 更新基础代码分类
- `POST /api/base_code_sort/delete` - 删除基础代码分类

## 服务地址

默认配置的服务地址：
- 本地开发环境：http://localhost:9090

可根据实际环境在 Apifox 中修改服务地址。
