package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/model"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

type DownloadHandler struct {
	minioEndpoint string
	fileRpc       *storagerpc.FileClient
}

func NewDownloadHandler(fileRpc *storagerpc.FileClient, minioEndpoint string) *DownloadHandler {
	return &DownloadHandler{fileRpc: fileRpc, minioEndpoint: minioEndpoint}
}

func (its *DownloadHandler) Download(c *gin.Context) {
	biz := c.MustGet("user").(*model.Biz)
	tenant := c.MustGet("tenant").(*model.Tenant)

	bucketName := c.Param("bucket_name")
	filePath := c.Param("file_path")

	its.fileRpc.GetFile()

}
