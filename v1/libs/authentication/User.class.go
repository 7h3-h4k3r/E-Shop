// authentication/auth.class.go

package authentication

import (
	"e-commerce/v1/libs"
	"net/http"
	"time"
	"strings"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)


type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:created_at" bson:"created_at"`
}


var user User

func __hashPass__(password string)([]byte,error){
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func __verifyPass__(hashed []byte,password string)error{
	return bcrypt.CompareHashAndPassword(hashed, []byte(password))
}

func Login(c *gin.Context) {
	
	var login struct{

		Username string `json:"username":binding:"required,username`
		Password string `json:"password":binding:"required,password`
	}
	
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(400, gin.H{
			"error": "credentials_miss_matching",
		})
		return
	}

	getcollection := libs.GetColl()
	filter := bson.M{"name":login.Username}
	if err := getcollection.FindOne(context.TODO(),filter).Decode(&user);err!=nil{
		if err == mongo.ErrNoDocuments{
			c.JSON(400,gin.H{
				"Login_Err" : "username not found",
			})
			return 
		}
	}
	
	if err := __verifyPass__([]byte(user.Password),login.Password);err!=nil{
		c.JSON(401,gin.H{
			"authenticate_Err" : "invalidate username (or) password",
		})
		return
	}

	access_token , err := GenareateToken(login.Username)
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

	var signup struct{
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password": binding:"required,min=3,max=16"`
		Email string `json:"email":binding:"required,email"`
	}

	if err:=c.ShouldBindJSON(&signup);err!=nil{
		c.JSON(400, gin.H{
			"error": "credentials_missing",
		})
		return 
	}
	hashed_password , _ := __hashPass__(signup.Password)
	
	user := User{
		ID : primitive.NewObjectID(),
		Name : strings.TrimSpace(signup.Username),
		Email : strings.ToLower(signup.Email),
		Password : string(hashed_password),
		CreatedAt : time.Now(),
	}
	
	getcollection := libs.GetColl()
	if _ , err := getcollection.InsertOne(context.Background(),user);err!=nil{

		if mongo.IsDuplicateKeyError(err){
			c.JSON(409,gin.H{
				"error" : "credentials already exist's ",
			})
			return
		}
		c.JSON(500,gin.H{
			"error":"database error",
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
		return
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