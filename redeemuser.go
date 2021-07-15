package main

import (
	//"fmt"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type itemlist struct {
	Item   string  `json:"Item", db:"Item"`
	Price  float32 `json:"Price", db:"Price"`
	RollNO int     `json:"RollNO", db:"RollNO"`
	ID     int     `json:"ID", db:"ID"`
	Action int     `json:"Action", db:"Action"`
}

func Redeemuser(c *gin.Context) {
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
	var i itemlist
	roll := token.Claims.(jwt.MapClaims)["Roll_NO"]
	tmp := fmt.Sprintf("%v", roll)
	i.RollNO, _ = strconv.Atoi(tmp)
	//Input JSON shuould have only Item name and price
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS redeemlog (ID INTEGER PRIMARY KEY , RollNO INTEGER, Item TEXT, Price FLOAT, Status TEXT )")
	statement.Exec()
	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in creating Database")
		return
	}
	statement, err = db.Prepare("INSERT INTO redeemlog (RollNO, Item, Price, Status) VALUES(?,?,?,'Pending')")
	statement.Exec(i.RollNO, i.Item, i.Price)

	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in Inserting User Data in Database")
		return
	}
}
