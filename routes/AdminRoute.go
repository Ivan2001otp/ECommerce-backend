package routes

import (
	middleware "ECommerce-Backend/middlewares"
	controller "ECommerce-Backend/controllers"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AdminRoutes(incomingRoutes *gin.Engine){
	rateLimiter := rate.NewLimiter(5,10);//5 requests per second

	incomingRoutes.Use(middleware.Authenticate());

	incomingRoutes.GET("/users",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.GetUsers(),
		);

}