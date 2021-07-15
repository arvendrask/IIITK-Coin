package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Redeemadminlist(c *gin.Context) {
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
	rows, err := db.Query("Select * from redeemlog")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error in Database")
		return
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	c.JSON(200, tableData)

}
