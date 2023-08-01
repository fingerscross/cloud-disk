# cloud-disk

GOPROXY=https://goproxy.cn,direct

建表
初始化项目：go mod init cloud-disk
 alt + enter


go get -u github.com/go-sql-driver/mysql

go get xorm.io/xorm

go get -u github.com/zeromicro/go-zero
go install github.com/zeromicro/go-zero/tools/goctl@latest
goctl -v
goctl api new core  //这里叫core服务

#启动服务
go run core.go -f etc/core-api.yaml
go get github.com/zeromicro/go-zero/rest/handler@v1.5.4
go get github.com/zeromicro/go-zero/rest/token@v1.5.4
go get github.com/zeromicro/go-zero/core/utils@v1.5.4

安装goctl插件


goctl api go -api core.api -dir . -style go_zero

#发送验证码
go get github.com/jordan-wright/email

#存验证码用redis
go get github.com/redis/go-redis/v9

#本地启动redis
redis-server.exe redis.windows.conf
#进入redis操作
redis-cli.exe -h 127.0.0.1 -p 6379

#存identity 用uuid生成
go get github.com/satori/go.uuid
