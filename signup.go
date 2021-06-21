package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Userdata struct {
	Password string `json:"Password", db:"Password"`
	Username string `json:"Username", db:"Username"`
	RollNO   int    `json:"RollNO", db:"RollNO"`
}

func Signup(c *gin.Context) {
	var u Userdata
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	statement, error := db.Prepare("CREATE TABLE IF NOT EXISTS users (Username TEXT, Password TEXT, RollNO INTEGER PRIMARY KEY, Coins INTEGER )")
	statement.Exec()
	if error != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in creating Database")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 16)
	statement, error = db.Prepare("INSERT INTO users (Username, Password, RollNO, Coins) VALUES(?,?,?,0)")
	statement.Exec(u.Username, string(hashedPassword), u.RollNO)

	if err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		c.JSON(http.StatusInternalServerError, "Error in Inserting User Data in Database")
		return
	}
	c.JSON(200, u)
}
