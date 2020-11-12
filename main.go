package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/skanehira/clipboard-image/v2"
	"golang.org/x/oauth2"
)

var configFile = ".github_token"

var version = "0.0.1"

type Git struct {
	User  string
	Email string
}

func getToken() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		configFile := filepath.Join(homeDir, configFile)
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			return "", err
		}

		token = strings.Trim(string(b), "\r\n")
	}
	if token == "" {
		return "", errors.New("github token is empty")
	}

	return token, nil
}

func getGitUser() (Git, error) {
	git := Git{
		User:  "unknown",
		Email: "unknown",
	}

	cmd := exec.Command("git", "config", "user.name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return git, errors.New(string(out))
	} else {
		git.User = strings.TrimRight(string(out), "\r\n")
	}

	cmd = exec.Command("git", "config", "user.email")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return git, err
	} else {
		git.Email = strings.TrimRight(string(out), "\r\n")
	}

	return git, nil
}

func run(args []string) error {
	var (
		fileName string
		r        io.Reader
		err      error
	)

	repo := args[0]

	// upload file from clipboard
	if len(args) == 1 {
		fileName = time.Now().Format("20060102150405") + ".png"
		r, err = clipboard.Read()
		if err != nil {
			return err
		}
	}

	// upload specified file
	if len(args) > 1 {
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

	url, err := upload(repo, fileName, contents)
	if err != nil {
		return err
	}
	fmt.Println(url)

	return nil
}

func upload(repo, fileName string, contents []byte) (string, error) {
	token, err := getToken()
	if err != nil {
		return "", err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	git, err := getGitUser()
	if err != nil {
		return "", err
	}
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("upload file " + fileName),
		Content: contents,
		Committer: &github.CommitAuthor{
			Name:  github.String(git.User),
			Email: github.String(git.Email),
		},
	}
	resp, _, err := client.Repositories.CreateFile(ctx, git.User, repo, fileName, opts)
	if err != nil {
		return "", err
	}
	return *resp.Content.DownloadURL, nil
}

func main() {
	name := "upimg"
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fs.SetOutput(os.Stdout)
		fmt.Printf(`%[1]s - Upload image file to GitHub repository

VERSION: %s

USAGE:
  $ %[1]s repo [file]

EXAMPLE:
  $ %[1]s images sample.png
  $ %[1]s images
`, name, version)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return
		}
		os.Exit(1)
	}

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		return
	}

	if err := run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
