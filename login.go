package main

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)
var db *sql.DB
type User struct {
	Password string `json:"Password", db:"Password"`
	Username string `json:"Username", db:"Username"`
	RollNO int `json:"RollNO", db:"RollNO"`
}

func Login(c *gin.Context) {
	//fmt.Println("12")
	var u User
	err := c.ShouldBindJSON(&u);
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	db,err = sql.Open("sqlite3","./mydb.db")
	if err!=nil {
		panic(err)
	}
	row := db.QueryRow("SELECT Password FROM users WHERE Username=$1",u.Username)
	var pwd string
	if err:=row.Scan(&pwd); err!=nil{
		c.JSON(http.StatusUnprocessableEntity, "NO user found")
		return
	}


	if err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(u.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		c.JSON(http.StatusUnauthorized,"Provided Password is Incorrect")
		return
	}
	//c.JSON(200, pwd)

	//compare the user from the request, with the one we defined:
	token, err := CreateToken(u.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in creating Token")
		return
	}
	c.JSON(http.StatusOK, token)
}

func CreateToken(username string) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
