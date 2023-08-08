package logic

import (
	"cloud-disk/core/help"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicCreateLogic {
	return &ShareBasicCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicCreateLogic) ShareBasicCreate(req *types.ShareBasicCreateRequest, useridentity string) (resp *types.ShareBasicCreateReply, err error) {
	uuid := help.UUID()
	ur := new(models.UserRepository)
	get, err := l.svcCtx.Engine.Where("identity = ?", req.UserRepositoryIdentity).Get(ur)
	if err != nil {
		return nil, err
	}
	if !get {
		return nil, errors.New("数据不存在")
	}

	data := &models.ShareBasic{
		Identity:           uuid,
		UserIdentity:       useridentity,
		RepositoryIdentity: ur.RepositoryIdentity,
		ExpiredTime:        req.ExpiredTime,
	}
	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return nil, err
	}

	resp = &types.ShareBasicCreateReply{Identity: uuid}
	return
}
