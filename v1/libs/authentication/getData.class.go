// authentication/auth.class.go
package authentication

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func GetId() primitive.ObjectID {
	return user.ID
}

func GetUsername() string{
	return user.Name
}
