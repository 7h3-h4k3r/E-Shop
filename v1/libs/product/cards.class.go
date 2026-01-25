package product

import(
	"github.com/gin-gonic/gin"
	"context"
	"strings"
	"time"
	"e-commerce/v1/libs"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
	"e-commerce/v1/libs/authentication"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	
)

type Cart struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"uid"`
	Items     map[string]int     `bson:"items"`
	CreatedAt time.Time          `bson:"created_at"`
}

func productExists(id string,is_return bool)(Item ,error){
	collection := libs.GetProductDb()
	var item Item
	filter := bson.M{"public_id": id}
	err := collection.FindOne(context.Background(), filter)

	if is_document := err.Err();is_document!=nil{
		if is_document == mongo.ErrNoDocuments {
			return item,fmt.Errorf("product id invalid")
		}

		return item,fmt.Errorf("internal service error")
	}


	if is_return{
		
		err.Decode(&item)
		return item ,nil
	}
	return item,nil
}

func cartExists(ctx context.Context, userID primitive.ObjectID) error {
	collection := libs.GetCartDb()

	err := collection.FindOne(ctx, bson.M{"uid": userID}).Err()

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, Cart{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Items:     map[string]int{},
			CreatedAt: time.Now(),
		})
	}

	return err
}


func AddCartItem(c *gin.Context) {
	var cart struct{
		Product_id string `json:"protect_id":binding:"required,protect_id`
		Quantity int `json:"quantity":binding:"required,quantity`
	}

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	if cart.Quantity < 1{
		cart.Quantity = 1
	}

	ctx := context.Background()
	userID := authentication.GetId()
	collection := libs.GetCartDb()

	if err := cartExists(ctx, userID); err != nil {
		c.JSON(500, gin.H{"error": "cart initialization failed"})
		return
	}
	
	if _,err := productExists(cart.Product_id,false);err!=nil{
		c.JSON(409,gin.H{"error":err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"items." + cart.Product_id: cart.Quantity,
		},
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"uid": userID},
		update,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "database error"})
		return
	}

	c.JSON(200, gin.H{"message": "cart updated"})
}


func GetCartsItem(c *gin.Context){

	uid := authentication.GetId()
	getcollection := libs.GetCartDb()

	var usercart Cart
	if err:= getcollection.FindOne(context.Background(),bson.M{"uid":uid}).Decode(&usercart);err!=nil{
		if err == mongo.ErrNoDocuments{
			c.JSON(404,gin.H{
				"error" : "user cart empty",
			})
			return 
		}
		c.JSON(500,gin.H{
			"error" : "Internal service error ",
		})
		return
	}

	response := []Item{}

	for product_id , quantity := range usercart.Items{
		var item Item

		item , err := productExists(product_id,true)
		if err!=nil{
			c.JSON(404,gin.H{
				"error" : err.Error(),
			})
			return
		}
		item.Quantity = quantity 
		response = append(response,item)
	}

	if len(response)==0{
		c.JSON(200,gin.H{
			"state" : "empty card Item",
		})
		return
	}

	c.JSON(200,gin.H{
		"cart" : response,
	})
	return
}

func DelCartItem(c *gin.Context){

	pid := c.Param("pid")
	product_id := strings.TrimSpace(pid)
	if product_id == ""{
		c.JSON(400,gin.H{
			"error" : "product ID messup"
		})
		return
	}
	if _ , err := productExists(pid,false);err!=nil{
		c.JSON(400,gin.H{
			"error" : err.Error(),
		})
		return
	}
	uid := authentication.GetId()
	getcollection := libs.GetCartDb()

	res , err:= getcollection.UpdateOne(
		context.Background(),
		bson.M{"uid":uid}, 
		bson.M{
			"$unset":bson.M{
				"items." + product_id : "",
				},
			},
	)

	if err!=nil{
		c.JSON(409,gin.H{
			"error":"internal service error",
		})
		return
	}


	c.JSON(200,gin.H{
		"message" : "succesfully remove the item from the Cart",
	})

}