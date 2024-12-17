package integration

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jorgebaptista/octo-manager/internal/githubapi"
)

func SetupRouter(ghClient *githubapi.Client) *gin.Engine {
	router := gin.Default()

	router.POST("/repos", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&req); err != nil || req.Name == "" {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		repo, err := ghClient.CreateRepo(c.Request.Context(), req.Name)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "Repository created", "name": *repo.Name})
	})

	router.DELETE("/repos/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(400, gin.H{"error": "repository name is required"})
			return
		}

		err := ghClient.DeleteRepo(c.Request.Context(), name)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Repository deleted", "repo": name})
	})

	router.GET("/repos", func(c *gin.Context) {
		repos, err := ghClient.ListRepos(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var repoNames []string
		for _, repo := range repos {
			if repo.Name != nil {
				repoNames = append(repoNames, *repo.Name)
			}
		}

		c.JSON(200, gin.H{"repositories": repoNames})
	})

	router.GET("/repos/:name/pulls", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(400, gin.H{"error": "repository name is required"})
			return
		}

		nStr := c.DefaultQuery("n", "-1")
		n, err := strconv.Atoi(nStr)
		if err != nil || n < -1 {
			c.JSON(400, gin.H{"error": "invalid value for n"})
			return
		}

		prs, err := ghClient.ListPullRequests(c.Request.Context(), name, n)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"repository": name, "pull_requests": prs, "count": len(prs)})
	})

	return router
}
