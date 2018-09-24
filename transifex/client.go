package transifex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	organisation string
	project      string
}

func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("https://www.transifex.com/api/2/project/%s/%s", c.project, path)
}

func (c *Client) request(method, path string, r, v interface{}) error {
	var body io.Reader
	if r != nil {
		b := bytes.NewBuffer(nil)
		if err := json.NewEncoder(b).Encode(r); err != nil {
			return fmt.Errorf("encode request: %s", err.Error())
		}
		body = b
	}

	req, err := http.NewRequest(method, c.buildURL(path), body)
	if err != nil {
		return fmt.Errorf("create request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("execute: %s", err.Error())
	}
	return c.decodeResp(resp, v)
}

func (c *Client) decodeResp(r *http.Response, v interface{}) error {
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	if r.StatusCode >= 400 {
		var errResp ErrResponse
		if err := d.Decode(&errResp); err != nil {
			return fmt.Errorf("error decode: %s", err.Error())
		}
		return errResp
	}
	if v != nil {
		if err := d.Decode(v); err != nil {
			return fmt.Errorf("response decode: %s", err.Error())
		}
	}
	return nil
}

func (c *Client) Project() (p Project, err error) {
	return p, c.request("GET", "", nil, &p)
}

func (c *Client) ListResources() (r []Resource, err error) {
	return r, c.request("GET", "resources/", nil, &r)
}

func (c *Client) CreateResource(u UploadResourceRequest) (r Response, err error) {
	return r, c.request("POST", "resources/", u, &r)
}

func (c *Client) UpdateResource(slug, content string) (r Response, err error) {
	return r, c.request("PUT", fmt.Sprintf("resource/%s/content/", slug), map[string]string{"slug": slug, "content": content}, &r)
}

func (c *Client) DeleteResource(slug string) (err error) {
	return c.request("DELETE", fmt.Sprintf("resource/%s/", slug), nil, nil)
}

func (c *Client) UpdateTranslation(slug, lang, content string) (r Response, err error) {
	return r, c.request("PUT", fmt.Sprintf("resource/%s/translation/%s", slug, lang), map[string]string{"content": content}, &r)
}

func (c *Client) GetTranslation(slug, lang string) (r map[string]interface{}, err error) {
	return r, c.request("GET", fmt.Sprintf("resource/%s/translation/%s", slug, lang), nil, &r)
}
