package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	//fmt.Println("12")
	var u Userdata
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	row := db.QueryRow("SELECT Password FROM users WHERE RollNO=$1", u.RollNO)
	var pwd string
	if err := row.Scan(&pwd); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "NO user found")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(u.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		c.JSON(http.StatusUnauthorized, "Provided Password is Incorrect")
		return
	}
	//c.JSON(200, pwd)

	//compare the user from the request, with the one we defined:
	token, err := CreateToken(u.RollNO)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in creating Token")
		return
	}
	c.JSON(http.StatusOK, token)
}

func CreateToken(Roll int) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["Roll_NO"] = Roll
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS")))
	if err != nil {
		return "", err
	}
	return token, nil
}
