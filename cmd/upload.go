package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
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
		if len(args) < 3 {
			return errors.New("owner, repo and branch is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		useClip, err := cmd.PersistentFlags().GetBool("clip")
		if err != nil {
			return
		}
		dir, err := cmd.PersistentFlags().GetString("dir")
		if err != nil {
			return
		}

		if !useClip && len(args) == 3 {
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
			path := path.Join(dir, time.Now().Format("20060102150405.999999999")+".png")
			files[path] = buf.Bytes()
		} else {
			for _, fileName := range args[3:] {
				b, err := ioutil.ReadFile(fileName)
				if err != nil {
					exitError(fmt.Errorf("failed to read file %s: %w", fileName, err))
				}

				path := filepath.Join(dir, filepath.Base(fileName))
				files[path] = b
			}
		}

		owner := args[0]
		repo := args[1]
		branch := args[2]
		upload(owner, repo, branch, files)
	},
}

func upload(owner, repo string, branch string, files map[string][]byte) {
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
			Branch: &branch,
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
		fmt.Println(*repo.Content.HTMLURL + "?raw=true")
	}
}

func init() {
	uploadCmd.PersistentFlags().Bool("clip", false, "upload from clipboard")
	uploadCmd.PersistentFlags().String("dir", time.Now().Format("20060102150405"), "file directory")
	uploadCmd.SetUsageFunc(func(*cobra.Command) error {
		fmt.Print(`
Usage:
  ghf up {owner} {repo} {branch} [file...] [flags]

Examples:
  $ ghf up skanehira images main sample1.png sample2.png
  $ ghf up skanehira images main --clip
  $ ghf up skanehira images main sample.png --dir gorilla

Flags:
      --dir    file directory
      --clip   upload from clipboard
  -h, --help   help for up
`)
		return nil
	})
	rootCmd.AddCommand(uploadCmd)
}
