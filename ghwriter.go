package ghwriter

// utility code for writing to files in Github via their API, used for
// automation that writes updated config files, etc into GH.
//

import (
	"context"
	"fmt"
	"io"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

type Writer struct {
	authtoken string
	author    string
	email     string
	branch    string
}

func NewWriter() *Writer {
	return &Writer{
		author: "ghwriter",
		branch: "main",
	}
}

func (w *Writer) ReadConfigFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Can't open GitHub config file %q: %v", filename, err)
	}
	err = yaml.Unmarshal(data, w)
	if err != nil {
		return fmt.Errorf("Can't parse GitHub config file %q: %v", filename, err)
	}
	return nil
}

func (w *Writer) SetAuthToken(authtoken string) *Writer {
	w.authtoken = authtoken
	return w
}

func (w *Writer) SetAuthor(author string) *Writer {
	w.author = author
	return w
}

func (w *Writer) SetEmail(email string) *Writer {
	w.email = email
	return w
}

func (w *Writer) SetBranch(branch string) *Writer {
	w.branch = branch
	return w
}

func (w *Writer) SetOrganization(organization string) *Writer {
	w.organization=organization
	return w
}

func (w *Writer) SetRepo(repo string) *Writer {
	w.repo=repo
	return w
}

func (w *Writer) commitAuthor() *github.CommitAuthor {
	return &github.CommitAuthor{
		Name:  w.author,
		Email: w.email,
	}
}

// WriteFile writes to a file in Github.
func (w *Writer) WriteFile(ctx context.Context, filename, commitMessge string, content []byte) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: w.authtoken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &hithub.RepositoryContentFileOptions{
		Message:   commitMessage,
		Content:   content,
		Branch:    w.branch,
		Committer: w.commitAuthor(),
	}

	_, _, err := client.Repositories.CreateFile(ctx, w.organization, w.repo, filename, opts)
	if err != nil {
		return fmt.Errorf("Unable to write to repo: %v", err)
	}
	return nil
}
