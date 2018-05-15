package source

import (
	"context"
	"strings"
	"testing"
)

func TestRepo(t *testing.T) {
	r, err := NewRepo("https://github.com/go-tent/tent")
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Update(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Tree("refs/remotes/origin/v1"); err != nil {
		t.Fatal(err)
	}
}

func TestGit(t *testing.T) {
	r, err := NewRepo("https://github.com/go-tent/tent")
	if err != nil {
		t.Fatal(err)
	}
	tree, err := r.Tree("refs/remotes/origin/v1")
	if err != nil {
		t.Fatal(err)
	}
	baseTest(t, NewGit(context.Background(), tree, func(s string) bool {
		return strings.HasPrefix(s, "source/common")
	}), 2)
}
