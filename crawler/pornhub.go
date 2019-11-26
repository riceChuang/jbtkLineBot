package crawler

import (
	"github.com/riceChuang/jbtkLineBot/boltdb"
)

type PornHubCrawler struct {
	ContentUrl chan string
	ImageUrl   chan string
	Stop       chan bool
	db         *boltdb.Boltdb
}



func NewPornHubCrawler(db *boltdb.Boltdb) *PornHubCrawler {
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

func (b *PornHubCrawler) RunImage(url string) {
	b.GetMainPage(url)
}

func (b *PornHubCrawler) GetMainPage(url string) {
}