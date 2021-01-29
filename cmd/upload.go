package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"

	"github.com/skanehira/clipboard-image/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("owner and repo is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		useClip, err := cmd.PersistentFlags().GetBool("clip")
		if err != nil {
			return
		}

		if !useClip && len(args) == 2 {
			printError("file is required")
			cmd.Usage()
			return
		}

		var (
			fileName string
			reader   io.Reader
		)

		if useClip {
			fileName = time.Now().Format("20060102150405") + ".png"
			reader, err = clipboard.Read()
			if err != nil {
				exitError(fmt.Errorf("failed to get contents from clipboard: %w", err))
			}
		} else {
			fileName = filepath.Base(args[2])
			reader, err = os.Open(args[2])
			if err != nil {
				exitError(fmt.Errorf("failed to open %s: %w", fileName, err))
			}
		}

		contents, err := ioutil.ReadAll(reader)
		if err != nil {
			exitError(fmt.Errorf("failed to read contents: %w", err))
		}

		owner := args[0]
		repo := args[1]
		url, err := upload(owner, repo, fileName, contents)
		if err != nil {
			exitError(fmt.Errorf("failed to upload: %w", err))
		}
		fmt.Println(url)
	},
}

func upload(owner, repo, fileName string, contents []byte) (string, error) {
	email := viper.GetString("email")

	ctx := context.Background()
	client := githubClient(ctx)
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("upload file " + fileName),
		Content: contents,
		Committer: &github.CommitAuthor{
			Name:  github.String(owner),
			Email: github.String(email),
		},
	}
	resp, _, err := client.Repositories.CreateFile(ctx, owner, repo, fileName, opts)
	if err != nil {
		return "", err
	}
	return *resp.Content.DownloadURL, nil
}

func init() {
	uploadCmd.PersistentFlags().Bool("clip", false, "upload from clipboard")
	uploadCmd.SetUsageFunc(func(*cobra.Command) error {
		fmt.Print(`
Usage:
  ghf up {owner} {repo} [file] [flags]

Examples:
  $ ghf up skanehira images sample.png
  $ ghf up skanehira images --clip

Flags:
      --clip   upload from clipboard
  -h, --help   help for up
`)
		return nil
	})
	rootCmd.AddCommand(uploadCmd)
}
