#### api模块脚本命令

user模块为例子
```sh
cd script
# 生成user模块的api
sh apigen user
# 重置user模块的api
sh apireset user
```

#### rpc模块脚本命令

user模块为例子
```sh
cd script
# 生成user模块的rpc
sh rpcgen user
# 重置user模块的rpc
sh rpcreset user
```

#### model脚本命令
生成user表为例子
```sh
cd app/{对应的模块}/
# amigo为数据库名称当前项目固定 user为表名按需修改
modelgen amigo user
```
