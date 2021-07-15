package crawler

import (
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"sync"
)

type PornHubCrawler struct {
	ContentUrl  chan string
	ImageUrl    chan string
	Stop        chan bool
	imageLength int32
	db          *boltdb.Boltdb
}

var (
	pornHubCraw        *PornHubCrawler
	pornHubCrawlerOnce = &sync.Once{}
)

func NewPornHubCrawler(db *boltdb.Boltdb) *PornHubCrawler {
	if pornHubCraw != nil {
		return pornHubCraw
	}
	pornHubCrawlerOnce.Do(func() {
		pornHubCraw = &PornHubCrawler{
			ContentUrl: make(chan string, 3000),
			ImageUrl:   make(chan string, 3000),
			Stop:       make(chan bool, 2),
			db:         db,
		}

	})
	return pornHubCraw
}

func (pc *PornHubCrawler) RunCrawlerImage(url string) {
	go pc.GetMainPage(url)
}

func (pc *PornHubCrawler) GetMainPage(url string) {
}

func (pc *PornHubCrawler) GetImageLength() int32 {
	return pc.imageLength
}
