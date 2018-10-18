package main // import "github.com/hahanein/faz"

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("http://www.faz.net/rss/aktuell/politik/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var v RSS
	err = xml.Unmarshal(body, &v)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, i := range v.Channel.Items {
		wg.Add(1)
		go func(i Item) {
			defer wg.Done()

			raw, err := Data(i.Link)
			if err != nil {
				return
			}

			var data LdJson

			if err := json.Unmarshal([]byte(raw), &data); err != nil {
				return
			}

			if data.Type == "Article" || data.Type == "NewsArticle" {
				fmt.Print(data.Headline, "\n\n  ", data.ArticleBody, "\n\n")
			}
		}(i)
	}

	wg.Wait()
}

type Article struct {
	Title string
	Body  string
}

func DataFromNode(n *html.Node) (string, error) {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, a := range n.Attr {
			if a.Val == "application/ld+json" {
				return n.FirstChild.Data, nil
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if s, err := DataFromNode(c); err == nil {
			return s, nil
		}
	}

	return "", errors.New("JSON-LD not found")
}

func Data(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	return DataFromNode(node)
}
