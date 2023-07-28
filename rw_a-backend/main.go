package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

var (
	accessToken = os.Getenv("GH_API_TOKEN")
)

type GitHubIssue struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	State    string `json:"state"`
	IssueUrl string `json:"html_url"`
}

type UrlParts struct {
	Scheme   string
	Host     string
	Path     string
	RawQuery string
}

type GhRepoUri struct {
	Owner     string
	Repo      string
	ExtraInfo []string
}

func parseUrl(repoUrl string) (*UrlParts, error) {
	u, err := url.Parse(repoUrl)
	if err != nil {
		return nil, err
	}
	urlParts := &UrlParts{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     u.Path,
		RawQuery: u.RawQuery,
	}
	return urlParts, nil
}

func splitPath(path string) (*GhRepoUri, error) {
	parts := strings.Split(path, "/")
	parts = parts[1:]

	if len(parts) < 2 {
		return nil, nil
	}

	gitHubRepo := &GhRepoUri{
		Owner: parts[0],
		Repo:  parts[1],
	}

	if len(parts) > 2 {
		gitHubRepo.ExtraInfo = parts[2:]
	}

	return gitHubRepo, nil
}

func listOpenIssues(c *gin.Context) {
	repoUrl := c.PostForm("repo")
	// parse repoUrl and split it
	parsedUrl, _ := parseUrl(repoUrl)
	splitUrl, _ := splitPath(parsedUrl.Path)

	owner := splitUrl.Owner
	repo := splitUrl.Repo

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.IssueListByRepoOptions{
		State: "all",
	}
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": err.Error()})
		return
	}

	var openIssues []GitHubIssue
	for _, issue := range issues {
		openIssues = append(openIssues, GitHubIssue{
			Number:   issue.GetNumber(),
			Title:    issue.GetTitle(),
			State:    issue.GetState(),
			IssueUrl: issue.GetHTMLURL(),
		})
	}

	c.HTML(http.StatusOK, "issues.html", gin.H{
		"Owner":    owner,
		"Repo":     repo,
		"Issues":   openIssues,
		"Number":   openIssues,
		"IssueUrl": openIssues,
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "RW GH Issues Renderer",
			},
		)
	})
	router.POST("/issues", listOpenIssues)

	port := "8080"

	router.Run(":" + port)
}
