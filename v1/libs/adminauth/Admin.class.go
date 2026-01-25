package adminauth

// authentication/auth.class.go

import (

	"encoding/json"
	"os"
	"github.com/gin-gonic/gin"
	"e-commerce/v1/libs/authentication"
	"net/http"

)
type Admin struct{
	ADMIN string `json:"ADMIN"`
}

type Admin_user struct{
	AdminLoginKey string `json:"key"`
}


func getenv() string{
	var key Admin

	file , err := os.ReadFile("../aenv.json")

	if err!=nil{
		return " "
	}

	if err := json.Unmarshal(file,&key);err!=nil{
		return " "
	}
	return key.ADMIN
}

func AdminAuthenticate(c *gin.Context){
	var key Admin_user
		
	if err := c.ShouldBindJSON(&key);err!=nil{
		c.JSON(401,gin.H{
			"error" : "Admin key Missing not validated",
		})
		return
	}

	if err := getenv();err!=key.AdminLoginKey{
		c.JSON(401,gin.H{
			"error" : "Admin-authentication Credential's Invalidate",
		})
		return
	}
	
	access_token , err := authentication.GenareateToken("admin")
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"authenticate_Err" : "server side Authoraization access_token Problem",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"access-token" :access_token,
		"message" : "Welcome Back Admin , How are you , i process a very well",
		"bool" : true,
	})
}