package v1

import (
	"github.com/harranali/authority"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ahmadcahyana/arctic-fox/constants"
	"github.com/ahmadcahyana/arctic-fox/datatransfers"
	"github.com/ahmadcahyana/arctic-fox/handlers"
	"github.com/ahmadcahyana/arctic-fox/models"
)

func GETUser(c *gin.Context) {
	var err error
	var userInfo datatransfers.UserInfo
	if err = c.ShouldBindUri(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}

	var user models.User
	if user, err = handlers.Handler.RetrieveUser(userInfo.Username); err != nil {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "cannot find user"})
		return
	}
	auth := authority.Resolve()
	ok, errPermission := auth.CheckPermission(user.ID, "create")
	if errPermission != nil || !ok {
		c.JSON(http.StatusNotFound, datatransfers.Response{Error: "You don't have permission"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.Response{Data: datatransfers.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
	}})
}

func PUTUser(c *gin.Context) {
	var err error
	var user datatransfers.UserUpdate
	if err = c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, datatransfers.Response{Error: err.Error()})
		return
	}
	if err = handlers.Handler.UpdateUser(uint(c.GetInt(constants.IsAuthenticatedKey)), user); err != nil {
		c.JSON(http.StatusNotModified, datatransfers.Response{Error: "failed updating user"})
		return
	}
	c.JSON(http.StatusOK, datatransfers.Response{Data: user})
}
