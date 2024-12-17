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

// func TestCreateRepo(t *testing.T) {
// 	mockClient := &mocks.MockGitHubClient{}
// 	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
// 	router := setupRouter(ghClient)

// 	reqBody := `{"name":"test-repo"}`
// 	req, _ := http.NewRequest("POST", "/repos", strings.NewReader(reqBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	if w.Code != http.StatusCreated {
// 		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
// 	}

// 	expected := `{"message":"Repository created","name":"test-repo"}`
// 	if strings.TrimSpace(w.Body.String()) != expected {
// 		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
// 	}
// }

// func TestDeleteRepo(t *testing.T) {
// 	mockClient := &mocks.MockGitHubClient{
// 		Repos: []*github.Repository{
// 			{Name: github.String("test-repo")},
// 		},
// 	}
// 	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
// 	router := setupRouter(ghClient)

// 	req, _ := http.NewRequest("DELETE", "/repos/test-repo", nil)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	if w.Code != http.StatusOK {
// 		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
// 	}

// 	expected := `{"message":"Repository deleted","repo":"test-repo"}`
// 	if strings.TrimSpace(w.Body.String()) != expected {
// 		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
// 	}
// }

// func TestListRepos(t *testing.T) {
// 	mockClient := &mocks.MockGitHubClient{
// 		Repos: []*github.Repository{
// 			{Name: github.String("repo1")},
// 			{Name: github.String("repo2")},
// 		},
// 	}
// 	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
// 	router := setupRouter(ghClient)

// 	req, _ := http.NewRequest("GET", "/repos", nil)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	if w.Code != http.StatusOK {
// 		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
// 	}

// 	expected := `{"repositories":["repo1","repo2"]}`
// 	if strings.TrimSpace(w.Body.String()) != expected {
// 		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
// 	}
// }

// func TestListPullRequests(t *testing.T) {
// 	mockClient := &mocks.MockGitHubClient{
// 		PullRequests: []*github.PullRequest{
// 			{Number: github.Int(1)},
// 			{Number: github.Int(2)},
// 		},
// 	}
// 	ghClient := githubapi.NewTestClient(mockClient, "test-owner")
// 	router := setupRouter(ghClient)

// 	req, _ := http.NewRequest("GET", "/repos/test-repo/pulls?n=1", nil)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	if w.Code != http.StatusOK {
// 		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
// 	}
// }
