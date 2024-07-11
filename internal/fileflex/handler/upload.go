package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"eim/internal/minio"
	"eim/internal/model"
	"eim/internal/model/consts"
	storagerpc "eim/internal/storage/rpc"
)

type UploadHandler struct {
	minioEndpoint string
	storageRpc    *storagerpc.Client
}

func NewUploadHandler(storageRpc *storagerpc.Client, minioEndpoint string) *UploadHandler {
	return &UploadHandler{storageRpc: storageRpc, minioEndpoint: minioEndpoint}
}

func (its *UploadHandler) Upload(c *gin.Context) {
	biz := c.MustGet("user").(*model.Biz)
	tenant, err := its.storageRpc.GetTenant(biz.TenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("get tenant -> %w", err).Error()})
		return
	}

	bucketName := tenant.Attributes[consts.FileflexBucket].String()
	userName := tenant.Attributes[consts.FileflexUser].String()
	password := tenant.Attributes[consts.FileflexPasswd].String()

	minioManager, err := minio.NewManager(&minio.Config{
		Endpoint:        its.minioEndpoint,
		AccessKeyId:     userName,
		SecretAccessKey: password,
		UseSSL:          false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("new minio manager -> %w", err).Error()})
		return
	}

	file, _ := c.FormFile("file")
	reader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Errorf("open file -> %w", err).Error()})
		return
	}

	key := fmt.Sprintf("%s/%s/%s", biz.BizId, time.Now().Format("2006-01-02"), file.Filename)
	err = minioManager.UploadObject(bucketName, key, reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("upload object -> %w", err).Error()})
		return
	}

	downloadUrl, err := minioManager.ShareObject(bucketName, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("share object -> %w", err).Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "upload success", "url": downloadUrl})
}
