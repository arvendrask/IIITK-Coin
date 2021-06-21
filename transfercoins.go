package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
type transaction struct{
	Sender int `json:"Sender"`
	Recipient int `json:"Recipient"`
	Coins int `json:"Coins", db:"Coins"`
}
func Transfercoins(c* gin.Context){
	var t transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	ctx := context.Background()
	tx,err:= db.BeginTx(ctx,nil)
	if err != nil {
		log.Fatal(err)
	}
	row:= tx.QueryRow("SELECT Coins FROM users WHERE RollNO=$1",t.Sender)
	var bal int
	if err:=row.Scan(&bal); err!=nil{
		c.JSON(http.StatusUnprocessableEntity, "NO such Sender exist in databse")
		return
	}
	if bal < (t.Coins) {
		c.JSON(http.StatusUnprocessableEntity, "Not enough Balance in Sender's Account to proceed with the Transaction")
		tx.Rollback()
		return
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET Coins=Coins-$1 WHERE RollNO=$2", t.Coins, t.Sender)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Updating Database")
		tx.Rollback()
		return
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET Coins=Coins+$1 WHERE RollNO=$2", t.Coins, t.Recipient)
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
