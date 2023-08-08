package logic

import (
	"cloud-disk/core/help"
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadChunkCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkCompleteLogic {
	return &FileUploadChunkCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkCompleteLogic) FileUploadChunkComplete(req *types.FileUploadChunkCompleteRequest) (resp *types.FileUploadChunkCompleteReply, err error) {

	co := make([]cos.Object, 0)
	for _, object := range req.CosObjects {
		co = append(co, cos.Object{
			ETag:       object.Etag,
			PartNumber: object.PartNumber,
		})
	}
	err = help.CosPartUploadComplete(req.Key, req.UploadId, co)
	if err != nil {
		return
	}
	return
}
