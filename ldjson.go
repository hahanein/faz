package main

const dat = `
{"@context":"http://schema.org","@type":"VideoObject","name":"Livestream: Staats- und Regierungschefs beim Nato-Gipfel in Brüssel","description":"Die 29 Staats- und Regierungschefs der Nato kommen in Brüssel zusammen. Bestimmt wird das Treffen von den Forderungen des amerikanischen Präsidenten, die Europäer müssten mehr Geld für ihre Verteidigung ausgeben. Verfolgen Sie den Gipfel im Livestream. ","thumbnailUrl":["https://www.faz.net/ppmedia/aktuell/2708388933/1.5685486/format_top1_portrait/54215390.jpg","https://www.faz.net/ppmedia/aktuell/2708388933/1.5685486/format4_3-schema_org_thumbnail/54215390.jpg","https://www.faz.net/ppmedia/aktuell/2708388933/1.5685486/mmo-uebersichtsseite-aufmacher-retina/54215390.jpg"],"uploadDate":"2018-07-11T10:11:10+0200","contentUrl":"http://player.zdf.de/zdf/faznet/index-faz.php?id\u003dphoenix-live-beitrag-100\u0026startTime\u003d\u0026stopTime\u003d","embedUrl":"http://www.faz.net/1.5685472?service\u003dembedded"}
`

type LdJson struct {
	Context string `json:"@context"`
	Type    string `json:"@type"`

	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ThumbnailUrl []string `json:"thumbnailUrl"`
	UploadDate   string   `json:"uploadDate"`
	ContentUrl   string   `json:"contentUrl"`
	EmbedUrl     string   `json:"embedUrl"`

	ArticleBody string `json:"articleBody"`
}

// Simple
//
// type LdJson struct {
// 	Context string `json:"@context"`
// 	Type    string `json:"@type"`
//
// 	Name         string   `json:"name"`
// 	Description  string   `json:"description"`
// 	ThumbnailUrl []string `json:"thumbnailUrl"`
// 	UploadDate   string   `json:"uploadDate"`
// 	ContentUrl   string   `json:"contentUrl"`
// 	EmbedUrl     string   `json:"embedUrl"`
// }
