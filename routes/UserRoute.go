package routes

import (
	middleware "ECommerce-Backend/middlewares"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func UserRoutes(incomingRoutes *gin.Engine){
	rateLimiter := rate.NewLimiter(5,10);//5 requests per second

	incomingRoutes.GET("/users",
		middleware.RateLimitMiddleWare(rateLimiter),
		middleware.Authorize("admin"),
		
		);
	incomingRoutes.POST("/users/signup",middleware.RateLimitMiddleWare(rateLimiter),);
	incomingRoutes.POST("/users/login",middleware.RateLimitMiddleWare(rateLimiter),);
}