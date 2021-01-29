package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

		files := map[string][]byte{}

		if useClip {
			reader, err := clipboard.Read()
			if err != nil {
				exitError(fmt.Errorf("failed to get contents from clipboard: %w", err))
			}
			buf := &bytes.Buffer{}
			if _, err := io.Copy(buf, reader); err != nil {
				exitError(fmt.Errorf("failed to read contents from clipboard: %w", err))
			}
			fileName := time.Now().Format("20060102150405") + ".png"
			files[fileName] = buf.Bytes()
		} else {
			for _, fileName := range args[2:] {
				b, err := ioutil.ReadFile(fileName)
				if err != nil {
					exitError(fmt.Errorf("failed to read file %s: %w", fileName, err))
				}

				files[filepath.Base(fileName)] = b
			}
		}

		owner := args[0]
		repo := args[1]
		upload(owner, repo, files)
	},
}

func upload(owner, repo string, files map[string][]byte) {
	email := viper.GetString("email")

	ctx := context.Background()
	client := githubClient(ctx)

	for name, contents := range files {
		opts := &github.RepositoryContentFileOptions{
			Message: github.String("upload file " + name),
			Content: contents,
			Committer: &github.CommitAuthor{
				Name:  github.String(owner),
				Email: github.String(email),
			},
		}
		repo, resp, err := client.Repositories.CreateFile(ctx, owner, repo, name, opts)
		if resp.StatusCode == 422 {
			printError("same file is already exists")
			continue
		}
		if err != nil {
			printError(err)
			continue
		}
		fmt.Println(*repo.Content.DownloadURL)
	}
}

func init() {
	uploadCmd.PersistentFlags().Bool("clip", false, "upload from clipboard")
	uploadCmd.SetUsageFunc(func(*cobra.Command) error {
		fmt.Print(`
Usage:
  ghf up {owner} {repo} [file...] [flags]

Examples:
  $ ghf up skanehira images sample1.png sample2.png
  $ ghf up skanehira images --clip

Flags:
      --clip   upload from clipboard
  -h, --help   help for up
`)
		return nil
	})
	rootCmd.AddCommand(uploadCmd)
}
