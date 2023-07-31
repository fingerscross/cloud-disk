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
go run core.go -f etc/core-api.yaml
go get github.com/zeromicro/go-zero/rest/handler@v1.5.4
go get github.com/zeromicro/go-zero/rest/token@v1.5.4
go get github.com/zeromicro/go-zero/core/utils@v1.5.4

安装goctl插件


