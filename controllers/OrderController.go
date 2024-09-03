package controllers

import (
	database "ECommerce-Backend/database"
	models "ECommerce-Backend/models"
	utils "ECommerce-Backend/utils"
	helper "ECommerce-Backend/helper"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func CreateOrder() gin.HandlerFunc{
	return (c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);
		var order models.Order;


		if err=c.BindJSON(&order);err!=nil{
			defer cancel();
			utils.LogMessage(err.Error());
			c.JSON(http.StatusBadRequest,gin.H{"error":"Failed to parse Order!"});
			return;
		}

		validate := validator.New();
		validationErr := validate.Struct(order);

		if validationErr!=nil{
			defer cancel();
			utils.LogMessage(validationErr.Error());
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()});
			return;
		}

		order.Updated_at,_ := time.Parse(time.RFC3339,time.Now().Format(time.RFC3339));
		order.Created_at,_ := time.Parse(time.RFC3339,time.Now().Format(time.RFC3339));
		
		order.ID = primitive.NewObjectID();

		order.Order_id = order.ID.Hex();
		

		//calculating totalprice.
		 total_price := 0.0;

		for _,productMap :=range order.Ordered_products{
			for product,quantity := range productMap{
				total_price += product.Price*float64(quantity);
			}
		}

		order.Total_price = utils.TransformToFixed(total_price);

		result,err := productCollection.InsertOne(ctx,order);
		defer cancel();

		if err!=nil{
			utils.LogMessage(err.Error());
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to insert new order !"});
			return;
		}

		utils.LogMessage(result);

		c.JSON(http.StatusOK,result);

	}
}

func FetchAllOrders() gin.HandlerFunc{
	return (c *gin.Context){

	}
}

func FetchOrderById() gin.HandlerFunc{
	return (c *gin.Context){

	}
}

func UpdateOrderById() gin.HandlerFunc{
	return (c *gin.Context){

	}
}