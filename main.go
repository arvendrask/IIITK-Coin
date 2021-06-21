package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB
var(
	router = gin.Default()
)
func main() {
	var err error
	db,err = sql.Open("sqlite3","./mydb.db")
	if err!=nil {
		panic(err)
	}
	router.POST("/signup", Signup )
	router.POST("/login",Login)
	router.GET("/secretpage",Secretpage)
	router.POST("/awardcoins", Awardcoins)
	router.GET("/getcoins", Getcoins)
	router.POST("/transfercoins",Transfercoins)
	log.Fatal(router.Run(":8080"))


}