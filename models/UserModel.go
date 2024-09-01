package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID			 primitive.ObjectID		`bson:"_id"`
	User_id		 string 				`json:"user_id"`
	First_name		 string					`json:"first_name" validate:"required,min=2,max=100"`
	Last_name		 string					`json:"last_name" validate:"required,max=100"`
	Email		 *string				`json:"email" validate:"email,required,min=6,max=50"`
	Phone		 *string				`json:"phone" validate:"max=10"`
	Address		 *string				`json:"address" validate:"max=1024"`
	Password	 string 				`json:"password" validate:"required,max=10,min=3"`
	Created_at 	time.Time				`json:"created_at"`
	Updated_at	time.Time				`json:"updated_at"`
	
}	