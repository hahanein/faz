package main // import "github.com/hahanein/faz"

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	var v RSS

	resp, err := http.Get("http://www.faz.net/rss/aktuell/politik/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(body, &v)
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range v.Channel.Items {
		fmt.Println(i.Title)
	}

	fmt.Println()
	fmt.Println()

	for _, i := range v.Channel.Items {
		fmt.Println(strings.ToUpper(i.Title))
		fmt.Println()

		raw, err := Data(i.Link)
		if err != nil {
			log.Fatal(err)
		}

		var data LdJson

		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			log.Fatal(err)
		}

		fmt.Println("\t", data.ArticleBody)
		fmt.Println()
		fmt.Println()
	}
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
