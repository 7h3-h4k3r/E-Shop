package authentication

import (
	"e-commerce/v1/libs/envread"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct{
	Username string
	jwt.RegisteredClaims

}


func getKey() []byte{
	envread.Getenv()
	return []byte(envread.Env_file_read.JwtKey)
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