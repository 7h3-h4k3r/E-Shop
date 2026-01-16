// Database.class.go
package libs

import (
	"context"
	"e-commerce/v1/libs/envread"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var Client *mongo.Client

func ConnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	envread.Getenv()
	uri := envread.Env_file_read.MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	return nil
}



func GetColl()*mongo.Collection{
	envread.Getenv()
	return  Client.Database(envread.Env_file_read.DataBase).Collection(envread.Env_file_read.Collection)

}	

func GetRefreshColl()*mongo.Collection{
	envread.Getenv()
	fmt.Println(envread.Env_file_read.DataBase,envread.Env_file_read.Collection_r)
	return  Client.Database(envread.Env_file_read.DataBase).Collection(envread.Env_file_read.Collection_r)

}
