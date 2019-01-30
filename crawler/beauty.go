package crawler

import (
	"fmt"
	"github.com/gocolly/colly"

	"strings"
)

var ImageMap []string

func GetFileList() {

	// Instantiate default collector
	mainPage := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.ptt.cc"),
	)

	contentPage := mainPage.Clone()

	mainPage.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link

		if e.Text != "批踢踢實業坊" {
			//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			contentPage.Visit(e.Request.AbsoluteURL(link))
		}
	})

	contentPage.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, ".jpg") && strings.Contains(link, "https") {
			fmt.Printf("my image: %q -> %s\n", e.Text, link)
			ImageMap = append(ImageMap, link)
		}
	})

	mainPage.Visit("https://www.ptt.cc/bbs/Beauty/index.html")
	contentPage.Wait()
}
