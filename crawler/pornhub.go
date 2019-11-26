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
	b := &PornHubCrawler{
		ContentUrl: make(chan string, 3000),
		ImageUrl:   make(chan string, 3000),
		Stop:       make(chan bool, 2),
		db:         db,
	}
	
	return b
}

func (b *PornHubCrawler) RunImage(url string) {
	b.GetMainPage(url)
}

func (b *PornHubCrawler) GetMainPage(url string) {
}