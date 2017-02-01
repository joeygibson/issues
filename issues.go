package main

import (
	"flag"
	"fmt"
	"os"
	"net/url"
	"strings"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"strconv"
)

var (
	numberOfIssues = flag.Int("issues", -1, "maximum number of issues to retrieve")
	repo = flag.String("repo", "", "repository to fetch issues from")
)

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	if repo == nil {
		fmt.Errorf("No repo specified\n")
		flag.Usage()
		os.Exit(1)
	}

	u, err := url.Parse(*repo)

	if err != nil {
		fmt.Errorf("Invalid repo URL", err)
		os.Exit(2)
	}

	path := strings.Split(u.EscapedPath(), "/")
	if len(path) < 3 {
		fmt.Errorf("Invalid Github URL")
		os.Exit(2)
	}

	client := github.NewClient(nil)

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var issues []*github.Issue

	// if numberOfIssues < 0, we fetch all issues
	for len(issues) <= *numberOfIssues || *numberOfIssues < 0 {
		newIssues, resp, err := client.Issues.ListByRepo(path[1], path[2], opt)
		if err != nil {
			fmt.Errorf("error fetching issues for %v/%v", path[1], path[2], err)
			os.Exit(3)
		}

		issues = append(issues, newIssues...)

		if resp.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = resp.NextPage
	}

	if len(issues) == 0 {
		fmt.Println("No issues found")
		os.Exit(0)
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"index", "number", "created_at", "title"})

	for index, issue := range issues {
		if index == *numberOfIssues {
			break
		}

		row := make([]string, 4)
		row[0] = strconv.Itoa(index)
		row[1] = strconv.Itoa(*issue.Number)
		row[2] = issue.CreatedAt.String()
		row[3] = *issue.Title

		table.Append(row)
	}

	table.SetRowLine(true)
	table.Render()
}
