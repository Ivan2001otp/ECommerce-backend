package routes

import (
	controller "ECommerce-Backend/controllers"
	middleware "ECommerce-Backend/middlewares"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	rateLimiter := rate.NewLimiter(5, 10) //5 requests per second

	//checks authorization.
	incomingRoutes.Use(middleware.Authenticate())

	//user ops
	incomingRoutes.GET("/users",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.GetUsers(),
	)

	
	//product ops
	incomingRoutes.POST("/admin/product/add",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.AddProduct(),
	)

	incomingRoutes.GET("/admin/products", middleware.RateLimitMiddleWare(rateLimiter),
		controller.GetAllProducts(),
	)

	incomingRoutes.GET("/admin/products/:product_id", middleware.RateLimitMiddleWare(rateLimiter),
		controller.GetProductById(),
	)

	incomingRoutes.PATCH("/admin/product/:product_id", middleware.RateLimitMiddleWare(rateLimiter),
		controller.UpdateProductById(),
	)


	//category ops
	incomingRoutes.POST("/admin/category/add",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.AddCategory(),
	)

	incomingRoutes.GET("/admin/categories",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.GetAllCategory())

	incomingRoutes.GET("/admin/categories/:category_id",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.GetCategoryById())

	incomingRoutes.PATCH("/admin/categories/:category_id",
		middleware.RateLimitMiddleWare(rateLimiter),
		controller.UpdateCategoryById(),
	)

}
