package main

import (
	"github.com/gin-gonic/gin"
	_ "fmt"
	"fmt"
)

func main() {
	// Disable Console Color
	// gin.DisableConsoleColor()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/get", getting)
	router.POST("/post", posting)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}

func getting(c *gin.Context)  {
	c.String(200,"HelloWorld")
}

type test struct {
	Name string `json:"name"`
	Age int `json:"age"`
}
func posting(c *gin.Context)  {
	var a test
	c.BindJSON(&a)
	c.String(200,fmt.Sprintf("my name is %s and my age is %d",a.Name,a.Age))
	c.JSON(200,a)
}