package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/minio"
	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/internal/model/consts"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

type UploadHandler struct {
	minioEndpoint string
	tenantRpc     *storagerpc.TenantClient
}

func NewUploadHandler(tenantRpc *storagerpc.TenantClient, minioEndpoint string) *UploadHandler {
	return &UploadHandler{tenantRpc: tenantRpc, minioEndpoint: minioEndpoint}
}

func (its *UploadHandler) Upload(c *gin.Context) {
	biz := c.MustGet("user").(*model.Biz)
	tenant, err := its.tenantRpc.GetTenant(biz.TenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("get tenant -> %w", err).Error()})
		return
	}

	enabled, _ := strconv.ParseBool(tenant.Attributes[consts.FileflexEnabled])
	if !enabled {
		c.JSON(http.StatusBadRequest, gin.H{"message": "fileflex is disabled"})
		return
	}

	bucketName := tenant.Attributes[consts.FileflexBucket]
	userName := tenant.Attributes[consts.FileflexUser]
	password := tenant.Attributes[consts.FileflexPasswd]

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
