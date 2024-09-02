package helper

import (
	database "ECommerce-Backend/database"
	"ECommerce-Backend/utils"

	jwt "github.com/dgrijalva/jwt-go"
)


var SECRET_KEY string = "A7d45Rz9pQ2wVb8Xs";

type SignedDetails struct{
	Email string
	First_name string
	Last_name string
	Uid string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user");

func GenerateAllTokens(email string ,first_name string ,last_name string, user_type string,uid string)(signedToken string,refreshToken string,err Error){
	claims := &SignedDetails{
		Email:email,
		First_name:first_name,
		Last_name:last_name,
		Uid:uid,
		StandardClaims:jwt.StandardClaims{
			//ideal expiration should be kept not more than 30mins.
			ExpiresAt:time.Now().Local().Add(time.Hour*time.Duration(24)).Unix()
		},
	}

	refresh_claims := &SignedDetails{
		StandardClaims:jwt.StandardClaims{
			ExpiresAt:time.Now().Local().Add(time.Hour * time.Duration(48)).Unix(),
		},
	}

	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY));

	if err!=nil{
		utils.LogMessage("Something went wrong while generating  tokens");
		log.Panic(err);
		return;
	}

	refresh_token,err2 := jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY));
	if err2!=nil{
		utils.LogMessage("Something went wrong while generating  refresh tokens")
		log.Panic(err2);
		return;
	}

	return token,refresh_token,nil;
}

func UpdateAllTokens(signedToken string,signedRefreshToken string,userId string){
	var  ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var updateObj primitive.D;

	Updated_at,_ := time.Parse(time.RFC3339,time.Now().Format(time.RFC3339));
	updateObj = append(updateObj,bson.E{"updated_at",Updated_at});
	upsert := true;

	filter := bson.M{"user_id",userId};

	option := options.UpdateOptions{
		Upsert:&upsert,
	}

	_,err := userCollection.UpdateOne{
		ctx,filter,
		bson.D{
			{"$set",updateObj},
		},
		&option,
	}

	defer cancel();

	if err!=nil{
		utils.LogMessage("Something went wrong while updating tokens!");
		log.Panic(err.Error());
		return;
	}

	return;
}

func ValidateToken(signedToken string)(claims *SignedDetails,msg string){
	token,err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token)(interface{},error){
			return []byte(SECRET_KEY),nil;
		},
	)

	if err!=nil{
		utils.LogMessage("Jwt parsewithClaims throwed error!");
		return;
	}

	claims,ok := token.Claims.(*SignedDetails);

	if !ok{
		msg = fmt.Sprintf("The token is invalid");
		utils.LogMessage(msg);
		return;
	}

	if 	claims.ExpiresAt<time.Now().Local().Unix(){
		msg= "Token is expired";
		return;
	}

	return claims,msg;
}