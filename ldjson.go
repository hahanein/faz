package main

type LdJson struct {
	Context string `json:"@context"`
	Type    string `json:"@type"`

	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ThumbnailUrl []string `json:"thumbnailUrl"`
	UploadDate   string   `json:"uploadDate"`
	ContentUrl   string   `json:"contentUrl"`
	EmbedUrl     string   `json:"embedUrl"`

	Headline    string `json:"headline"`
	ArticleBody string `json:"articleBody"`
}
