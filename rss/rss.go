package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Get fetches an RSS-Feed from a given URL.
func Get(url string) (RSS, error) {
	var res RSS

	resp, err := http.Get(url)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = xml.Unmarshal(body, &res)

	return res, err
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type Image struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Link  string `xml:"link"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Channel struct {
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Copyright     string   `xml:"copyright"`
	Category      string   `xml:"category"`
	Language      string   `xml:"language"`
	Docs          string   `xml:"docs"`
	Generator     string   `xml:"generator"`
	TTL           string   `xml:"ttl"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Image         Image    `xml:"image"`
	AtomLink      AtomLink `xml:"atom:link"`
	Items         []Item   `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}
