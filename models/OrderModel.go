package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID       			primitive.ObjectID		`bson:"_id"`
	Order_id 			string					`json:"order_id"`
	User_id				*string					`json:"user_id" validate:"required"`
	Order_date			time.Time				`json:"order_date" validate:"required"`
	Ordered_quantity	int						`json:"ordered_quantity"`
	Payment_mode		*string					`json:"payment_mode" validate:"eq=CARD|eq=CASH|eq="`
	Payment_status		*string					`json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`
	Ordered_products	[]map[Product]int		`json:"ordered_products"`
	Created_at			time.Time				`json:"created_at"`
	Updated_at			time.Time				`json:"updated_at"`
}