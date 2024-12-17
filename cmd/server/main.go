package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

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

	// Create repo
	router.POST("/repos", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&req); err != nil || req.Name == "" {
			c.JSON(400, gin.H{"error": "invalid request"})
		}

		repo, err := ghClient.CreateRepo(context.Background(), req.Name)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "Repository created", "name": *repo.Name})
	})

	// Delete repo
	router.DELETE("/repos/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(400, gin.H{"error": "repository name is required"})
			return
		}

		err := ghClient.DeleteRepo(context.Background(), name)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Repository deleted", "repo": name})
	})

	// List all repos
	router.GET("/repos", func(c *gin.Context) {
		repos, err := ghClient.ListRepos(context.Background())
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

	// List N open pull requests for a repo
	router.GET("/repos/:name/pulls", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(400, gin.H{"error": "repository name is required"})
			return
		}

		// Default is -1 means no limit
		nStr := c.DefaultQuery("n", "-1")
		n, err := strconv.Atoi(nStr)
		if err != nil || n < -1 {
			c.JSON(400, gin.H{"error": "invalid value for n"})
			return
		}

		prs, err := ghClient.ListPullRequests(context.Background(), name, n)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"reposiroty": name, "pull_requests": prs, "count": len(prs)})
	})

	// Start server
	port := ":8080"
	fmt.Printf("Server running on port http://localhost%s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
