// authentication/auth.class.go
package authentication

import (
	
	"e-commerce/v1/libs"
	"net/http"
	"time"
	"fmt"
	"strings"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)


type SignupReq struct{
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password": binding:"required,min=3,max=16"`
	Email string `json:"email":binding:"required,email"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
}

type LoginReq struct{
	Username string `json:"username"`
	Password string `json:"password"`
}


func __hashPass__(password string)([]byte,error){
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func __verifyPass__(hashed []byte,password string)error{
	return bcrypt.CompareHashAndPassword(hashed, []byte(password))
}


func Login(c *gin.Context) {
	
	var req LoginReq
	var user User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "credentials_miss_matching",
		})
		return
	}
	getcollection := libs.GetColl()
	filter := bson.M{"name":req.Username}
	if err := getcollection.FindOne(context.TODO(),filter).Decode(&user);err!=nil{
		if err == mongo.ErrNoDocuments{
			c.JSON(400,gin.H{
				"Login_Err" : "username not found",
			})
			return 
		}
	}
	
	if err := __verifyPass__([]byte(user.Password),req.Password);err!=nil{
		c.JSON(401,gin.H{
			"authenticate_Err" : "invalidate username (or) password",
		})
		return
	}

	access_token , err := GenareateToken(req.Username)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"authenticate_Err" : "server side Authoraization access_token Problem",
		})
		return
	}
	
	refresh_token , err := GetRefreshToken(user.Name,user.ID)
	
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"authenticate_Err" : "server side Authoraization refresh_token Problem",
			"err" : err.Error(),
		})
		return
	}

	
	c.JSON(http.StatusOK,gin.H{
		"access-token" :access_token,
		"refresh-token":refresh_token,
		"message" : "user succesfully logged",
		"bool" : true,
	})
	
}

func Signup(c *gin.Context){
	var in_val SignupReq

	if err:=c.ShouldBindJSON(&in_val);err!=nil{
		c.JSON(400, gin.H{
			"error": "credentials_missing",
		})
		return 
	}
	hashed_password , _ := __hashPass__(in_val.Password)
	user := User{
		ID : primitive.NewObjectID(),
		Name : strings.TrimSpace(in_val.Username),
		Email : strings.ToLower(in_val.Email),
		Password : string(hashed_password),
		CreatedAt : time.Now(),
	}
	getcollection := libs.GetColl()
	if _ , err := getcollection.InsertOne(context.Background(),user);err!=nil{
		c.JSON(500,gin.H{
			"error":"database InsertProblem Signup()",
		})
		return
	}
	
	c.JSON(201,gin.H{
		"message":"signup success",
		"userId" : user.ID,
		"status" : true,
	})
}


func Refresh(c *gin.Context) {

	var r_token RefreshID
	
	var r_token_struct RefreshToken

	if err := c.ShouldBindJSON(&r_token);err!=nil{
		c.JSON(401,gin.H{
			"error" : "missing credentials refresh token missing",
		})
		return
	}

	if err := GetRDetails(&r_token_struct,r_token.Refresh_token_hashed);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error" : err.Error(),
		})
		fmt.Println("yes this side the GetRDetails")
		return
	}else{
		fmt.Println("crossing the GetRDetails")
	}

	if err := SetRefreshToken(r_token_struct.Id);err!=nil{
		c.JSON(401,gin.H{
			"error" : err.Error(),
		})
		return
	}


	access_token , err := GenareateToken(r_token_struct.Username)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"authenticate_Err" : "server side Authoraization access_token Problem",
		})
		return
	}
	
	refresh_token , err := GetRefreshToken(r_token_struct.Username,r_token_struct.UserId)
	
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"authenticate_Err" : "server side Authoraization refresh_token Problem",
			"err" : err.Error(),
		})
		return
	}

	
	c.JSON(http.StatusOK,gin.H{
		"access-token" :access_token,
		"refresh-token":refresh_token,
		"message" : "user succesfully logged",
		"bool" : true,
	})

	
}	