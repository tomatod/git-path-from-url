package main

import (
	"errors"
	"fmt"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// ------------------------
// Util functions
// ------------------------
func getGitRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stderr = nil
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("current directory is not inside a git repository")
	}
	return strings.TrimSpace(string(output)), nil
}

// ------------------------
// URL Convertors
// ------------------------
type urlConvertor interface {
	GetLocalPath() (string, error)
}

type gitHubHttpsUrlConvertor struct {
	rawUrl  string
	host    string
	owner   string
	repo    string
	subPath string
	baseDir string
}

func newGitHubHttpsUrlConvertor(raw string) (*gitHubHttpsUrlConvertor, error) {
	u, err := neturl.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	host := u.Host
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}

	segments := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(segments) < 2 {
		return nil, errors.New("invalid github url: path must contain owner and repo")
	}
	owner := segments[0]
	repo := strings.TrimSuffix(segments[1], ".git")

	subPath := ""
	if len(segments) > 2 {
		if segments[2] == "blob" || segments[2] == "tree" {
			if len(segments) > 4 {
				subPath = filepath.Join(segments[4:]...)
			}
		} else {
			subPath = filepath.Join(segments[2:]...)
		}
	}

	baseDir, err := getGitRepoRoot()
	if err != nil {
		return nil, err
	}

	return &gitHubHttpsUrlConvertor{
		rawUrl:  raw,
		host:    host,
		owner:   owner,
		repo:    repo,
		subPath: subPath,
		baseDir: baseDir,
	}, nil
}

func (g *gitHubHttpsUrlConvertor) GetLocalPath() (string, error) {
	if g.owner == "" || g.repo == "" {
		return "", errors.New("invalid convertor state")
	}

	path := g.baseDir
	if g.subPath != "" {
		return filepath.Join(path, g.subPath), nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", path)
	} else if err != nil {
		return "", fmt.Errorf("failed to access path: %w", err)
	}

	return path, nil
}

// ------------------------
// Primary proccesses
// ------------------------
func getConvertorFromUrl(raw string) (urlConvertor, error) {
	u, err := neturl.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	host := u.Host
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}

	if strings.EqualFold(host, "github.com") && (u.Scheme == "https" || u.Scheme == "http") {
		return newGitHubHttpsUrlConvertor(raw)
	}

	return nil, fmt.Errorf("unsupported url: %s", raw)
}

func format(convertor urlConvertor) (string, error) {
	return convertor.GetLocalPath()
}

func process(raw string) (string, error) {
	conv, err := getConvertorFromUrl(raw)
	if err != nil {
		return "", err
	}

	return format(conv)
}

// ------------------------
// CLI process
// ------------------------
func main() {
	app := &cli.App{
		Name:  "git-path-from-url",
		Usage: "Convert URL to local path (base: current git repo root).",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("usage: git path-from-url <url>")
			}
			raw := c.Args().Get(0)

			result, err := process(raw)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
