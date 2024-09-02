package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	routes "ECommerce-Backend/routes"
)

func main(){
	fmt.Println("Hi");
	port := "8080";
	router:=gin.New();

	router.Use(gin.Logger());

	routes.AuthRoutes(router);
	routes.AdminRoutes(router);

	router.Run(":"+port);
}