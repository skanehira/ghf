package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/skanehira/clipboard-image/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var uploadCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("repo is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		useClip, err := cmd.PersistentFlags().GetBool("clip")
		if err != nil {
			return
		}

		if !useClip && len(args) == 1 {
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
			fileName = args[1]
			reader, err = os.Open(fileName)
			if err != nil {
				exitError(fmt.Errorf("failed to open %s: %w", fileName, err))
			}
		}

		contents, err := ioutil.ReadAll(reader)
		if err != nil {
			exitError(fmt.Errorf("failed to read contents: %w", err))
		}

		repo := args[0]
		url, err := upload(repo, fileName, contents)
		if err != nil {
			exitError(fmt.Errorf("failed to upload: %w", err))
		}
		fmt.Println(url)
	},
}

func upload(repo, fileName string, contents []byte) (string, error) {
	token := viper.GetString("token")
	user := viper.GetString("user")
	email := viper.GetString("email")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("upload file " + fileName),
		Content: contents,
		Committer: &github.CommitAuthor{
			Name:  github.String(user),
			Email: github.String(email),
		},
	}
	resp, _, err := client.Repositories.CreateFile(ctx, user, repo, fileName, opts)
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
  ghf up {repo} {file} [flags]

Examples:
  $ ghf up {repo} {file}
  $ ghf up {repo} --clip

Args:
  repo     repository
  file     file

Flags:
      --clip   upload from clipboard
  -h, --help   help for up

`)
		return nil
	})
	rootCmd.AddCommand(uploadCmd)
}
