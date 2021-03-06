package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"github.com/riceChuang/jbtkLineBot/config"
	"strings"
	"sync"
)

type JokerCrawler struct {
	ContentUrl  chan string
	ImageUrl    chan string
	maxLen      int32
	imageLength int32
	Stop        chan bool
	db          *boltdb.Boltdb
}

var (
	jokerCraw        *JokerCrawler
	jokerCrawlerOnce = &sync.Once{}
	JokerMap         = map[string]string{}
)

func NewJokerCrawler(db *boltdb.Boltdb) *JokerCrawler {
	cfg := config.GetConfig()
	if jokerCraw != nil {
		return jokerCraw
	}
	jokerCrawlerOnce.Do(func() {
		jokerCraw = &JokerCrawler{
			ContentUrl: make(chan string, 1000),
			ImageUrl:   make(chan string, 1000),
			Stop:       make(chan bool, 2),
			maxLen:     cfg.MaxJokerLen,
			db:         db,
		}
		jokerCraw.Initialize()
	})
	return jokerCraw
}

func (j *JokerCrawler) Initialize() {
	for i := 0; i < 1; i++ {
		go j.RunContenPage()
	}

	for i := 0; i < 2; i++ {
		go j.RunTextPage()
	}

	go func() {
		for {
			select {
			case <-j.Stop:
				for len(j.ContentUrl) > 0 {
					//log.Printf("content len: %v", len(j.ContentUrl))
					<-j.ContentUrl
				}
				for len(j.ImageUrl) > 0 {
					//log.Printf("image len: %v", len(j.ImageUrl))
					<-j.ImageUrl
				}
			}
		}
	}()
}

func (j *JokerCrawler) RunCrawlerImage(url string) {
	go  j.GetmainPage(url)
}

func (j *JokerCrawler) GetImageLength() int32 {
	return j.imageLength
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

		if strings.Contains(e.Text, "上頁") && j.imageLength < j.maxLen {
			//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			mainPage.Visit(e.Request.AbsoluteURL(link))
		}
	})

	mainPage.OnResponse(func(r *colly.Response) {
		j.addContentUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
	})

	mainPage.Visit(url)
	mainPage.Wait()
}

func (j *JokerCrawler) RunContenPage() {

	for url := range j.ContentUrl {
		if j.imageLength > j.maxLen {
			j.Stop <- true
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
			j.addImageUrl("https://" + r.Request.URL.Host + r.Request.URL.Path)
		})

		contentPage.Visit(url)
		contentPage.Wait()
	}
}

func (j *JokerCrawler) RunTextPage() {

	for url := range j.ImageUrl {
		if j.imageLength > j.maxLen {
			j.Stop <- true
			break
		}
		textPage := colly.NewCollector(
			// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
			colly.AllowedDomains("www.ptt.cc"),
		)

		textPage.OnHTML("#main-content", func(e *colly.HTMLElement) {

			titleValue := e.ChildText("div[class='article-metaline']:nth-child(3)>span[class='article-meta-value']")
			if !strings.Contains(titleValue, "笑話") {
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

			j.imageLength++
			jokerKey := fmt.Sprintf("joker-%d", j.imageLength)
			JokerMap[jokerKey] = pageContent
			if j.imageLength%50 == 0 {
				fmt.Println("now JokerLenght len : %d", j.imageLength)
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
