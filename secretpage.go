package main

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/twinj/uuid"
	"net/http"
	"os"
)

type info struct {
	Password string `json:"Password", db:"Password"`
	Username string `json:"Username", db:"Username"`
	RollNO int `json:"RollNO", db:"RollNO"`
}
var MySigningKey = []byte(os.Getenv("ACCESS_SECRET"))
var u info

func Secretpage(c *gin.Context){
	err := c.ShouldBindJSON(&u);
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	f:=0
	if c.Request.Header.Get("Authorization") != "" {

		token, err := jwt.Parse(string(c.Request.Header.Get("Authorization")[0]), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(("Invalid Signing Method"))
			}
			aud := u.Username
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAudience {
				return nil, fmt.Errorf(("invalid aud"))
			}
			/*	// verify iss claim
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf(("invalid iss"))
				}*/

			return MySigningKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
		}

		if token.Valid {
			f=1
		}

	} else {
		c.JSON(http.StatusUnprocessableEntity, "No Authorization Token provided")
	}
	if f==1 {
		db, err = sql.Open("sqlite3", "./mydb.db")
		if err != nil {
			panic(err)
		}
		row := db.QueryRow("SELECT Password FROM users WHERE Username=$1", u.Username)
		if err := row.Scan(&u.RollNO); err != nil {
			c.JSON(http.StatusUnprocessableEntity, "NO user found")
			return
		}
		c.JSON(200, u)
	}
}
