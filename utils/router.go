package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ahmadcahyana/arctic-fox/constants"
	"github.com/ahmadcahyana/arctic-fox/datatransfers"
)

func AuthOnly(c *gin.Context) {
	if !c.GetBool(constants.IsAuthenticatedKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.Response{Error: "user not authenticated"})
	}
}
