package transifex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func NewClient(apiKey, org, project string) *Client {
	return &Client{
		client: &http.Client{
			Transport: apiAuth(apiKey),
			Timeout:   10 * time.Second,
		},
		organisation: org,
		project:      project,
	}
}

type Client struct {
	client       *http.Client
	ticker       *time.Ticker
	organisation string
	project      string
}

func (c *Client) SetTicker(t *time.Ticker) { c.ticker = t }

func (c *Client) makeURL(path string, args ...interface{}) string {
	return fmt.Sprintf(
		"https://api.transifex.com/organizations/%s/projects/%s/%s",
		c.organisation, c.project, fmt.Sprintf(path, args...),
	)
}
func (c *Client) legacyURL(path string, args ...interface{}) string {
	return fmt.Sprintf(
		"https://www.transifex.com/api/2/project/%s/%s",
		c.project, fmt.Sprintf(path, args...),
	)
}

func (c *Client) request(method, url string, r, v interface{}) error {
	var body io.Reader
	if r != nil {
		b := bytes.NewBuffer(nil)
		if err := json.NewEncoder(b).Encode(r); err != nil {
			return fmt.Errorf("encode request: %s", err.Error())
		}
		body = b
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("create request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	if c.ticker != nil {
		<-c.ticker.C
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("execute: %s", err.Error())
	}
	b, ok := v.(*[]byte)
	if !ok {
		return c.decodeResp(resp, v)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error with %s", req.URL)
	}
	*b, err = ioutil.ReadAll(resp.Body)
	return err
}

func (c *Client) decodeResp(r *http.Response, v interface{}) error {
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		var errResp ErrResponse
		if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
			b, _ := ioutil.ReadAll(r.Body)
			fmt.Println(r.Request.URL, string(b))
			return fmt.Errorf("error decode: %s", err.Error())
		}
		return errResp
	}
	if v != nil {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			return fmt.Errorf("response decode: %s", err.Error())
		}
	}
	return nil
}

func (c *Client) Project() (p Project, err error) {
	return p, c.request("GET", c.legacyURL(""), nil, &p)
}

func (c *Client) AddLanguage(lang string, coord ...string) (err error) {
	return c.request("POST", c.legacyURL("languages/"), map[string]interface{}{"language_code": lang, "coordinators": coord}, nil)
}

func (c *Client) ListResources() (r []Resource, err error) {
	return r, c.request("GET", c.legacyURL("resources/"), nil, &r)
}

func (c *Client) UpdateName(slug, name string) (err error) {
	return c.request("PUT", c.legacyURL("resource/%s/", slug), map[string]string{"slug": slug, "name": name}, nil)
}

func (c *Client) ResourceDetail(slug string) (r ResourceDetail, err error) {
	return r, c.request("GET", c.makeURL("resources/%s/", slug), nil, &r)
}

func (c *Client) CreateResource(u UploadResourceRequest) (r Response, err error) {
	return r, c.request("POST", c.legacyURL("resources/"), u, &r)
}

func (c *Client) UpdateResource(slug, content string) (r Response, err error) {
	return r, c.request("PUT", c.legacyURL("resource/%s/content/", slug), map[string]string{"slug": slug, "content": content}, &r)
}

func (c *Client) DeleteResource(slug string) (err error) {
	return c.request("DELETE", c.legacyURL("resource/%s/", slug), nil, nil)
}

func (c *Client) UpdateTranslation(slug, lang, content string) (r Response, err error) {
	return r, c.request("PUT", c.legacyURL("resource/%s/translation/%s", slug, lang), map[string]string{"content": content}, &r)
}

func (c *Client) GetTranslation(slug, lang string) (r map[string]interface{}, err error) {
	return r, c.request("GET", c.legacyURL("resource/%s/translation/%s", slug, lang), nil, &r)
}

func (c *Client) GetTranslationFile(slug, lang string) (b []byte, err error) {
	return b, c.request("GET", c.legacyURL("resource/%s/translation/%s?file", slug, lang), nil, &b)
}

func (c *Client) GetStrings(slug, lang string) (r []ResourceString, err error) {
	return r, c.request("GET", c.legacyURL("resource/%s/translation/%s/strings/", slug, lang), nil, &r)
}

func (c *Client) SetStringTags(slug, hash string, tags ...string) (err error) {
	return c.request("PUT", c.legacyURL("resource/%s/source/%s/", slug, hash), map[string][]string{"tags": tags}, nil)
}
