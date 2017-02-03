package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "issues",
		Short: "Shows issues from a Github repo",
		Long:  "Provile a Github personal access token to access private repos.",
		Run:   CmdRoot,
	}
)

func CmdRoot(cmd *cobra.Command, _ []string) {
	path := getRepoPath(cmd)

	apiKey := viper.GetString("api.key")
	ts := loginToGithub(apiKey)

	client := github.NewClient(ts)

	numberOfIssues := viper.GetInt("number.of.issues")

	issues := getIssues(client, path, numberOfIssues)

	if len(issues) == 0 {
		fmt.Println("No issues found")
		os.Exit(0)
	}

	renderTable(issues, numberOfIssues)
}

func getIssues(client *github.Client, path []string, numberOfIssues int) []*github.Issue {
	var issues []*github.Issue

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

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

	return issues
}

func getRepoPath(cmd *cobra.Command) []string {
	repo, _ := cmd.Root().Flags().GetString("repo")

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

	return path
}

func renderTable(issues []*github.Issue, numberOfIssues int) {
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

func loginToGithub(apiKey string) *http.Client {
	var client *http.Client

	if apiKey != "" {
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: apiKey},
		)

		client = oauth2.NewClient(oauth2.NoContext, tokenSource)
	}

	return client
}

func setupCobraAndViper() {
	usr, _ := user.Current()

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(usr.HomeDir, ".issues"))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		// this is OK
	}

	viper.BindEnv("api.key", "ISSUES_API_KEY")

	rootCmd.Flags().StringP("key", "k", "", "Github API key")
	rootCmd.Flags().StringP("repo", "r", "", "Github repo to read")
	rootCmd.Flags().IntP("number-of-issues", "n", -1, "Number of issues to fetch")

	viper.BindPFlag("api.key", rootCmd.Flags().Lookup("key"))
	viper.BindPFlag("repo", rootCmd.Flags().Lookup("repo"))
	viper.BindPFlag("number.of.issues", rootCmd.Flags().Lookup("number-of-issues"))
}

func main() {
	setupCobraAndViper()
	rootCmd.Execute()
}
