package logic

import (
	"cloud-disk/core/define"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type UserFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListRequest, useridentity string) (resp *types.UserFileListReply, err error) {
	//使用xorm查数据
	uf := make([]*types.UserFile, 0)

	//分页参数
	size := req.Size
	if size == 0 {
		size = define.PageSize
	}
	page := req.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * size

	//数据连表查询
	resp = new(types.UserFileListReply)
	err = l.svcCtx.Engine.Table("user_repository").Where("parent_id=? AND user_identity=?", req.Id, useridentity).
		Select("user_repository.id,user_repository.identity,user_repository.repository_identity ,user_repository.ext,"+
			"user_repository.name,repository_pool.path,repository_pool.size").
		Join("LEFT", "repository_pool", "user_repository.repository_identity=repository_pool.identity").
		Where("user_repository.deleted_at = ? OR user_repository.deleted_at IS NULL", time.Time{}.Format(define.DateTime)).
		Limit(size, offset).Find(&uf)
	if err != nil {
		return
	}
	//查询文件总数
	count, err := l.svcCtx.Engine.Where("parent_id=? AND user_identity=?", req.Id, useridentity).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	resp.List = uf
	resp.Count = count
	return
}
