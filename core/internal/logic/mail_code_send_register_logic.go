package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/help"
	"cloud-disk/core/models"
	"context"
	"errors"
	"time"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MailCodeSendRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendRegisterLogic {
	return &MailCodeSendRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendRegisterLogic) MailCodeSendRegister(req *types.MailCodeSendRequest) (resp *types.MailCodeSendReply, err error) {

	//邮箱未被注册

	count, err := l.svcCtx.Engine.Where("email=?", req.Email).Count(new(models.UserBasic))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		err = errors.New("该邮箱已被注册")
		return
	}
	code := help.RandCode()
	//存储验证码 key-email
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))

	//发送给用户验证码
	err = help.MailSendCode(req.Email, code)
	if err != nil {
		return nil, err
	}
	return
}
