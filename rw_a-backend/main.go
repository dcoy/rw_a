package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.DisableConsoleColor()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// client := github.NewClient(nil)
	// 	token := os.Getenv("GITHUB_API_TOKEN")

	// issues := client.

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"path": "path", "value": "pong"})
	})

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"key": "key", "value": "value"})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "GH Issues Renderer",
			},
		)
	})

	// Get user value
	// r.GET("/user/:name", func(c *gin.Context) {
	// 	user := c.Params.ByName("name")
	// 	value, ok := db[user]
	// 	if ok {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	// 	}
	// })

	// // Authorized group (uses gin.BasicAuth() middleware)
	// // Same than:
	// // authorized := r.Group("/")
	// // authorized.Use(gin.BasicAuth(gin.Credentials{
	// //	  "foo":  "bar",
	// //	  "manu": "123",
	// //}))
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
