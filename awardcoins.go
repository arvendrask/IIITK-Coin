package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func Awardcoins(c *gin.Context) {
	MySigningKey := []byte(os.Getenv("ACCESS"))

	tokenstr := c.Request.Header.Get("Authorization")

	token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (interface{}, error) {
		//c.JSON(200,MySigningKey)
		if t.Method.Alg() != "HS256" {
			return []byte(""), fmt.Errorf("Invalid Signing Method1")
		}

		return MySigningKey, nil
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if !token.Valid {
		c.JSON(http.StatusUnprocessableEntity, "Invalid Token Provided2")
		return
	}
	roll := token.Claims.(jwt.MapClaims)["Roll_NO"]
	tmpr,_ := strconv.Atoi(fmt.Sprintf("%v", roll))
	row := db.QueryRow("SELECT Role FROM users WHERE RollNO=$1", tmpr)
	var role string
	if err := row.Scan(&role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	if role!="admin" {
		c.JSON(http.StatusUnprocessableEntity, "Not Authorized")
		return
	}
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
