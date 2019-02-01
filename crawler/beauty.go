package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

type BeautyCrawler struct {
	ContentUrl chan string
	ImageUrl   chan string
}

var ImageMap []string

func NewBeautyCrawler() *BeautyCrawler{
	b := &BeautyCrawler{
		ContentUrl: make(chan string, 10000),
		ImageUrl:   make(chan string, 10000),
	}

	for i := 0; i < 2; i++ {
		go b.RunContenPage()
	}

	for i := 0; i < 3; i++ {
		go b.RunImagePage()
	}

	return b
}

func (b *BeautyCrawler) RunImage(url string) {
	b.GetmainPage(url)
}

func (b *BeautyCrawler) GetmainPage(url string) {

	// Instantiate default collector
	mainPage := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.ptt.cc"),
	)

	mainPage.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		if strings.Contains(e.Text, "上頁") {
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			mainPage.Visit(e.Request.AbsoluteURL(link))
		}
	})

	mainPage.OnResponse(func(r *colly.Response) {
		b.addContentUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
	})

	mainPage.Visit(url)
	mainPage.Wait()
}

func (b *BeautyCrawler) RunContenPage() {

	for url := range b.ContentUrl {
		// Instantiate default collector
		contentPage := colly.NewCollector(
			// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
			colly.AllowedDomains("www.ptt.cc"),
		)

		contentPage.OnHTML("div[class]", func(e *colly.HTMLElement) {
			class := e.Attr("class")
			if class == "title" {
				link := e.ChildAttr("a[href]", "href")
				fmt.Printf("Link found: %q -> %s\n", e.Text, link)
				contentPage.Visit(e.Request.AbsoluteURL(link))
			}
		})

		contentPage.OnResponse(func(r *colly.Response) {
			//b.GetImagePage("https://" + r.Request.URL.Host + r.Request.URL.Path)
			b.addImageUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
		})

		contentPage.Visit(url)
		contentPage.Wait()
	}
}

func (b *BeautyCrawler) RunImagePage() {

	for url := range b.ImageUrl {
		imagePage := colly.NewCollector(
			// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
			colly.AllowedDomains("www.ptt.cc"),
		)

		imagePage.OnHTML("img", func(e *colly.HTMLElement) {
			link := e.Attr("src")
			if strings.Contains(link, ".jpg") && strings.Contains(link, "https") {
				fmt.Printf("my image: %q -> %s\n", e.Text, link)
				ImageMap = append(ImageMap, link)
				fmt.Println("now imagemap len : %v", len(ImageMap))
			}
		})

		imagePage.Visit(url)
		imagePage.Wait()
	}
}

func (b *BeautyCrawler) addContentUrl(url string ){
	b.ContentUrl <- url
}

func (b *BeautyCrawler) addImageUrl(url string){
	b.ImageUrl <- url
}