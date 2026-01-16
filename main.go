package main

import (
	"e-commerce/v1/libs"
	"e-commerce/v1/libs/authentication"
	"e-commerce/v1/libs/Middleware"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main(){

	// token , _ := authentication.GenareateToken("dharain")
	// fmt.Println(token)
	// err := authentication.AuthGranted("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImRoYXJhaW4iLCJleHAiOjE3Njc5NTEyMjAsImlhdCI6MTc2Nz.13R-ngV1GMi7ogiVReJ3qQko_oHDXYYwy8N5JLsIDGw")

	// if err := libs.ConnectMongo();err!=nil{
	// 	fmt.Println("DataBase connection is Failed")
	// }else{
	// 	fmt.Println("its Connected succesfully");
	// }

	// collection := libs.GetColl()
	// user := libs.User{
	// 	Name: "dharani",

	// }
	 
	// _ , err := collection.InsertOne(context.TODO(),user)
	// if err!= nil{
	// 	fmt.Println("somethink is Wrong")
	// }
	libs.ConnectMongo()
	
	route := gin.Default()
	route.GET("/",func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,gin.H{
			"HomePage" : "hi every one",
		})
	})
	route.GET("/login",authentication.Login)
	route.GET("/signup",authentication.Signup)
	
	fmt.Println("service start port with :8080")
	protected := route.Group("api/")
	protected.Use(Middleware.JwtMiddleware())
	protected.GET("/refresh",authentication.Refresh)
	protected.GET("/product",func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"message":"welcome to the Product Daseboard Page",
		})
	})
	route.Run(":8080")
	
}