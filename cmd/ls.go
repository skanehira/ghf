package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
		owner := args[0]
		repo := args[1]
		branch := args[2]
		useFuzzy, err := cmd.PersistentFlags().GetBool("f")
		if err != nil {
			exitError(fmt.Errorf("failed to get flags: %w", err))
		}
		files := ls(owner, repo, branch, useFuzzy)

		builder := &strings.Builder{}

		for _, f := range files {
			builder.WriteString(fmt.Sprintln(*f.Path))
		}

		fmt.Println(strings.TrimRight(builder.String(), "\r\n"))
	},
}

func ls(owner, repo, branch string, useFuzzy bool) []github.TreeEntry {
	ctx := context.Background()
	client := githubClient(ctx)
	tree, _, err := client.Git.GetTree(ctx, owner, repo, branch, true)
	if err != nil {
		exitError(fmt.Errorf("failed to get branch info: %w", err))
	}

	var entries []github.TreeEntry
	for _, e := range tree.Entries {
		if *e.Type == "tree" {
			continue
		}
		entries = append(entries, e)
	}

	var files []github.TreeEntry
	if useFuzzy {
		idx, err := fuzzyfinder.FindMulti(
			entries,
			func(i int) string {
				return *entries[i].Path
			},
		)
		if err != nil {
			exitError(err)
		}

		for _, i := range idx {
			files = append(files, entries[i])
		}
	} else {
		for _, e := range entries {
			files = append(files, e)
		}
	}
	return files
}

func init() {
	lsCmd.PersistentFlags().Bool("f", false, "fuzyy selector")
	lsCmd.SetUsageFunc(func(*cobra.Command) error {
		fmt.Print(`
Usage:
  ghf ls {owner} {repo} {branch} [--f]

Examples:
  $ ghf ls skanehira ghf master
  $ ghf ls skanehira ghf master --f

Args:
  owner    owner
  repo     repository
  branch   branch

Flags:
      --f      fuzyy selector
  -h, --help   help for ls
`)
		return nil
	})
	rootCmd.AddCommand(lsCmd)
}
