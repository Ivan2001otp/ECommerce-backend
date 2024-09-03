package routes

import (
	middleware "ECommerce-Backend/middlewares"
	controller "ECommerce-Backend/controllers"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	rateLimiter := rate.NewLimiter(5,10);//5 requests per second

	
	
	incomingRoutes.POST("/users/signup",middleware.RateLimitMiddleWare(rateLimiter),controller.SignUp());
	incomingRoutes.POST("/users/login",middleware.RateLimitMiddleWare(rateLimiter),controller.Login());
}