package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID 					primitive.ObjectID   `bson:"_id"`
	Invoice_id			string				 `json:"invoice_id"`
	Description			*string	/*optional*/			 `json:"description"`
	User_id				*string				 `json:"user_id" validate:"required"`	
	Order_id			*string			  	 `json:"order_id" validate:"required"`
	Created_at			time.Time			`json:"created_at"`
}