package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
)
var(
	router = gin.Default()
)
func main() {
	//router.POST("/signup", Signup )
	//log.Fatal(router.Run(":8080"))

	router.POST("/login",Login)
	router.GET("/secretpage",Secretpage)
	log.Fatal(router.Run(":8080"))

	//router.POST("/secretpage",Secretpage)
	//log.Fatal(router.Run(":8080"))



}