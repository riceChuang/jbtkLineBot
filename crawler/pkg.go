package crawler

import "github.com/riceChuang/jbtkLineBot/boltdb"

var(
	beautyCrawler *BeautyCrawler
	dcardSex *DcardCrawler
	joker *JokerCrawler
	pornHub *PornHubCrawler
)

func Initialize(){
	beautyCrawler = NewBeautyCrawler(boltdb.DB())
	dcardSex = NewDcrdCrawler(boltdb.DB())
	joker = NewJokerCrawler(boltdb.DB())
	//pornHub = NewPornHubCrawler(boltdb.DB())
}

