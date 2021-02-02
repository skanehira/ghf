package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/go-github/github"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "del file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("owner, repo and branch name is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		owner := args[0]
		repo := args[1]
		branch := args[2]

		del(owner, repo, branch, ls(owner, repo, branch, true))
	},
}

func del(owner, repo string, branch string, files []github.TreeEntry) {
	email := viper.GetString("email")

	ctx := context.Background()
	client := githubClient(ctx)

	for _, file := range files {
		opts := &github.RepositoryContentFileOptions{
			Message: github.String("delete file " + *file.Path),
			SHA:     file.SHA,
			Committer: &github.CommitAuthor{
				Name:  github.String(owner),
				Email: github.String(email),
			},
			Branch: &branch,
		}
		_, _, err := client.Repositories.DeleteFile(ctx, owner, repo, *file.Path, opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		fmt.Println("deleted", *file.Path)
	}
}

func init() {
	delCmd.SetUsageFunc(func(*cobra.Command) error {
		fmt.Print(`
Usage:
  ghf del {owner} {repo} {branch}

Examples:
  $ ghf del skanehira ghf master

Args:
  owner    owner
  repo     repository
  branch   branch

Flags:
  -h, --help   help for ls
`)
		return nil
	})
	rootCmd.AddCommand(delCmd)
}
