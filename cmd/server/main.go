package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jorgebaptista/octo-manager/internal/githubapi"
)

func main() {
	// Initialize GitHub client
	ghClient, err := githubapi.NewClient()
	if err != nil {
		log.Fatalf("Failed to create GitHub client: %v", err)
	}

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

	// List all repos (placeholder)
	router.GET("/repos", func(c *gin.Context) {
		ctx := context.Background()
		repos, err := ghClient.ListRepos(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// todo extract only the names for now
		var repoNames []string
		for _, repo := range repos {
			if repo.Name != nil {
				repoNames = append(repoNames, *repo.Name)
			}
		}

		c.JSON(200, gin.H{"repositories": repoNames})
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
