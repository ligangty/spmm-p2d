package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseContent(pageContent string, top int) []map[string]string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageContent))
	if err != nil {
		fmt.Printf("Error parsing the html content: %s\n", err)
		os.Exit(1)
	}
	reposInfo := []map[string]string{}
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		info := map[string]string{}
		url, _ := s.Find("h1 > a[data-hydro-click]").Attr("href")
		info["url"] = url[1:]
		reposInfo = append(reposInfo, info)
		description := strings.TrimSpace(s.Find("p").Text())
		info["description"] = description
		language := strings.TrimSpace(s.Find("span[itemprop=programmingLanguage]").Text())
		info["language"] = language
		s.Find("a[class='Link--muted d-inline-block mr-3']").Each(func(i int, s *goquery.Selection) {
			content := strings.TrimSpace(s.Text())
			switch i {
			case 0:
				info["star"] = content
			case 1:
				info["fork"] = content
			default:
			}
		})
	})
	return reposInfo[:top]
}
