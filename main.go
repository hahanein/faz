package main // import "github.com/hahanein/faz"

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hahanein/faz/ldjson"
	"github.com/hahanein/faz/rss"
	"golang.org/x/sync/errgroup"
)

const VERSION = "3.1.0"

func main() {
	var (
		printVersion = flag.Bool("version", false, "Versionsnummer anzeigen und Programm beenden")

		resorts = flag.String("resorts", "", "Resorts (getrennt durch Kommata)")
		format  = flag.String("format", "plaintext", "Ausgabeformat")
	)

	flag.Parse()

	if *printVersion {
		fmt.Printf("faz v%s\n", VERSION)
		return
	}

	as, err := getAllArticles(strings.Split(*resorts, ","))
	if err != nil {
		panic(err)
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

func getAllArticles(resorts []string) ([]ldjson.Article, error) {
	if len(resorts) == 0 {
		return getArticles("")
	}

	articlesByResort := make([][]ldjson.Article, len(resorts))

	var g errgroup.Group
	for i := range resorts {
		n := i
		g.Go(func() error {
			as, err := getArticles(resorts[n])
			if err != nil {
				return err
			}

			articlesByResort[n] = as

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return mergeArticles(articlesByResort), nil
}

func getArticles(resort string) ([]ldjson.Article, error) {
	var res []ldjson.Article

	r, err := rss.Get(fmt.Sprintf("https://www.faz.net/rss/aktuell/%s", resort))
	if err != nil {
		return res, err
	}

	for _, i := range r.Channel.Items {
		a, err := ldjson.GetArticle(i.Link)
		if err == ldjson.ErrNotFound {
			continue
		} else if err != nil {
			return res, err
		}

		res = append(res, a)
	}

	return res, nil
}

func mergeArticles(xs [][]ldjson.Article) []ldjson.Article {
	var y []ldjson.Article

	for i, _ := range xs {
		for j, _ := range xs[i] {
			missing := true

			for _, a := range y {
				if a.MainEntityOfPage.Id == xs[i][j].MainEntityOfPage.Id {
					missing = false
					break
				}
			}

			if missing {
				y = append(y, xs[i][j])
			}
		}
	}

	return y
}
