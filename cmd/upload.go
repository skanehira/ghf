package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
	"github.com/skanehira/clipboard-image/v2"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var uploadCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file",
	RunE: func(cmd *cobra.Command, args []string) error {
		useClip, err := cmd.PersistentFlags().GetBool("clip")
		if err != nil {
			return err
		}

		argsLen := len(args)

		if argsLen < 1 {
			return cmd.Usage()
		}

		//
		if useClip && len(args) < 2 {

		}

		// 引数が2つある場合は repo と filename として処理
		if len(args) > 1 {

			return nil
		}

		var (
			fileName string
			r        io.Reader
		)
		if useClip {
			fileName = time.Now().Format("20060102150405") + ".png"
			r, err = clipboard.Read()
		} else {
			fileName = filepath.Base(args[1])
			r, err = os.Open(args[1])
			if err != nil {
				return err
			}
		}

		contents, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		repo := ""
		url, err := upload(repo, fileName, contents)
		if err != nil {
			return err
		}
		fmt.Println(url)

		return nil
	},
}

func upload(repo, fileName string, contents []byte) (string, error) {
	// TODO get info from config
	token := ""
	user := ""
	email := ""

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
	rootCmd.AddCommand(uploadCmd)
}
