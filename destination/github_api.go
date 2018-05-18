package destination

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/tent.v1/item"
)

// RepoCfg specifies the Repo used by GithubAPI
type RepoCfg struct {
	Owner  string
	Repo   string
	Branch string
}

// NewGihubAPI returns a new GithubAPI destination.
func NewGihubAPI(ctx context.Context, client *github.Client, config RepoCfg) *GithubAPI {
	return &GithubAPI{
		ctx:    ctx,
		client: client.Repositories,
		config: config,
	}
}

// GithubAPI stores files using github API.
type GithubAPI struct {
	ctx    context.Context
	client *github.RepositoriesService
	config RepoCfg
}

// Hash returns file hash using Github API.
func (g *GithubAPI) Hash(ctx context.Context, i item.Item) (string, error) {
	var o = github.RepositoryContentGetOptions{Ref: g.config.Branch}
	contents, _, _, err := g.client.GetContents(ctx, g.config.Owner, g.config.Repo, i.Name(), &o)
	if err != nil {
		return "", err
	}
	return contents.GetSHA(), nil
}

// Create adds a new Item to the Destination.
func (g *GithubAPI) Create(ctx context.Context, i item.Item) error {
	return g.request(ctx, "create", i, "")
}

// Update overwrites an Item in the Destination.
func (g *GithubAPI) Update(ctx context.Context, i item.Item, hash string) error {
	return g.request(ctx, "update", i, hash)
}

// Delete removes an Item from the Destination.
func (g *GithubAPI) Delete(ctx context.Context, i item.Item, hash string) error {
	return g.request(ctx, "delete", i, hash)
}

func (g *GithubAPI) request(ctx context.Context, action string, i item.Item, hash string) error {
	var o = github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("%s %s", strings.Title(action), path.Base(i.Name()))),
		Branch:  github.String(g.config.Branch),
	}
	if action != "delete" {
		r, err := i.Content()
		if err != nil {
			return err
		}
		defer r.Close()
		if o.Content, err = ioutil.ReadAll(r); err != nil {
			return err
		}
	}
	if action != "create" {
		o.SHA = github.String(hash)
	}
	var err error
	switch action {
	case "create":
		_, _, err = g.client.CreateFile(ctx, g.config.Owner, g.config.Repo, i.Name(), &o)
	case "delete":
		_, _, err = g.client.DeleteFile(ctx, g.config.Owner, g.config.Repo, i.Name(), &o)
	case "update":
		_, _, err = g.client.UpdateFile(ctx, g.config.Owner, g.config.Repo, i.Name(), &o)

	}
	return err
}
