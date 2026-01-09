package libs

import (
	"e-commerce/v1/libs/authentication"
	"net/http"

	"github.com/gin-gonic/gin"
)



func Login(c *gin.Context) {
	username := c.Query("username");
	password := c.Query("password");
	if username == "" || password == ""{
		c.JSON(http.StatusNotFound,gin.H{
				"error" : "Credential Missing",
		})
		return
	}
	if username == "admin" && password == "pass@!23"{
		token ,err := authentication.GenareateToken(username)
		if err==nil{
			c.JSON(http.StatusOK,gin.H{
				"access-token" : token,
				"bool" : true,
			})
			return
		}
		c.JSON(http.StatusUnauthorized,gin.H{
				"authenticate" : err,
			})

	}else{
		c.JSON(http.StatusBadRequest,gin.H{
			"authenticate" : "invalidate username (or) password",
		})
	}
}