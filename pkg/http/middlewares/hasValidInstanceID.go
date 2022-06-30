package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tekenradar/content-service/pkg/http/helpers"
)

func HasValidInstanceID(validIDs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		instanceID := c.Param("instanceID")
		if instanceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A valid InstanceID is missing"})
			c.Abort()
			return
		}

		err := helpers.CheckInstanceID(validIDs, instanceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
