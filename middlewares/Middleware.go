package middlewares

import (
	helper "ECommerce-Backend/helper"
	utils "ECommerce-Backend/utils"
	"net/http"
	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleWare(limiter *rate.Limiter)gin.HandlerFunc{
	return func(c *gin.Context){
		if(limiter.Allow()){
			c.Next();
		}else{
			http.Error(c.Writer,"Too many requests,wait for a moment",http.StatusTooManyRequests);
			c.Abort();
		}
	}
}

func Authorize(roles ...string) gin.HandlerFunc{
	return func(c *gin.Context){
		role,exists := c.Get("role");

		utils.LogMessage("Role is ");
		utils.LogMessage(role)
		if !exists{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Unauthorized"});
			return;
		}

		for _,allowedRoles := range roles{
			if role == allowedRoles{
				c.Next();
				return;
			}
		}
	}
}

func Authenticate() gin.HandlerFunc{
	return func (c *gin.Context)  {
		clientToken := c.Request.Header.Get("token")

		if clientToken==""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"No Authorization header provided!"});
			c.Abort();
			return;
		}

		claims,err := helper.ValidateToken(clientToken);

		if err!=""{
			utils.LogMessage("Something went wrong while authenticating!");
			c.JSON(http.StatusInternalServerError,gin.H{"error":err});
			c.Abort();
			return;
		}

		c.Set("email",claims.Email);
		c.Set("first_name",claims.First_name);
		c.Set("last_name",claims.Last_name);
		c.Set("uid",claims.Uid);
		c.Set("role",claims.User_role);
		c.Next();
	}
}