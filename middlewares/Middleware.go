package middlewares

import (
	helper "ECommerce-Backend/helper"
	utils "ECommerce-Backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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

		

		// if cookieErr!=nil{
		// 	utils.LogMessage(cookieErr.Error());
		// 	c.JSON(http.StatusInternalServerError,gin.H{"error":"No Authorization header provided!"});
		// 	c.Abort();
		// 	return;
		// }
		// if clientToken==""{
		// 	c.JSON(http.StatusInternalServerError,gin.H{"error":"No Authorization header provided!"});
		// 	c.Abort();
		// 	return;
		// }

		
		

		claims,err := helper.ValidateToken(clientToken);

		if claims.ExpiresAt < time.Now().Unix(){
			utils.LogMessage("Token is expired!");
			c.AbortWithStatus(http.StatusUnauthorized);
			return;
		}

		if err!=""{
			utils.LogMessage("Something went wrong while authenticating!");
			c.JSON(http.StatusInternalServerError,gin.H{"error":err});
			c.Abort();
			return;
		}

		origin := c.Request.Header.Get("Origin");
		referer := c.Request.Header.Get("Referer");

		if origin!="" && origin!=c.Request.Host && referer!="" && referer!=c.Request.Host{
			utils.LogMessage("something wrong with origin and referer!");
			c.AbortWithStatus(http.StatusForbidden);
			return;
		}

		// xsrfHeader := c.Request.Header.Get("X-XSRF-TOKEN");
		
		// if xsrfHeader!="" && xsrfHeader!=claims.{
		// 	c.AbortWithStatus(http.StatusForbidden);
		// 	return;
		// }

		c.Set("email",claims.Email);
		c.Set("first_name",claims.First_name);
		c.Set("last_name",claims.Last_name);
		c.Set("uid",claims.Uid);
		c.Set("role",claims.User_role);
		c.Next();
	}
}