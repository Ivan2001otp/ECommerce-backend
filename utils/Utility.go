package utils

import (
	"fmt"
	"math"
)

func LogMessage(message interface{}) {
	fmt.Println(message)
}

func round(num float64) int{
	return int(num+math.Copysign(0.5,num));
}

func TransformToFixed(num float64,precision int) float64{
	output:=math.Pow(10,float64(precision));
	return float64(round((num*output)/output));
}