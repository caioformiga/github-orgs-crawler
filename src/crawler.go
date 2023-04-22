package src

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type CrawlerCtrl struct {
	home       string
	userPage   string
	reposPage  string
	pattern    string
	pageNumber int
	collector  *colly.Collector
	repos      []string
	ignoreMap  map[string]string
}

func NewCrawlerCtrl(gitHome, user, repo string) *CrawlerCtrl {
	ctrl := &CrawlerCtrl{
		home:       gitHome,
		userPage:   user,
		reposPage:  repo,
		pageNumber: 1,
		collector:  colly.NewCollector(),
		repos:      []string{},
		ignoreMap:  map[string]string{},
	}

	ctrl.pattern = "a.d-inline-block"
	ctrl.loadIgnoreMap()
	return ctrl
}

func (ctrl *CrawlerCtrl) loadIgnoreMap() {
	if len(ctrl.ignoreMap) == 0 {
		ctrl.ignoreMap["/klever-io/.github"] = ".github"
	}
}

func (ctrl *CrawlerCtrl) ignoreItem(item string) bool {
	_, ok := ctrl.ignoreMap[item]
	return ok
}

// auth should get authorized access to github repositories by token
func (ctrl *CrawlerCtrl) auth() {
}

func (ctrl *CrawlerCtrl) FetchRepos() []string {
	fmt.Printf("fetching...")

	// show progress on repositories page
	ctrl.collector.OnRequest(func(r *colly.Request) {
		fmt.Printf(".")
	})

	// get all patterns on html page and add on slice
	ctrl.collector.OnHTML(ctrl.pattern, func(e *colly.HTMLElement) {
		// get only link itens
		item := e.Attr("href")

		// validate if item is signup link
		if strings.Contains(item, "signup") {
			return
		}

		// validate if item should be ignored
		if ctrl.ignoreItem(item) {
			return
		}

		// add on slice item with url prefix
		repoURL := fmt.Sprintf("%s%s", ctrl.home, item)
		ctrl.repos = append(ctrl.repos, repoURL)

	})

	ctrl.pageNumber = 1
	qty := 0
	for {
		// qty before run crawler actions
		qty = len(ctrl.repos)

		// print detaled info
		if os.Getenv("DEBUG") == "true" {
			fmt.Printf("URL: %s Page %d\n", ctrl.reposPage, ctrl.pageNumber)
		}

		// update ctrl.repos if has repos on that page
		ctrl.collector.Visit(fmt.Sprintf(ctrl.reposPage, ctrl.userPage, ctrl.pageNumber))

		// no more pages to visit
		if qty == len(ctrl.repos) {
			fmt.Println(".")
			break
		}

		time.Sleep(time.Millisecond * 100)
		ctrl.pageNumber++
	}

	return ctrl.repos
}
