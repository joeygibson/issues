# Issues
This is a command-line tool that started out as the code in [this blog post](https://blog.bartfokker.nl/issue-table/). 
I've modified it to use the [Cobra](https://github.com/spf13/cobra) option parser, and the [Viper](https://github.com/spf13/viper)
companion library to support other ways of passing in some options.

# Running
The easiest way to run it is to just specify the URL of the repo, such as 
```bash
issues github.com/joeygibson/issues
```


## Options
```bash
Provile a Github personal access token to access private repos.

Usage:
  issues [flags]

Flags:
  -k, --key string             Github API key
  -n, --number-of-issues int   Number of issues to fetch (default -1)
```
  
# Accessing Private Repos
In order to get to private repos, you need to provide a 
[Github personal token](https://github.com/settings/tokens). There are three ways to provide
this key.

## Config file
Create a file called `config.yml` in `~/.issues` that looks like this
```yml
api.key: your-access-token
```
You can also have a `configl.yml` file in the current directory.

## Environment Variable
You can also set `ISSUES_API_KEY=your-access-token`. This is handy to temporarily
use a different key than what's in your config file.

## Command-line option
Finally, you can pass `-k your-access-token` on the command line.

