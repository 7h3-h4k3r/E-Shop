

package authentication

import (
	"e-commerce/v1/libs/envread"
	"e-commerce/v1/libs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"crypto/sha256"
    "encoding/hex"
	"crypto/rand"
    "encoding/base64"
	"fmt"

	"context"
	"strconv"

)


type Claims struct{
	Username string
	jwt.RegisteredClaims

}

type RefreshID struct{
	Refresh_token_hashed string `json:"refresh"`
}


type RefreshToken struct{
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Username string `bson:"username"`
	UserId primitive.ObjectID `bson:"user_id"`
	CreatedAt time.Time `bson:"create_at"`
	ExpiresAt string `bson:"expire_at"`
	Refresh_token_hashed string `bson:"hashed_token"`
	Revoke bool `bson:"revoke"`
}
func getKey() []byte{
	envread.Getenv()
	return []byte(envread.Env_file_read.JwtKey)
}


func generateSecureRandom() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(token string) string {
    sum := sha256.Sum256([]byte(token))
    return hex.EncodeToString(sum[:])
}

func getrkey() []byte{
	envread.Getenv()
	return []byte(envread.Env_file_read.JwtRKey)
}


func GenareateToken(username string)(string , error){

	claim := Claims{
		Username : username,
		RegisteredClaims : jwt.RegisteredClaims{
			ExpiresAt : jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt : jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claim)
	return token.SignedString(getKey())
}

func AuthGranted(token_string string) error {
	token , err := jwt.ParseWithClaims(
		token_string,
		&Claims{},
		func (token *jwt.Token)(any,error){
			return getKey(),nil
		},
	)

	if err!=nil || !token.Valid{
		return err
	}
	
	return nil
}



func GetRefreshToken(user string, user_id primitive.ObjectID)(string , error){
	refresh_token , err := generateSecureRandom()
	if err!=nil{
		return "Error",fmt.Errorf("generateSecurerandom() func is Err")
	}
	refresh_token_hashed := hashToken(refresh_token)

	r_token := RefreshToken{
		Username : user,
		UserId : user_id,
		CreatedAt : time.Now(),
		ExpiresAt : strconv.FormatInt(time.Now().Add(7 * 24 * time.Hour).Unix(),10),
		Refresh_token_hashed : refresh_token_hashed,
		Revoke : false,
	}

	getcollection := libs.GetRefreshColl()
	
	if _ ,err := getcollection.InsertOne(context.Background(),r_token);err!=nil{
		return "Error", fmt.Errorf("insert refresh_token_base_Error :%v",getcollection)
	}

	return refresh_token,nil
}



func SetRefreshToken(tokenId primitive.ObjectID)error{

	filter := bson.M{
		"_id" : tokenId,
		"revoke":false,
	}

	update := bson.M{
		"$set": bson.M{
			"revoke": true,
		},
	}

	getcollection := libs.GetRefreshColl()

	revoking , err := getcollection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return err
	}

	if revoking.MatchedCount == 0{
		return fmt.Errorf("no token registered in the database ")
	}
	return nil
}

func GetRDetails(r_token *RefreshToken , token_id string) error{

	refresh_token_hashed := hashToken(token_id)
 
	filter := bson.M{
		// "expire_at": bson.M{
		// 	"$gt": time.Now(),
		// },
		"hashed_token": refresh_token_hashed,
		"revoke":    false,
		
	}
	fmt.Println(refresh_token_hashed)

	getcollection := libs.GetRefreshColl()
	if  err := getcollection.FindOne(context.Background(),filter).Decode(&r_token);err!=nil{
		if err == mongo.ErrNoDocuments{
			return fmt.Errorf("refresh token are not found or expired")
		}
		return err
	}

	return nil

}