package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"
)

func Getcoins(c *gin.Context) {
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
	var balance float32
	roll := token.Claims.(jwt.MapClaims)["Roll_NO"]
	tmp, _ := strconv.Atoi(fmt.Sprintf("%v", roll))
	row := db.QueryRow("SELECT Coins FROM users WHERE RollNO=$1", tmp)
	if err := row.Scan(&(balance)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "NO such user exist in Database")
		return
	}
	c.JSON(200, balance)
}
