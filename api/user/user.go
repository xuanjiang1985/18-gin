package user

import (
	"18-gin/help"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Name(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"code": help.Statecode.Success,
		"msg":  "成功",
		"content": gin.H{
			"name": name,
		},
	})
}
