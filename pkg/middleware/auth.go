package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"eim/pkg/log"
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

	c.Set("user", user)

	c.Next()
}
