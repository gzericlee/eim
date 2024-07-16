package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gzericlee/eim/pkg/log"
)

func (its *ginMiddleware) Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.Replace(token, "Basic ", "", 1)

	user, err := its.authRpc.CheckToken(token)
	if err != nil {
		log.Error("Error check auth token", zap.Error(err))
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	tenant, err := its.tenantRpc.GetTenant(user.TenantId)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("get tenant -> %w", err))
		return
	}

	c.Set("user", user)
	c.Set("tenant", tenant)

	c.Next()
}
