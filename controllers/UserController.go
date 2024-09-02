package controllers

import (
	database "ECommerce-Backend/database"
	helper "ECommerce-Backend/helper"
	"ECommerce-Backend/models"
	"ECommerce-Backend/utils"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		err:= helper.CheckUserType(c,"admin");
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Cannot access /getusers endpoint because u are not admin!"});
			return;
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

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
			utils.LogMessage("Warning Parsing error -> "+err.Error())
		}

		matchStage := bson.D{
			{
				Key: "$match", Value: bson.D{{}},
			},
		}
		
		groupStage := bson.D{
			{
				"$group",bson.D{
					{"_id",bson.D{{"_id","null"}}},
					{"total_count",bson.D{{"$sum",1}}},
					{"data",bson.D{{"$push","$$ROOT"}}},
				},
			},
		}

		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count",Value: 1},
					{Key: "user_items",
						Value: bson.D{
							{Key: "$slice",
							Value:[]interface{}{"$data",startIndex,recordPerPage}},
						},
					},
				},
			},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while listing users"})
			return
		}

		var allUser []bson.M

		if err = result.All(ctx, &allUser); err != nil {
			utils.LogMessage("Something went wrong with GETUSERS()->" + err.Error())
			log.Fatal(err.Error())
			return
		}

		c.JSON(http.StatusOK, allUser[0])

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		validate := validator.New()

		if err := c.BindJSON(&user); err != nil {
			defer cancel()
			msg := "Something went wrong on bindJSON in signUp()"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		validationErr := validate.Struct(&user)

		if validationErr != nil {
			defer cancel()
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": *user.Email})

		if err != nil {
			defer cancel()
			utils.LogMessage("Duplicate email exists - " + *user.Email)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cancel()

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email already exists!"})
			return
		}

		phoneNumCount, err := userCollection.CountDocuments(ctx, bson.M{"phone_number": *user.Phone})

		if err != nil {
			defer cancel()
			utils.LogMessage("Something went wrong while checking duplicate phone number - " + *user.Phone)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if phoneNumCount > 0 {
			defer cancel()
			utils.LogMessage("This pohne num already exists " + *user.Phone)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Phone Number already exists !"})
			return
		}

		hashedPassword := helper.HashPassword(user.Password)
		user.Password = hashedPassword

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, user.First_name, user.Last_name, user.Role, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			utils.LogMessage("user not able to insert.Something went wrong!")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to insert user data ->" + err.Error()})
			return
		}

		utils.LogMessage(result)

		c.JSON(http.StatusOK, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var foundUser models.User
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			defer cancel()
			utils.LogMessage("Something went wrong ->" + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			defer cancel()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Target email not found while logging"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, foundUser)
	}
}
