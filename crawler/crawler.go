package crawler

import "sync"

type Type int

const (
	Beauty Type = 1
	DcardSex Type = 2
)

type Crawler interface {
	RunImage(url string)
}

var (
	crawlerInstanceMap = map[Type]Crawler{}
	mu                 = sync.Mutex{}
)

func GetCrawlerByType(sourceType Type) (instance Crawler, err error) {
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

func createCrawlerBySourceType(sourceType Type) (crawler Crawler, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	switch sourceType {
	case Beauty:
		crawler = _b
	case DcardSex:
		crawler =  _dsex
	default:
	}

	return
}
