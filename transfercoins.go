package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

type transaction struct {
	Sender    int       `json:"Sender", db:"Sender"`
	Recipient int       `json:"Recipient", db:"Recipient"`
	Coins     float32   `json:"Coins", db:"Coins"`
	Time      time.Time `json:"Time", db:"Time"`
}

func Transfercoins(c *gin.Context) {
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
	var t transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	roll := token.Claims.(jwt.MapClaims)["Roll_NO"]
	tmp := fmt.Sprintf("%v", roll)
	t.Sender, _ = strconv.Atoi(tmp)
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Code to check if Sender's Balance is enough to complete transaction
	row := tx.QueryRow("SELECT Coins FROM users WHERE RollNO=$1", t.Sender)
	var bal float32
	if err := row.Scan(&bal); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	if bal < (t.Coins) {
		c.JSON(http.StatusUnprocessableEntity, "Not enough Balance in Sender's Account to proceed with the Transaction")
		tx.Rollback()
		return
	}

	//Code to check if Sender's have participated in sufficient events
	row = tx.QueryRow("SELECT COUNT(*) FROM history WHERE Action='AwardCoins' AND Recipient=$1", t.Sender)
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

	//Code to check if reciever is permitted to recieve coins
	row = tx.QueryRow("Select Role FROM users WHERE RollNO=$1", t.Recipient)
	var r string
	if err := row.Scan(&r); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	if r != "General" {
		c.JSON(http.StatusUnprocessableEntity, "Recipient is not allowed to earn coin")
		tx.Rollback()
		return
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET Coins=Coins-$1 WHERE RollNO=$2", t.Coins, t.Sender)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Updating Database")
		tx.Rollback()
		return
	}
	var tax float32
	tax = 0.98
	if (t.Sender / 10000) != (t.Recipient / 10000) {
		tax = 0.67
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET Coins = CASE WHEN Coins+($1*$2) <= 100 THEN Coins+($1*$2) ELSE 100 END WHERE RollNO=$3", tax, t.Coins, t.Recipient)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Updating Database")
		tx.Rollback()
		return
	}
	t.Time = time.Now()
	_, err = tx.ExecContext(ctx, "INSERT INTO history (Sender, Recipient, Coins, Time, Action) VALUES($1,$2,$3,$4,'TransferCoins')", t.Sender, t.Recipient, t.Coins, t.Time)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Updating Database")
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
