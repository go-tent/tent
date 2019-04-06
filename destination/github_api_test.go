package destination

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"github.com/go-tent/tent/item"
)

func TestGithubAPI(t *testing.T) {
	env := ensureEnv(t, "GITHUB", "TOKEN", "OWNER", "REPO", "BRANCH")
	ctx := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env["TOKEN"]},
	)))
	dest := NewGihubAPI(ctx, client, RepoCfg{Owner: env["OWNER"], Repo: env["REPO"], Branch: env["BRANCH"]})
	item := item.Memory{ID: "test/.category.yml", Contents: []byte("index: 10\ntitle: this is a test")}

	if err := dest.Create(ctx, item); err != nil {
		t.Fatalf("Create: %s", err)
	}

	h, err := dest.Hash(ctx, item)
	if err != nil {
		t.Fatalf("Hash: %s", err)
	}

	item.Contents = append(item.Contents, []byte(", updated value")...)
	if err := dest.Update(ctx, item, h); err != nil {
		t.Fatalf("Update: %s", err)
	}

	h, err = dest.Hash(ctx, item)
	if err != nil {
		t.Fatalf("Hash: %s", err)
	}
	if err := dest.Delete(ctx, item, h); err != nil {
		t.Fatalf("Delete: %s", err)
	}
}

func ensureEnv(t *testing.T, prefix string, keys ...string) map[string]string {
	var m = make(map[string]string, len(keys))
	for _, k := range keys {
		name := fmt.Sprintf("%s_%s", prefix, k)
		if m[k] = os.Getenv(name); m[k] == "" {
			t.Fatalf("$%s not found", name)
		}
	}
	return m
}
