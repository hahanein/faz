package main // import "github.com/hahanein/faz"

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/hahanein/faz/ldjson"
	"github.com/hahanein/faz/rss"
)

const VERSION = "2.0.0"

func main() {
	var (
		printVersion = flag.Bool("version", false, "Versionsnummer anzeigen und Programm beenden")

		resort = flag.String("resort", "", "Resort")
		format = flag.String("format", "plaintext", "Ausgabeformat")
	)

	flag.Parse()

	if *printVersion {
		fmt.Printf("faz v%s\n", VERSION)
		return
	}

	r, err := rss.Get(fmt.Sprintf("https://www.faz.net/rss/aktuell/%s", *resort))
	if err != nil {
		panic(err)
	}

	var as []ldjson.Article
	for _, i := range r.Channel.Items {
		a, err := ldjson.GetArticle(i.Link)
		if err == ldjson.ErrNotFound {
			continue
		} else if err != nil {
			panic(err)
		}

		as = append(as, a)
	}

	switch *format {
	case "plaintext":
		for _, a := range as {
			fmt.Print(a.Plaintext())
		}

	case "json":
		enc := json.NewEncoder(os.Stdout)

		if err := enc.Encode(&as); err != nil {
			panic(err)
		}

	default:
		panic("Unbekanntes Ausgabeformat")
	}
}
