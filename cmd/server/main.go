package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define routes (endpoints)

	// todo 1. Create repo (placeholder)
	router.POST("/repos", func(c *gin.Context) {
		c.JSON(201, gin.H{"message": "repo created (placeholder)"})
	})

	// todo 2. Delete repo (placeholder)
	router.DELETE("/repos/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(200, gin.H{"message": "repo deleted (placeholder)", "repo": name})
	})

	// todo 3. List all repos (placeholder)
	router.GET("/repos", func(c *gin.Context) {
		sampleRepos := []string{"repo1", "repo2", "repo3"}
		c.JSON(200, gin.H{"repositories": sampleRepos})
	})

	// todo 4. List N open pull requests for a repo (placeholder)
	router.GET("/repos/:name/pulls", func(c *gin.Context) {
		name := c.Param("name")
		n := c.Query("n")
		placeholderPRs := []string{"pr1", "pr2", "pr3"}
		c.JSON(200, gin.H{"repo": name, "pull_requests": placeholderPRs, "count_requested": n})
	})

	// Start server
	fmt.Println("Server running on port http://localhost:8080")
	router.Run(":8080")
}
