package product 

import(
	"e-commerce/v1/libs"
	"time"
	"context"
	"strconv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/google/uuid"
)


type B_data struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	PublicId string `bson:"public_id,omitempty"`
	ProductName string `json:"product" bson:"product"`
	Price int `json:"price" bson:"price"`
	Description string `json:"desc" bson:"desc"`
	Quantity int `json:"quantity" bson:"quantity`
	CreatedAt time.Time `bson:"create_at"`
}

type Item struct{
	ProductName string `json:"product" bson:"product"`
	Price int `json:"price" bson:"price"`
	Description string `json:"desc" bson:"desc"`
	Quantity int `json:"quantity" bson:"quantity"`
}


func InsertItem(c *gin.Context){

	var item Item

	if err := c.ShouldBindJSON(&item);err!=nil{
		c.JSON(400,gin.H{
			"error": "Product Details Messup",
		})
		return
	}
	
	var b_data = B_data{
		Id : primitive.NewObjectID(),
		PublicId : uuid.NewString(),
		ProductName : item.ProductName,
		Description : item.Description,
		Price : item.Price,
		Quantity : item.Quantity,
		CreatedAt : time.Now(),
	}
	getcollection := libs.GetProductDb()
	if _ , err := getcollection.InsertOne(context.Background(),b_data);err!=nil{
		if mongo.IsDuplicateKeyError(err){
			c.JSON(409,gin.H{
				"error":"Product already Found",
			})
			return
		}
		
		c.JSON(500,gin.H{
			"error":"database error",
		})
		return
	}

	c.JSON(200,gin.H{
		"id" :b_data.Id,
		"public_id":b_data.PublicId,
		"state" : true,
		"Message":"Item successfully added",
	})
}


func Products(c *gin.Context){


	
	ctx , cancel := context.WithTimeout(context.Background(),time.Second * 10)
	defer cancel()

	PageStr := c.PostForm("page")

	page , _ := strconv.Atoi(PageStr)

	if page < 1{
		page = 1 
	}

	limit := int64(10)
	skip := int64((page - 1) * 10)

	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip)

	getcollection := libs.GetProductDb()

	cursor , err := getcollection.Find(ctx,bson.M{},opts)

	if err !=nil{
		c.JSON(500,gin.H{
			"error":"database error",
		})
		return
	}

	defer cursor.Close(ctx)
	var users []Item

	cursor.All(ctx,&users)

	c.JSON(200,gin.H{
		"Products" : users,
	})


}



func DelItem(c *gin.Context){
	
	item := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(item)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid id"})
        return
    }

	getcollection := libs.GetProductDb()

	res , err := getcollection.DeleteOne(c,bson.M{"_id":objID})

	if err!=nil{
		c.JSON(500,gin.H{
			"error":"database error",
		})
		return
	}

	if res.DeletedCount == 0{
		c.JSON(409,gin.H{
			"error":"Product notFound",
		})
		return
	}

	c.JSON(200,gin.H{"message":"Product Deleted Successfully ",})
}



func Product(c *gin.Context){
	var p_data B_data
	product_id := c.Param("puid")

	if product_id == ""{
		c.JSON(409,gin.H{
			"error" : "Product credentials Missing",
		})
		return
	}


	getcollection := libs.GetProductDb()

	filter := bson.M{"public_id":product_id}
	if err:=getcollection.FindOne(context.TODO(),filter).Decode(&p_data);err!=nil{
		c.JSON(500,gin.H{
			"error" : "product not found",
		})
		return
	}

	c.JSON(200,gin.H{
		"product" : p_data.ProductName,
		"price" : p_data.Price,
		"description" : p_data.Description,
	})
}


func UpdatePQ(c *gin.Context){

	var item struct{
		PublicId string `json:"gid"`
		Price int `json:"price",omitempty`
		Quantity int  `json:"quantity",omitempty`
	}

	if err := c.ShouldBindJSON(&item);err!=nil{
		c.JSON(400,gin.H{
			"error":"Detail's messup",
		})
		return
	}

	filter := bson.D{{"public_id",item.PublicId}}
	update := mongo.Pipeline{
		{
			{"$set", bson.M{
				"quantity": bson.M{
					"$max": bson.A{
						0,
						bson.M{"$add": bson.A{"$quantity", item.Quantity}},
					},
				},
				"price": bson.M{
					"$max": bson.A{
						0,
						bson.M{"$add": bson.A{"$price", item.Price}},
					},
				},
			}},
		},
	}


	getcollection := libs.GetProductDb()


	res , err := getcollection.UpdateOne(context.Background(),filter,update)

	if err!=nil{
		c.JSON(500,gin.H{
			"error" : err.Error(),
		})
		return
	}
	if res.MatchedCount != 1{
		c.JSON(409,gin.H{
			"error":"Product notfound",
		})
		return
	}

	if res.ModifiedCount != 1 {
		c.JSON(409,gin.H{
			"error":"Product not found or Product already in 0",
		})
		return
	}

	c.JSON(200,gin.H{
		"state" : true,
		"message" : "Price succesfully Updated [id]: "+item.PublicId,
	})
}