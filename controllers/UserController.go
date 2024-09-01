package controllers

import (
	database "ECommerce-Backend/database"
	"ECommerce-Backend/models"
	"ECommerce-Backend/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		startIndex, err3 := strconv.Atoi(c.Query("startIndex"))

		if err3 != nil {
			defer cancel();
			utils.LogMessage("something wrong with pagination")
			return
		}

		matchStage := bson.D{
			{
				Key: "$match", Value: bson.D{{}},
			},
		}

		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: bson.D{{"$size", "$user_items"}}},
					{Key: "user_items",
						Value: bson.D{
							{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
						},
					},
				},
			},
		}

		result,err := userCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage,
			projectStage,
		})

		defer cancel();


		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error while listing users"});
			return;
		}

		var allUser []bson.M;

		if err = result.All(ctx,&allUser);err!=nil{
			utils.LogMessage("Something went wrong with GETUSERS()->"+err.Error())
			log.Fatal(err.Error());
			return;
		}

		c.JSON(http.StatusOK,allUser);

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

		var foundUser models.User;
		var user models.User;

		if err:=c.BindJSON(&user);err!=nil{
			defer cancel();
			utils.LogMessage("Something went wrong ->"+err.Error());
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()});
			return;
		}

		err := userCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser);

		if err!=nil{
			defer cancel();
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Target email not found while logging"});
			return;
		}
	}
}


func HashPassword(password string) string{
	bytes,err := bcrypt.GenerateFromPassword([]byte(password),14);
	if err!=nil{
		utils.LogMessage("hashing went wrong !");
		log.Panic(err);
		return "";
	}

	return string(bytes);
}

func VerifyPassword(userPassword string,providedPassword string)(bool,error){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword),[]byte(userPassword));

	check:=true;
	msg:=""

	if err!=nil{
		msg = fmt.Sprintf("Login email or password is incorrect!")
		check =false;
		utils.LogMessage(err.Error());
		return check,err;
	}

	return check,nil;
}