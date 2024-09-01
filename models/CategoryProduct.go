package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductCategory struct {
	ID 					primitive.ObjectID		`bson:"_id"`
	Category_id			string					`json:"category_id" validate:"required"`
	Name				string					`json:"name"`
	Created_at			time.Time				`json:"created_at"`
	Updated_at			time.Time				`json:"updated_at"`
	Category 			string					`json:"category"`
}