package crawler

import (
	"github.com/gocolly/colly"
	"fmt"

)




func GetFileList(){
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.ptt.cc"),
		colly.MaxDepth(1),
	)
	//https://www.ptt.cc/bbs/Beauty/index.html
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link

		if e.Text != "批踢踢實業坊"{
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {

	})
	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.ptt.cc/bbs/Beauty/index.html")
}