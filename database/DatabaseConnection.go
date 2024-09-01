package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	util "ECommerce-Backend/utils"
)

func DBinstance() *mongo.Client{
	connectionUrl := "mongodb://localhost:27017";
	fmt.Println(connectionUrl);

	client,err := mongo.NewClient(options.Client().ApplyURI(connectionUrl));

	if err !=nil{
		util.LogMessage("Invalid connection url to Mongo->"+err.Error());
		log.Fatal(err)
		
	}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second);

	defer cancel();

	err = client.Connect(ctx)

	if err!=nil{
		util.LogMessage("Something went wrong after during mongo connection ! "+err.Error());
		log.Fatal(err)

	}

	util.LogMessage("Connected to mongo successfully!");
	return client;
}