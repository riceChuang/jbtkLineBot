package crawler

import (
	"github.com/gocolly/colly"
	"github.com/riceChuang/jbtkLineBot/boltdb"

	"fmt"
	"strings"
)

type JokerCrawler struct {
	ContentUrl chan string
	ImageUrl   chan string
	db         *boltdb.Boltdb
}

var (
	JokerLenght   int
	titlePosition = 0
	JokerMap = map[string]string{}
)

func NewJokerCrawler(db *boltdb.Boltdb) *JokerCrawler {
	j := &JokerCrawler{
		ContentUrl: make(chan string, 1000),
		ImageUrl:   make(chan string, 1000),
		db:         db,
	}

	for i := 0; i < 2; i++ {
		go j.RunContenPage()
	}

	for i := 0; i < 3; i++ {
		go j.RunTextPage()
	}

	return j
}

func (j *JokerCrawler) RunImage(url string) {
	j.GetmainPage(url)
}

func (j *JokerCrawler) GetmainPage(url string) {

	// Instantiate default collector
	mainPage := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.ptt.cc"),
	)

	mainPage.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link

		if strings.Contains(e.Text, "上頁") {
			//fmt.Printf("Link found: %q -> %s\n", e.Text, link)

			mainPage.Visit(e.Request.AbsoluteURL(link))
		}
	})

	mainPage.OnResponse(func(r *colly.Response) {
		if JokerLenght > 500 {
			return
		}
		j.addContentUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
	})

	mainPage.Visit(url)
	mainPage.Wait()
}

func (j *JokerCrawler) RunContenPage() {

	for url := range j.ContentUrl {
		if JokerLenght > 500 {
			return
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
			//b.GetImagePage("https://" + r.Request.URL.Host + r.Request.URL.Path)
			j.addImageUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
		})

		contentPage.Visit(url)
		contentPage.Wait()
	}
}

func (j *JokerCrawler) RunTextPage() {

	for url := range j.ImageUrl {
		if JokerLenght > 500 {
			return
		}
		textPage := colly.NewCollector(
			// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
			colly.AllowedDomains("www.ptt.cc"),
		)

		textPage.OnHTML("#main-content", func(e *colly.HTMLElement) {

			titleValue := e.ChildText("div[class='article-metaline']:nth-child(3)>span[class='article-meta-value']")
			if !strings.Contains(titleValue, "笑話"){
				return
			}
			timeValue := e.ChildText("div[class='article-metaline']:nth-child(4)>span[class='article-meta-value']")
			//找文章的字首跟字尾
			contentFirstPos := strings.Index(e.Text, timeValue)
			if contentFirstPos == -1 {
				return
			}
			contentLastPos := strings.Index(e.Text, "--")
			if contentLastPos == -1 {
				return
			}
			//取代換行
			content := strings.Replace(e.Text[contentFirstPos+len(timeValue):contentLastPos], "\n", " ", -1)
			content = strings.Replace(content, "              ", "\n", -1)
			pageContent := fmt.Sprintf("%v\n%v", titleValue, content)
			//fmt.Println(pageContent)

			JokerLenght++
			jokerKey := fmt.Sprintf("joker-%d", JokerLenght)
			JokerMap[jokerKey] = pageContent
			if JokerLenght%50 == 0 {
				fmt.Println("now JokerLenght len : %d", JokerLenght)
			}
		})

		textPage.Visit(url)
		textPage.Wait()
	}
}

func (j *JokerCrawler) addContentUrl(url string) {
	j.ContentUrl <- url
}

func (j *JokerCrawler) addImageUrl(url string) {
	j.ImageUrl <- url
}
