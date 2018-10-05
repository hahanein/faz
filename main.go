package main // import "github.com/hahanein/faz"

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("http://www.faz.net/rss/aktuell/politik/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var v RSS
	err = xml.Unmarshal(body, &v)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	var as []LdJson
	mux := new(sync.Mutex)

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
				mux.Lock()
				as = append(as, data)
				mux.Unlock()
			}
		}(i)
	}

	wg.Wait()

	for _, a := range as {
		fmt.Println(strings.ToUpper(a.Headline), "\n")
		fmt.Println(a.ArticleBody, "\n\n")
	}
}

type Article struct {
	Title string
	Body  string
}

func Data(link string) (string, error) {
	var res string

	resp, err := http.Get(link)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return res, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Val == "application/ld+json" {
					res = n.FirstChild.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return res, nil
}
