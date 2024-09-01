package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID 					primitive.ObjectID 			`bson:"_id"`
	Category_id			*string						`json:"category_id" validate:"required"`
	Product_name		string						`json:"product_name" validate:"required,max=100,min=2"`
	Description			string						`json:"description"  validate:"required,max=256"`
	Stocks				int							`json:"stocks" validate:"required"`	
	Price				float64						`json:"price" validate:"required"`
	Product_image		*string						`json:"product_image" validate:"required"`
	Created_at			time.Time					`json:"created_at"`
	Updated_at			time.Time					`json:"updated_at"`
}