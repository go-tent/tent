package transifex

import (
	"os"
	"testing"
	"time"
)

var (
	_APIKey  = os.Getenv("TX_API_KEY")
	_Org     = os.Getenv("TX_ORG")
	_Project = os.Getenv("TX_PROJ")
)

func TestClient(t *testing.T) {
	const (
		slug    = "test-resource"
		lang    = "es"
		content = `
some:
  key: value
  number: 10
another:
  another: some text
  one: we are happy
  two: he is sad
  deep:
    nested: what's wrong
`
		extra = `
another: this is it.`
		translation = `some:
  key: valor
  number: 10
another:
  another: algun texto
  one: somos felices
  two: es triste
  deep:
    nested: que pasa
another: eso es.
`
	)
	c := NewClient(_APIKey, _Org, _Project)
	c.SetTicker(time.NewTicker(time.Hour / 6000))
	proj, err := c.Project()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proj)
	res, err := c.ListResources()
	if err != nil {
		t.Fatal(err)
	}
	var exists bool
	for _, r := range res {
		if r.Slug != slug {
			continue
		}
		exists = true
		break
	}
	if exists {
		if err := c.DeleteResource(slug); err != nil {
			t.Fatal(err)
		}
	}
	r, err := c.CreateResource(UploadResourceRequest{
		BaseResource:       BaseResource{Slug: slug, Name: slug + `.yml`, I18nType: `YAML_GENERIC`},
		AcceptTranslations: true, Content: content})
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%+v", r)
	r, err = c.UpdateResource(slug, content+extra)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%+v", r)
	r, err = c.UpdateTranslation(slug, lang, translation)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%+v", r)
	tx, err := c.GetTranslation(slug, lang)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%+v", tx)
	file, err := c.GetTranslationFile(slug, lang)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%s", file)
}
