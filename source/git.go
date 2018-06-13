package source

import (
	"context"
	"io"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/tent.v1/item"
)

// NewGit returns a new Source
func NewGit(ctx context.Context, commit *object.Commit, filters ...PathFilter) (Git, error) {
	tree, err := commit.Tree()
	if err != nil {
		return Git{}, err
	}
	src := Git{ch: make(chan gitItem)}
	go src.walk(ctx, tree, filters)
	return src, nil
}

// Git takes a git repo as origin
type Git struct {
	ch chan gitItem
}

func (g Git) walk(ctx context.Context, tree *object.Tree, filters []PathFilter) {
	done := ctx.Done()
	tree.Files().ForEach(func(file *object.File) error {
		select {
		case <-done:
			return nil
		default:
			for _, f := range filters {
				if !f(file.Name) {
					return nil
				}
			}
			g.ch <- gitItem{name: file.Name, blob: file.Blob}
			return nil
		}
	})
	close(g.ch)
}

// Next implements the Source interface
func (g Git) Next() (item.Item, error) {
	item, ok := <-g.ch
	if !ok {
		return nil, nil
	}
	return item, nil
}

// Item represents a Item in a git repo
type gitItem struct {
	name string
	blob object.Blob
}

// Name implements the Item interface
func (i gitItem) Name() string {
	return i.name
}

// Content implements the Item interface
func (i gitItem) Content() (io.ReadCloser, error) {
	return i.blob.Reader()
}

// NewRepo returns a new repo
func NewRepo(url string) (*Repo, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: url})
	if err != nil {
		return nil, err
	}
	return &Repo{URL: url, repo: repo}, nil
}

// Repo is a git repository
type Repo struct {
	URL  string
	repo *git.Repository
}

// Update pulls new contents
func (r *Repo) Update() error {
	err := r.repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

// Commit returns the Commit for the given reference (example: refs/remotes/origin/master)
func (r *Repo) Commit(reference string) (*object.Commit, error) {
	hash, err := r.repo.Reference(plumbing.ReferenceName(reference), false)
	if err != nil {
		return nil, err
	}
	return r.repo.CommitObject(hash.Hash())
}
