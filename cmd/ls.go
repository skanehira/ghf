package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/ktr0731/go-fuzzyfinder"
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
		}

		useFuzzy, err := cmd.PersistentFlags().GetBool("f")
		if err != nil {
			exitError(fmt.Errorf("failed to get flags: %w", err))
		}

		var files []github.TreeEntry
		for _, e := range tree.Entries {
			if *e.Type == "tree" {
				continue
			}
			files = append(files, e)
		}

		if useFuzzy {
			idx, err := fuzzyfinder.FindMulti(
				files,
				func(i int) string {
					return *files[i].Path
				},
			)
			if err != nil {
				exitError(err)
			}

			for _, i := range idx {
				fmt.Println(*files[i].URL)
			}
		} else {
			for _, e := range files {
				fmt.Println(*e.Path)
			}
		}
	},
}

func init() {
	lsCmd.PersistentFlags().Bool("f", false, "fuzyy selector")
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
