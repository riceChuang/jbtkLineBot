package crawler

import (
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"sync"
)

type CrawlerType int

const (
	Beauty CrawlerType = iota + 1
	DcardSex
	Joker
	PornHub
)

type Crawler interface {
	RunCrawlerImage(url string)
	GetImageLength() int32
}

var (
	crawlerInstanceMap = map[CrawlerType]Crawler{}
	mu                 = sync.Mutex{}
)

func GetCrawlerByType(sourceType CrawlerType) (instance Crawler, err error) {
	mu.Lock()
	defer mu.Unlock()

	var exist bool

	if instance, exist = crawlerInstanceMap[sourceType]; !exist {
		instance, err = createCrawlerBySourceType(sourceType)
		if err != nil {
			return nil, err
		}
		if instance != nil {
			crawlerInstanceMap[sourceType] = instance
		}
	}

	return

}

func createCrawlerBySourceType(sourceType CrawlerType) (crawler Crawler, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	switch sourceType {
	case Beauty:
		crawler = NewBeautyCrawler(boltdb.DB())
	case DcardSex:
		crawler = NewDcardCrawler(boltdb.DB())
	case Joker:
		crawler = NewJokerCrawler(boltdb.DB())
	case PornHub:
		crawler = NewPornHubCrawler(boltdb.DB())
	default:
	}

	return
}
