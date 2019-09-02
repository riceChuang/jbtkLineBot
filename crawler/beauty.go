package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"net/http"
	"strings"
)

type BeautyCrawler struct {
	ContentUrl chan string
	ImageUrl   chan string
	Stop       chan bool
	db         *boltdb.Boltdb
}

var (
	ImageLength = 0
	maxImage    = 200
)

func NewBeautyCrawler(db *boltdb.Boltdb) *BeautyCrawler {
	b := &BeautyCrawler{
		ContentUrl: make(chan string, 3000),
		ImageUrl:   make(chan string, 3000),
		Stop:       make(chan bool, 2),
		db:         db,
	}

	for i := 0; i < 1; i++ {
		go b.RunContentPage()
	}

	for i := 0; i < 3; i++ {
		go b.RunImagePage()
	}

	go func() {
		for {
			select {
			case <-b.Stop:
				for len(b.ImageUrl) > 1 {
					//log.Printf("image len: %v", len(b.ImageUrl))
					<-b.ImageUrl
				}
				for len(b.ContentUrl) > 1 {
					//log.Printf("content len: %v", len(b.ContentUrl))
					<-b.ContentUrl
				}
			}
		}
	}()

	return b
}

func (b *BeautyCrawler) RunImage(url string) {
	b.GetMainPage(url)
}

func (b *BeautyCrawler) GetMainPage(url string) {

	// Instantiate default collector
	mainPage := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.ptt.cc"),
	)
	mainPage.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link

		if strings.Contains(e.Text, "上頁") && ImageLength < maxImage {
			//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			mainPage.Visit(e.Request.AbsoluteURL(link))
		} else {
			return
		}
	})

	mainPage.OnResponse(func(r *colly.Response) {
		//log.Println("https://" + r.Request.URL.Host + r.Request.URL.Path)
		b.addContentUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
	})
	cookies := []*http.Cookie{
		{
			Name:   "over18",
			Value:  "1",
			Domain: "www.ptt.cc",
		},
	}
	mainPage.SetCookies(url, cookies)
	mainPage.Visit(url)
	mainPage.Wait()
}

func (b *BeautyCrawler) RunContentPage() {

	for url := range b.ContentUrl {
		if ImageLength > maxImage {
			b.Stop <- true
			break
		}
		// Instantiate default collector
		contentPage := colly.NewCollector(
			// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
			colly.AllowedDomains("www.ptt.cc"),
		)

		contentPage.OnHTML("div[class]", func(e *colly.HTMLElement) {
			class := e.Attr("class")
			if class == "title" {
				link := e.ChildAttr("a[href]", "href")
				//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
				contentPage.Visit(e.Request.AbsoluteURL(link))
			}
		})

		contentPage.OnResponse(func(r *colly.Response) {
			b.addImageUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
		})
		cookies := []*http.Cookie{
			{
				Name:   "over18",
				Value:  "1",
				Domain: "www.ptt.cc",
			},
		}
		contentPage.SetCookies(url, cookies)
		contentPage.Visit(url)
		contentPage.Wait()
	}
}

func (b *BeautyCrawler) RunImagePage() {

	for url := range b.ImageUrl {
		if ImageLength > maxImage {
			b.Stop <- true
			break
		}
		imagePage := colly.NewCollector(
			// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
			colly.AllowedDomains("www.ptt.cc"),
		)

		imagePage.OnHTML("img", func(e *colly.HTMLElement) {
			link := e.Attr("src")
			if strings.Contains(link, ".jpg") && strings.Contains(link, "https") {
				ImageLength++
				beautyKey := fmt.Sprintf("beauty-%d", ImageLength)
				b.db.Insert(beautyKey, link)
				if ImageLength%50 == 0 {
					fmt.Println("now beauty images len : %d", ImageLength)
				}

			}
		})
		cookies := []*http.Cookie{
			{
				Name:   "over18",
				Value:  "1",
				Domain: "www.ptt.cc",
			},
		}
		imagePage.SetCookies(url, cookies)
		imagePage.Visit(url)
		imagePage.Wait()
	}
}

func (b *BeautyCrawler) addContentUrl(url string) {
	b.ContentUrl <- url
}

func (b *BeautyCrawler) addImageUrl(url string) {
	b.ImageUrl <- url
}
