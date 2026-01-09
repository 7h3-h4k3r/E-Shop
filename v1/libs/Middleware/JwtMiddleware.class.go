package Middleware

import (
	"e-commerce/v1/libs/authentication"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		auth := c.GetHeader("Authorization")

		if auth == ""{
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		auth_data := strings.Split(auth," ")

		if len(auth_data)!=2 || auth_data[0] != "Bearer"{
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} 

		if err := authentication.AuthGranted(auth_data[1]);err!=nil{
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.JSON(http.StatusAccepted,gin.H{
			"state" : "authGranted",
			"bool" : true,
		})

		c.Next()
	}
}