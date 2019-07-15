package ldjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Article struct {
	Context          string           `json:"@context"`
	Type             string           `json:"@type"`
	MainEntityOfPage MainEntityOfPage `json:"mainEntityOfPage"`

	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ThumbnailUrl []string `json:"thumbnailUrl"`
	UploadDate   string   `json:"uploadDate"`
	ContentUrl   string   `json:"contentUrl"`
	EmbedUrl     string   `json:"embedUrl"`

	Headline    string `json:"headline"`
	ArticleBody string `json:"articleBody"`

	Authors Authors `json:"author"`
}

type MainEntityOfPage struct {
	Type string `json:"@type"`
	Id   string `json:"@id"`
	Url  string `json:"url"`
}

type Author struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

type Authors []Author

func (as *Authors) UnmarshalJSON(b []byte) error {
	type Alias Authors

	var tmp Alias

	if err := json.Unmarshal(b, &tmp); err == nil {
		*as = (Authors)(tmp)

		return nil
	}

	var a Author

	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	*as = Authors{a}

	return nil
}

var ErrNotFound = errors.New("Nicht gefunden")

func (a Article) Plaintext() string {
	var authors string

	for i, author := range a.Authors {
		authors = authors + strings.TrimSpace(author.Name)

		if i+1 == len(a.Authors) {
			continue
		}

		authors = authors + ", "
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n\n",
		strings.TrimSpace(a.Headline),
		authors,
		strings.TrimSpace(a.ArticleBody),
	)
}

// ldjsonFromNode recursively scans an HTML-Node for JSON-LD content and
// returns the first occurence. If not JSON-LD content can be found it returns
// an error.
func ldjsonFromNode(n *html.Node) (string, error) {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, a := range n.Attr {
			if a.Val == "application/ld+json" {
				return n.FirstChild.Data, nil
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if s, err := ldjsonFromNode(c); err == nil {
			return s, nil
		}
	}

	return "", ErrNotFound
}

// GetArticle scans a webpage for an Article in the JSON-LD format and returns
// the first occurence. Otherwise it returns an error.
func GetArticle(url string) (Article, error) {
	var res Article

	resp, err := http.Get(url)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return res, err
	}

	raw, err := ldjsonFromNode(node)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return res, err
	}

	if res.Type == "Article" || res.Type == "NewsArticle" {
		return res, nil
	} else {
		return res, ErrNotFound
	}
}
