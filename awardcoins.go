package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type coin struct{
	RollNO int `json:"RollNO"`
	Coins int `json:"Coins"`
}
func Awardcoins(c* gin.Context){
	var ac coin
	if err := c.ShouldBindJSON(&ac); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	statement, error := db.Prepare("UPDATE users SET Coins = ? WHERE RollNO = ?");
	statement.Exec(ac.Coins, ac.RollNO)
	if error != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in updating coin in Database")
		return
	}

}
