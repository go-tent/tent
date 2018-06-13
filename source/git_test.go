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
	if _, err := r.Commit("refs/remotes/origin/v1"); err != nil {
		t.Fatal(err)
	}
}

func TestGit(t *testing.T) {
	r, err := NewRepo("https://github.com/go-tent/tent")
	if err != nil {
		t.Fatal(err)
	}
	commit, err := r.Commit("refs/remotes/origin/v1")
	if err != nil {
		t.Fatal(err)
	}
	git, err := NewGit(context.Background(), commit, func(s string) bool {
		return strings.HasPrefix(s, "source/common")
	})
	if err != nil {
		t.Fatal(err)
	}
	baseTest(t, git, 2)
}
