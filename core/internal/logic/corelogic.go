package logic

import (
	"bytes"
	"cloud-disk/core/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoreLogic {
	return &CoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoreLogic) Core(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	data := make([]*models.UserBasic, 0)
	err = models.Engine.Find(&data)
	if err != nil {
		log.Print("Get UserBasic Error:", err)
	}
	fmt.Println(data) //直接输出struct会输出地址 需要转化
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("Marshal Error", err)
	}

	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", " ") //转成byte buffer 再用string形式展示
	if err != nil {
		log.Println("Json indent Error", err)
	}
	fmt.Println(dst.String())
	resp = new(types.Response)
	resp.Message = dst.String()

	return
}
