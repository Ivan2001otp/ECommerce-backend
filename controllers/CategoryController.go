package controllers

import (
	database "ECommerce-Backend/database"
	models "ECommerce-Backend/models"
	utils "ECommerce-Backend/utils"
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

var categoryCollection *mongo.Collection = database.OpenCollection(database.Client,"category");

func AddCategory() gin.HandlerFunc{
	return func (c *gin.Context)  {
			var productCategory models.ProductCategory;

			var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

			if err:= c.BindJSON(&productCategory);err!=nil{
				defer cancel();
				utils.LogMessage("something went wrong while binding on adding category")
				c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()});
				return;
			}
			validate := validator.New();

			validationErr := validate.Struct(productCategory);

			if validationErr!=nil{
				defer cancel();
				utils.LogMessage("validation in create category went wront!");
				c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()});
				return;
			}

			productCategory.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339));
			productCategory.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339));

			productCategory.ID = primitive.NewObjectID();

			productCategory.Category_id = productCategory.ID.Hex();

			result,err := categoryCollection.InsertOne(ctx,productCategory);

			defer cancel();

			if err!=nil{
				utils.LogMessage(err.Error());
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Cannot insert newly created category"});
				return;
			}

			c.JSON(http.StatusOK,result);
		}

}

func GetAllCategory() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);
		defer cancel();

		result,err := categoryCollection.Find(context.TODO(),bson.M{});

		if err!=nil{
			utils.LogMessage("Failed to fetch all categories");
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()});
			return;
		}

		var allCategories []bson.M;

		if err=result.All(ctx,&allCategories);err!=nil{
			utils.LogMessage("all-Categories threw error!");
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()});
			return;
		}

		c.JSON(http.StatusOK,allCategories);

	}
}

func GetCategoryById() gin.HandlerFunc{
	return func(c *gin.Context) {
			var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

			categoryId := c.Param("category_id");

			var productCategory models.ProductCategory;

			err := categoryCollection.FindOne(ctx,bson.M{"category_id":categoryId}).Decode(&productCategory);

			if err!=nil{
				defer cancel();
				utils.LogMessage(err.Error());
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Error while fetching category by ID"});
				return;
			}

			c.JSON(http.StatusOK,productCategory);

	}
}

func UpdateCategoryById() gin.HandlerFunc{
	return func(c *gin.Context) {
			var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

			categoryId := c.Param("category_id");

			var productCategory models.ProductCategory;

			if err:= c.BindJSON(&productCategory);err!=nil{
				defer cancel();
				utils.LogMessage(err.Error());
				c.JSON(http.StatusInternalServerError,gin.H{"error":"updateCategoryById - failed to bind go struct!"});
				return;	
			}

			filter := bson.M{"category_id":categoryId};

			var updateObj primitive.D;


			if productCategory.Category!=""{
				updateObj = append(updateObj, bson.E{"category",productCategory.Category});
			}

			if productCategory.Name!=""{
				updateObj = append(updateObj, bson.E{"name",productCategory.Name});
			}

			productCategory.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339));

			updateObj = append(updateObj, bson.E{"updated_at",productCategory.Updated_at});

			upsert:=true;

			opt := options.UpdateOptions{
				Upsert:&upsert,
			}

			result ,err := categoryCollection.UpdateOne(ctx,filter,
				bson.D{
					{
						"$set",updateObj,
					},
				},
				&opt,
			);

			if err!=nil{
				defer cancel();
				utils.LogMessage(err.Error());
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to update category !"});
				return;
			}

			c.JSON(http.StatusOK,result);

	}
}