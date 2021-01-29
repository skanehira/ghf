package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("owner, repo and branch name is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		user := args[0]
		repo := args[1]
		branch := args[2]

		ctx := context.Background()
		client := githubClient(ctx)
		tree, _, err := client.Git.GetTree(ctx, user, repo, branch, true)
		if err != nil {
			exitError(fmt.Errorf("failed to get branch info: %w", err))
			return
		}

		files := make([]string, len(tree.Entries))
		for i, e := range tree.Entries {
			files[i] = *e.Path
			fmt.Println(*e.Path)
		}
	},
}

func init() {
	lsCmd.SetUsageFunc(func(*cobra.Command) error {
		fmt.Print(`
Usage:
  ghf ls {owner} {repo} {branch}

Examples:
  $ ghf ls skanehira ghf master

Args:
  repo     repository

Flags:
  -h, --help   help for ls
`)
		return nil
	})
	rootCmd.AddCommand(lsCmd)
}
