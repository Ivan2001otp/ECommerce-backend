package routes

import (
	controller "ECommerce-Backend/controllers"
	middleware "ECommerce-Backend/middlewares"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func UserRoute(incomingRoutes *gin.Engine) {
	rateLimiter := rate.NewLimiter(5, 10)

	// checks authorization.
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.POST("/order/:user_id", middleware.RateLimitMiddleWare(rateLimiter), controller.CreateOrder())
	incomingRoutes.GET("/order/:order_id", middleware.RateLimitMiddleWare(rateLimiter), controller.FetchOrderById())

	incomingRoutes.GET("/admin/orders",
	middleware.RateLimitMiddleWare(rateLimiter),controller.FetchAllOrders())

}
