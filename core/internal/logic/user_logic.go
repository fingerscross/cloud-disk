package logic

import (
	"cloud-disk/core/help"
	"cloud-disk/core/models"
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	user := new(models.UserBasic)
	get, err := models.Engine.Where("name=? AND password =?", req.Name, help.Md5(req.Password)).Get(user)
	if err != nil {
		return nil, err
	}
	if !get {
		return nil, errors.New("用户名或者密码错误")
	}
	//生成token
	token, err := help.GenerateToken(user.Id, user.Identity, user.Name)
	if err != nil {
		return nil, err
	}
	resp = new(types.LoginResponse)
	resp.Token = token

	return
}
