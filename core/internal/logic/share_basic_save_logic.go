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

type ShareBasicSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicSaveLogic {
	return &ShareBasicSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicSaveLogic) ShareBasicSave(req *types.ShareBasicSaveRequest, useridentity string) (resp *types.ShareBasicSaveReply, err error) {
	//获取资源详情
	rp := new(models.RepositoryPool)
	get, err := l.svcCtx.Engine.Where("identity = ?", req.RepositoryIdentity).Get(rp)
	if err != nil {
		return nil, err
	}
	if !get {
		return nil, errors.New("资源不存在")
	}
	//资源保存 user_repository

	ur := &models.UserRepository{
		Identity:           help.UUID(),
		UserIdentity:       useridentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                rp.Ext,
		Name:               rp.Name,
	}

	_, err = l.svcCtx.Engine.Insert(ur)
	if err != nil {
		return
	}

	resp = new(types.ShareBasicSaveReply)
	resp.Identity = ur.Identity
	return
}
