package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Awardcoins(c *gin.Context) {
	var th transaction
	if err := c.ShouldBindJSON(&th); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	statement, err := db.Prepare("UPDATE users SET Coins = CASE WHEN Coins+$1 <= 100 THEN Coins+$1 ELSE 100 END WHERE RollNO = $2")
	statement.Exec(th.Coins, th.Recipient)
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in updating coin in Database")
		return
	}
	th.Time = time.Now()
	// Admin Roll NO is already known for now is assigned on random
	th.Sender = 111111
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS history (Sender INTEGER , Recipient INTEGER, Coins FLOAT, Time TEXT, Action TEXT)")
	statement.Exec()
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in creating Database")
		return
	}
	tmp := fmt.Sprintf("%v", th.Time)
	statement, err = db.Prepare("INSERT INTO history (Sender, Recipient, Coins, Time, Action) VALUES(?,?,?,?,'AwardCoins')")
	statement.Exec(th.Sender, th.Recipient, th.Coins, tmp)

	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in Inserting User Data in Database")
		return
	}

}
