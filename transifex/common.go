package transifex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ErrResponse struct {
	ErrorCode string   `json:"error_code,omitempty"`
	Detail    string   `json:"detail,omitempty"`
	Priority  []string `json:"priority,omitempty"`
}

func (e ErrResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Detail)
}

type apiAuth string

func (a apiAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth("api", string(a))
	return http.DefaultTransport.RoundTrip(req)
}

type Project struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Slug      string        `json:"slug"`
	Tags      []interface{} `json:"tags"`
	Languages []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"languages"`
	TotalResources int `json:"total_resources"`
	SourceLanguage struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"source_language"`
	Type            string      `json:"type"`
	LogoURL         string      `json:"logo_url"`
	Description     string      `json:"description"`
	Stringcount     int         `json:"stringcount"`
	Wordcount       int         `json:"wordcount"`
	LongDescription string      `json:"long_description"`
	WebsiteURL      string      `json:"website_url"`
	Maintainers     []string    `json:"maintainers"`
	Created         time.Time   `json:"created"`
	LastUpdate      interface{} `json:"last_update"`
	Private         bool        `json:"private"`
	RepositoryURL   string      `json:"repository_url"`
	Archived        bool        `json:"archived"`
	Team            struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"team"`
	Stats struct {
		Es struct {
			Translated struct {
				Stringcount  int         `json:"stringcount"`
				LastActivity interface{} `json:"last_activity"`
				Name         string      `json:"name"`
				Wordcount    int         `json:"wordcount"`
				Percentage   float64     `json:"percentage"`
			} `json:"translated"`
			Reviewed1 struct {
				Stringcount  int         `json:"stringcount"`
				LastActivity interface{} `json:"last_activity"`
				Name         string      `json:"name"`
				Wordcount    int         `json:"wordcount"`
				Percentage   float64     `json:"percentage"`
			} `json:"reviewed_1"`
		} `json:"es"`
	} `json:"stats"`
}

type BaseResource struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	I18nType string `json:"i18n_type"`
	Priority string `json:"priority"`
	Category string `json:"category"`
}

type Resource struct {
	BaseResource
	SourceLanguage string `json:"source_language_code"`
}

type UploadResourceRequest struct {
	BaseResource
	Content            string `json:"content"`
	AcceptTranslations bool   `json:"accept_translations"`
}

type Language struct {
	Coordinators []string `json:"coordinators"`
	LanguageCode string   `json:"language_code"`
	Translators  []string `json:"translators"`
	Reviewers    []string `json:"reviewers"`
}

type Response struct {
	Added   int `json:"strings_added"`
	Updated int `json:"strings_updated"`
	Deleted int `json:"strings_delete"`
}

func (r *Response) UnmarshalJSON(raw []byte) error {
	var dst interface{}
	if err := json.Unmarshal(raw, &dst); err != nil {
		return err
	}
	switch v := dst.(type) {
	case []interface{}:
		r.Added = int(v[0].(float64))
		r.Updated = int(v[1].(float64))
		r.Deleted = int(v[2].(float64))
	case map[string]interface{}:
		r.Added = int(v["strings_added"].(float64))
		r.Updated = int(v["strings_updated"].(float64))
		r.Deleted = int(v["strings_delete"].(float64))
	default:
		return fmt.Errorf("Unkwown type %T", v)
	}
	return nil
}
