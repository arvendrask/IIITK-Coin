package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Redeemadminaction(c *gin.Context) {
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
	tmpr, _ := strconv.Atoi(fmt.Sprintf("%v", roll))
	row := db.QueryRow("SELECT Role FROM users WHERE RollNO=$1", tmpr)
	var role string
	if err := row.Scan(&role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	if role != "admin" {
		c.JSON(http.StatusUnprocessableEntity, "Not Authorized")
		return
	}
	var i itemlist
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	row = tx.QueryRow("SELECT Coins FROM users WHERE RollNO=$1", i.RollNO)
	var bal float32
	if err := row.Scan(&bal); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	if bal < (i.Price) {
		c.JSON(http.StatusUnprocessableEntity, "Not enough Balance in Sender's Account to proceed with the Transaction")
		tx.Rollback()
		return
	}
	row = tx.QueryRow("SELECT COUNT(*) FROM history WHERE Action='AwardCoins' AND Recipient=$1", i.RollNO)
	var reqev int
	if err := row.Scan(&reqev); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	//For now minimum no of event required is 1
	if reqev < 1 {
		c.JSON(http.StatusUnprocessableEntity, "Sender does not have required number of event participation to proceed with the Transaction")
		tx.Rollback()
		return
	}
	if i.Action == 0 {
		_, err = tx.ExecContext(ctx, "UPDATE redeemlog SET Status='Rejected' WHERE RollNO=$1", i.RollNO)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error in Updating Database")
			tx.Rollback()
			return
		}
	} else {
		_, err = tx.ExecContext(ctx, "UPDATE redeemlog SET Status='Approved' WHERE RollNO=$1", i.RollNO)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error in Updating Database")
			tx.Rollback()
			return
		}
		_, err = tx.ExecContext(ctx, "UPDATE users SET Coins=Coins-$1 WHERE RollNO=$2", i.Price, i.RollNO)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error in Updating Database")
			tx.Rollback()
			return
		}

	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
