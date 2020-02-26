## scan.service

浏览器后端服务
#####build
```
sh build.sh
```
   
#####deploy
修改配置文件：./conf/scan.conf
```
#####重要参数#####
Port     = "8359"                               #服务启动的端口，作为nginx的上游，为前端提供数据接口
Redis    = "xxx:8379"                           #redis缓存服务, 预留给加速用，当前可不配
Gate     = "gateway.inner.poc.com:8545"         #poc链网关节点rpc，浏览器通过此节点获取链数据
database = "xxx:xxx@2019@tcp(xxx:8306)/scan"    #数据库，格式化链上数据，以提供快速查询
```

#####API
参见：src/main.go 和 src/api

#####database
参见：src/model/create_table.go