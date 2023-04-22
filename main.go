package main

import (
	"fmt"
	"os"

	"github.com/caioformiga/crwaler/src"
	env "github.com/joho/godotenv"
)

func main() {
	env.Load("default.env", ".env")

	ctrl := src.NewCrawlerCtrl(os.Getenv("GIT_URL"), os.Getenv("GIT_ORG_USER"), os.Getenv("REPO_PAGE"))

	repos := ctrl.FetchRepos()

	for i := range repos {
		fmt.Println(repos[i])
	}
}
