package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	numberOfIssues int
	repo           string
	apikey         string
	rootCmd        = &cobra.Command{
		Use:   "issues",
		Short: "Shows issues from a Github repo",
		Long:  "Currently only works with public repos",
		Run:   CmdRoot,
	}
)

func CmdRoot(cmd *cobra.Command, _ []string) {
	if repo == "" {
		fmt.Printf("%v\n", "No repo specified")
		cmd.Help()
		os.Exit(1)
	}

	u, err := url.Parse(repo)

	if err != nil {
		fmt.Printf("Invalid repo URL: %v\n", err)
		os.Exit(2)
	}

	path := strings.Split(u.EscapedPath(), "/")

	if len(path) < 3 {
		fmt.Printf("Invalid Github URL: %v\n", u)
		os.Exit(2)
	}

	ts := LoginToGithub(apikey)

	client := github.NewClient(ts)

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var issues []*github.Issue

	// if numberOfIssues < 0, we fetch all issues
	for len(issues) <= numberOfIssues || numberOfIssues < 0 {
		newIssues, resp, err := client.Issues.ListByRepo(path[1], path[2], opt)
		if err != nil {
			fmt.Printf("error fetching issues for %v/%v: %v\n", path[1], path[2], err)
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
		if index == numberOfIssues {
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

func LoginToGithub(apiKey string) *http.Client {
	var client *http.Client

	if apiKey != "" {
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: apikey},
		)

		client = oauth2.NewClient(oauth2.NoContext, tokenSource)
	}

	return client
}

func main() {
	rootCmd.Flags().IntVarP(&numberOfIssues, "number-of-issues",
		"n", -1, "the number of issues to retrieve")
	rootCmd.Flags().StringVarP(&repo, "repo", "r", "", "the Github URL to retrieve")
	rootCmd.Flags().StringVarP(&apikey, "api-key", "a", "", "a Github personal API key")
	rootCmd.Execute()
}
