package controllers

import (
	database "ECommerce-Backend/database"
	models "ECommerce-Backend/models"
	utils "ECommerce-Backend/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order

		if err := c.BindJSON(&order); err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse Order!"})
			return
		}

		validate := validator.New()
		validationErr := validate.Struct(order)

		if validationErr != nil {
			defer cancel()
			utils.LogMessage(validationErr.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()

		order.Order_id = order.ID.Hex()

		order.Order_date = &order.Created_at

		//calculating totalprice.
		total_price := 0.0

		for _, productMap := range order.Ordered_products {

			total_price += productMap.Added_product.Price * float64(productMap.Quantity)

		}

		fmt.Println(total_price)
		order.Total_price = &total_price

		result, err := orderCollection.InsertOne(ctx, order)
		defer cancel()

		if err != nil {
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new order !"})
			return
		}

		utils.LogMessage(result)

		c.JSON(http.StatusOK, result)
	}
}

// this is only allowed to admin
func FetchAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// err := helper.CheckUserType(c, "admin")


		// if err != nil {
		// 	defer cancel()
		// 	utils.LogMessage(err.Error())
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "FetchAllOrders endpoint is only accessed to admin user type!"})
		// 	return
		// }

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))

		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		if err != nil {
			utils.LogMessage("Warning parsing error->" + err.Error())
		}

		matchStage := bson.D{
			{
				"$match", bson.D{{}},
			},
		}

		groupStage := bson.D{
			{
				"$group", bson.D{
					{"_id", "null"},
					{"total_count", bson.D{{"$sum", 1}}},
					{"data", bson.D{{"$push", "$$ROOT"}}},
				},
			},
		}

		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{
						"_id", 0,
					},
					{"total_count", 1},
					{Key: "order_items", Value: bson.D{
						{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
					}},
				},
			},
		}

		result, err := orderCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})

		defer cancel()

		if err != nil {
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the orders from db!"})
			return
		}

		var allOrders []bson.M

		if err = result.All(ctx, &allOrders); err != nil {
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse allOrders in golang struct!"})
			return
		}

		c.JSON(http.StatusOK, allOrders[0])

	}
}

func FetchOrderById() gin.HandlerFunc {

	return func(c *gin.Context) {
		// err := helper.CheckUserType(c, "admin")

		// if err != nil {
		// 	utils.LogMessage(err.Error())
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "FetchOrderById endpoint is only accessible to admin!"})
		// 	return
		// }

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderId := c.Param("order_id")
		var order models.Order
		fmt.Println("Order id -> "+orderId);

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)

		defer cancel()

		if err != nil {

			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while fetching order by id!"})
			return
		}

		c.JSON(http.StatusOK, order)

	}
}
