package main

import (
	"e-commerce/v1/libs"
	"e-commerce/v1/libs/authentication"
	"e-commerce/v1/libs/adminauth"
	"e-commerce/v1/libs/Middleware"
	"e-commerce/v1/libs/product"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main(){

	libs.ConnectMongo()
	
	route := gin.Default()
	route.GET("/",func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,gin.H{
			"HomePage" : "hi every one",
		})
	})
	route.POST("/login",authentication.Login)
	route.POST("/signup",authentication.Signup)
	route.POST("/admin",adminauth.AdminAuthenticate)
	fmt.Println("service start port with :8080")
	protected := route.Group("api/")
	protected.Use(Middleware.JwtMiddleware())
	protected.GET("/refresh",authentication.Refresh)

	// protected.GET("/product/:puid",product.GetProductData)
	protected.POST("/setcart",product.AddCartItem)
	protected.POST("/getcart",product.GetCartsItem)
	protected.POST("/delcart/:pid",product.DelCartItem)
	// protected.POST("/delcard/:pid",product.DelCardItem)
	protected.POST("/products",product.Products)

	admin_auth := route.Group("auth/")
	admin_auth.Use(Middleware.JwtMiddleware())
	admin_auth.POST("/additem",product.InsertItem)
	admin_auth.DELETE("/del/:id",product.DelItem)
	admin_auth.POST("/upqdata",product.UpdatePQ)
	admin_auth.POST("/Dashboard",product.Products)
	route.Run(":8080")
	
}