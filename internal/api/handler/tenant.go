package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"eim/internal/minio"
	"eim/internal/model"
	"eim/internal/model/consts"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/stringsutil"
)

type TenantHandler struct {
	minioManager *minio.Manager
	storageRpc   *storagerpc.Client
}

func NewTenantHandler(storageRpc *storagerpc.Client, minioManager *minio.Manager) *TenantHandler {
	return &TenantHandler{storageRpc: storageRpc, minioManager: minioManager}
}

func (its *TenantHandler) Register(c *gin.Context) {
	var tenant *model.Tenant
	err := c.ShouldBindJSON(&tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("bind json -> %w", err).Error()})
		return
	}

	err = its.checkTenant(tenant.TenantId)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant already exists"})
		return
	}

	err = its.storageRpc.InsertTenant(tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("save tenant -> %w", err).Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (its *TenantHandler) Update(c *gin.Context) {
	tenantId := c.Param("tenantId")
	err := its.checkTenant(tenantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("check tenant -> %w", err).Error()})
		return
	}

	var tenant *model.Tenant
	err = c.ShouldBindJSON(&tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("bind json -> %w", err).Error()})
		return
	}

	tenant.TenantId = tenantId
	err = its.storageRpc.UpdateTenant(tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("save tenant -> %w", err).Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (its *TenantHandler) EnableFileFlex(c *gin.Context) {
	tenantId := c.Param("tenantId")
	err := its.checkTenant(tenantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("check tenant -> %w", err).Error()})
		return
	}

	tenant, err := its.storageRpc.GetTenant(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("get tenant -> %w", err).Error()})
		return
	}

	bucketName := tenantId
	userName := stringsutil.RandomString(10)
	password := stringsutil.RandomString(10)
	err = its.minioManager.CreateUser(userName, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("create user -> %w", err).Error()})
		return
	}

	err = its.minioManager.CreateBucket(bucketName)
	if err != nil {
		_ = its.minioManager.RemoveUser(userName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("create bucket -> %w", err).Error()})
		return
	}

	err = its.minioManager.AttachBucketPolicy(bucketName, userName)
	if err != nil {
		_ = its.minioManager.RemoveUser(userName)
		_ = its.minioManager.RemoveBucket(bucketName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("attach bucket policy -> %w", err).Error()})
		return
	}

	tenant.Attributes[consts.FileflexEnabled] = fmt.Sprintf("%v", true)
	tenant.Attributes[consts.FileflexBucket] = bucketName
	tenant.Attributes[consts.FileflexUser] = userName
	tenant.Attributes[consts.FileflexPasswd] = password
	err = its.storageRpc.UpdateTenant(tenant)
	if err != nil {
		_ = its.minioManager.DetachBucketPolicy(bucketName, userName)
		_ = its.minioManager.RemoveUser(userName)
		_ = its.minioManager.RemoveBucket(bucketName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("save tenant -> %w", err).Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (its *TenantHandler) checkTenant(tenantId string) error {
	if tenantId == "" {
		return fmt.Errorf("tenantId is empty")
	}
	_, err := its.storageRpc.GetTenant(tenantId)
	if err != nil {
		return fmt.Errorf("get tenant -> %w", err)
	}
	return nil
}
