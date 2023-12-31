package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/help"
	"cloud-disk/core/models"
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	user := new(models.UserBasic)
	get, err := l.svcCtx.Engine.Where("name=? AND password =?", req.Name, help.Md5(req.Password)).Get(user)
	if err != nil {
		return nil, err
	}
	if !get {
		return nil, errors.New("用户名或者密码错误")
	}
	//生成token
	token, err := help.GenerateToken(user.Id, user.Identity, user.Name, define.TokenExpire)
	if err != nil {
		return nil, err
	}

	//用户刷新token的token
	refreshToken, err := help.GenerateToken(user.Id, user.Identity, user.Name, define.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}

	resp = new(types.LoginResponse)
	resp.Token = token
	resp.RefreshToken = refreshToken
	return
}
