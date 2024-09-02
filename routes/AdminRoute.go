package routes

import (
	middleware "ECommerce-Backend/middlewares"
	controller "ECommerce-Backend/controllers"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AdminRoutes(incomingRoutes *gin.Engine){
	rateLimiter := rate.NewLimiter(5,10);//5 requests per second

	//checks authorization.
	incomingRoutes.Use(middleware.Authenticate());
	
	incomingRoutes.GET("/users",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.GetUsers(),
		);

	incomingRoutes.POST("/admin/product/add",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.AddProduct(),
	)

	incomingRoutes.POST("/admin/category/add",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.AddCategory(),
	)

}