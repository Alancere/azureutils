package spec

import (
	"context"
	"os"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func CredGithubRepository() *github.Client {
	token := os.Getenv("GITHUB_TOKEN_JH")

	// oauth2 认证
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
