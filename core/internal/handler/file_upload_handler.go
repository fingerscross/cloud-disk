package handler

import (
	"cloud-disk/core/help"
	"cloud-disk/core/models"
	"crypto/md5"
	"fmt"
	"net/http"
	"path"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			httpx.Error(w, err)
			return
		}

		//判断文件是否存在
		b := make([]byte, fileHeader.Size)
		_, err = file.Read(b)
		//hash 针对文件的唯一值，即使是文件的副本也不会重复上传
		//identity是上传操作的唯一值
		hash := fmt.Sprintf("%x", md5.Sum(b))
		rp := new(models.RepositoryPool)
		get, err := svcCtx.Engine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		if get {
			httpx.OkJson(w, &types.FileUploadReply{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}

		//不存在则往COS中存储文件
		cosPath, err := help.CosUpload(r)
		if err != nil {
			return
		}
		//往logic中传递request
		req.Name = fileHeader.Filename
		req.Ext = path.Ext(fileHeader.Filename)
		req.Size = fileHeader.Size
		req.Hash = hash
		req.Path = cosPath

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
