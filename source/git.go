package source

import (
	"context"
	"io"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// NewGit returns a new Source
func NewGit(ctx context.Context, tree *object.Tree, filters ...PathFilter) Source {
	src := Git{ch: make(chan gitItem)}
	go src.walk(ctx, tree, filters)
	return src
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
func (g Git) Next() (Item, error) {
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
func NewRepo(address string) (*Repo, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: address})
	if err != nil {
		return nil, err
	}
	return &Repo{repo: repo}, nil
}

// Repo is a git repository
type Repo struct {
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

// Tree returns the tree for the given reference (example: refs/remotes/origin/master)
func (r *Repo) Tree(reference string) (*object.Tree, error) {
	hash, err := r.repo.Reference(plumbing.ReferenceName(reference), false)
	if err != nil {
		return nil, err
	}
	commit, err := r.repo.CommitObject(hash.Hash())
	if err != nil {
		return nil, err
	}
	return commit.Tree()
}
