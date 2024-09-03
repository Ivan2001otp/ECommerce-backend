package controllers

import (
	database "ECommerce-Backend/database"
	models "ECommerce-Backend/models"
	utils "ECommerce-Backend/utils"
	helper "ECommerce-Backend/helper"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// dummy image.
const freeNetWorkImageUrl = "https://cdn.pixabay.com/photo/2023/11/03/11/40/technology-8362813_640.jpg"

var productCollection *mongo.Collection = database.OpenCollection(database.Client, "product")

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		err1:= helper.CheckUserType(c,"admin");
		if err1!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Cannot access /getusers endpoint because u are not admin!"});
			return;
		}

		var product models.Product
		var belongedCategory models.ProductCategory

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.BindJSON(&product); err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Binding json while creating product threw error!"})
			return
		}

		validate := validator.New()
		validationErr := validate.Struct(product)

		if validationErr != nil {
			defer cancel()
			utils.LogMessage(validationErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			return
		}

		err := categoryCollection.FindOne(ctx, bson.M{"category_id": product.Category_id}).Decode(&belongedCategory)

		if err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The category not found !"})
			return
		}

		product.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		product.ID = primitive.NewObjectID()
		product.Product_id = product.ID.Hex()

		var num = utils.TransformToFixed(product.Price, 2)
		product.Price = num

		result, insertErr := productCollection.InsertOne(ctx, product)

		if insertErr != nil {
			defer cancel()
			utils.LogMessage(insertErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product!"})
			return
		}

		defer cancel()

		utils.LogMessage(result)

		c.JSON(http.StatusOK, result)

	}
}

func GetProductById() gin.HandlerFunc {
	return func(c *gin.Context) {
		err1:= helper.CheckUserType(c,"admin");
		if err1!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Cannot access /getusers endpoint because u are not admin!"});
			return;
		}

		productId := c.Param("product_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var product models.Product

		err := productCollection.FindOne(ctx, bson.M{"product_id": productId}).Decode(&product)

		if err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product do not exists!"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, product)
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {

		err1:= helper.CheckUserType(c,"admin");
		if err1!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Cannot access /getusers endpoint because u are not admin!"});
			return;
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))

		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage //skip limit
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		if err != nil {
			defer cancel()
			utils.LogMessage("Warning Parsing error -> " + err.Error())
		}

		matchStage := bson.D{{"$match", bson.D{{}}}} //take everything.
		groupStage := bson.D{
			{
				"$group", bson.D{
					{"_id", bson.D{{"_id", "null"}}},
					{"total_count", bson.D{{"$sum", 1}}},
					{"data", bson.D{{"$push", "$$ROOT"}}},
				},
			},
		}

		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{"_id", 0},
					{"total_count", 1},
					{Key: "product_items",
						Value: bson.D{
							{"$slice", []interface{}{"$data", startIndex, recordPerPage}},
						}},
				},
			},
		}
		result, err := productCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})

		defer cancel()

		if err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while listing all products"})
			return
		}

		var allProducts []bson.M

		if err = result.All(ctx, &allProducts); err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong during fetching all products!"})

			return
		}

		defer cancel()
		c.JSON(http.StatusOK, allProducts[0])

	}
}

func UpdateProductById() gin.HandlerFunc {
	return func(c *gin.Context) {
		err1:= helper.CheckUserType(c,"admin");
		if err1!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Cannot access /getusers endpoint because u are not admin!"});
			return;
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		product_id := c.Param("product_id")

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Binding JSON failed while updating product!"})
			return
		}

		filter := bson.M{"product_id": product_id}

		var updateObj primitive.D

		if product.Product_name != "" {
			updateObj = append(updateObj, bson.E{"product_name", product.Product_name})

		}
		if product.Description != "" {
			updateObj = append(updateObj, bson.E{"description", product.Description})
		}
		if product.Price > 0 {
			updateObj = append(updateObj, bson.E{"price", product.Price})
		}
		if *product.Product_image != "" {
			updateObj = append(updateObj, bson.E{"product_image", product.Product_image})
		}
		if product.Stocks >= 0 {
			updateObj = append(updateObj, bson.E{"stocks", product.Stocks})
		}

		product.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateObj = append(updateObj, bson.E{"updated_at", product.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := productCollection.UpdateOne(ctx, filter,
			bson.D{
				{
					"$set", updateObj,
				},
			}, &opt,
		)

		if err != nil {
			defer cancel()
			utils.LogMessage(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falied to update product!"})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
