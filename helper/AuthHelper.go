package helper

import (
	"ECommerce-Backend/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context,role string)(err error){
	userType := c.GetString("role")
	utils.LogMessage(c);
	
	err = nil;

	if userType==""{
		utils.LogMessage("The user Type is empty !");
	}else{
		utils.LogMessage("The user Type is not empty !");

	}

	utils.LogMessage("The role is "+userType);

	if userType!=role{
		err = errors.New("Unauthorized to access this resources");
		return err;
	}
	return err;
}