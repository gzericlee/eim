package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/minio"
	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/internal/model/consts"
	seqrpc "github.com/gzericlee/eim/internal/seq/rpc/client"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

type UploadHandler struct {
	minioEndpoint           string
	seqRpc                  *seqrpc.SeqClient
	tenantRpc               *storagerpc.TenantClient
	fileRpc                 *storagerpc.FileClient
	externalServiceEndpoint string
}

func NewUploadHandler(tenantRpc *storagerpc.TenantClient, seqTpc *seqrpc.SeqClient, fileRpc *storagerpc.FileClient, minioEndpoint, externalServiceEndpoint string) *UploadHandler {
	return &UploadHandler{tenantRpc: tenantRpc, seqRpc: seqTpc, fileRpc: fileRpc, minioEndpoint: minioEndpoint, externalServiceEndpoint: externalServiceEndpoint}
}

func (its *UploadHandler) Upload(c *gin.Context) {
	biz := c.MustGet("user").(*model.Biz)
	tenant := c.MustGet("tenant").(*model.Tenant)
	scope := c.GetString("scope")
	if scope == "" {
		scope = "*"
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

	today := time.Now().Format("2006-01-02")

	upFile, _ := c.FormFile("file")
	reader, err := upFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Errorf("open file -> %w", err).Error()})
		return
	}

	filePath := fmt.Sprintf("%s/%s/%s", biz.BizId, today, upFile.Filename)
	err = minioManager.UploadObject(bucketName, filePath, reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("upload object -> %w", err).Error()})
		return
	}

	fileId, err := its.seqRpc.SnowflakeId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("get snowflake id -> %w", err).Error()})
		return
	}

	file := &model.File{
		FileId:   fileId,
		FileName: upFile.Filename,
		FileType: filepath.Ext(upFile.Filename),
		FilePath: filePath,
		FileSize: upFile.Size,
		BizId:    biz.BizId,
		TenantId: tenant.TenantId,
		Attributes: map[string]string{
			"scope": scope,
		},
	}

	err = its.fileRpc.InsertFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("insert file -> %w", err).Error()})
		return
	}

	thumbPath := fmt.Sprintf("%d/24x24/%s", file.FileId, file.FileName)

	downloadUrl := fmt.Sprintf("%s/%d/%s", its.externalServiceEndpoint, fileId, upFile.Filename)
	thumbUrl := fmt.Sprintf("%s/%s/%s", its.externalServiceEndpoint, bucketName, thumbPath)

	c.JSON(http.StatusOK, gin.H{"message": "upload success", "download": downloadUrl, "thumb": thumbUrl, "explain": "24x24 can be changed to other sizes"})
}
