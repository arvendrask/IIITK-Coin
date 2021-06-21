package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
func Getcoins(c* gin.Context){
	var gc coin
	if err := c.ShouldBindJSON(&gc); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	row := db.QueryRow("SELECT Coins FROM users WHERE RollNO=$1",gc.RollNO)
	if err:=row.Scan(&(gc.Coins)); err!=nil{
		c.JSON(http.StatusUnprocessableEntity, "NO such user exist in Database")
		return
	}
	c.JSON(200,gc)
}
