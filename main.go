package main

import (
	"fmt"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"github.com/riceChuang/jbtkLineBot/config"
	"github.com/riceChuang/jbtkLineBot/crawler"
	"github.com/riceChuang/jbtkLineBot/line"
	"github.com/riceChuang/jbtkLineBot/route"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {

	boltdb.Initialize()
	config.InitialConfigPkg()

	config := config.GetConfig()
	crawlerTypesMap := map[crawler.CrawlerType]string{
		crawler.DcardSex: config.DcardUrl,
		//crawler.Beauty:   config.BeautyUrl,
	}

	for crawlerType, url := range crawlerTypesMap {
		crawlerWorker, err := crawler.GetCrawlerByType(crawlerType)
		if err != nil {
			logrus.Error("getCrawler Err:%v", err)
		}
		crawlerWorker.RunCrawlerImage(url)
	}

	line.InitBotClient()
	http.HandleFunc("/back", route.CallbackHandler)
	port := os.Getenv("PORT")
	addr := ""
	if port != "" {
		addr = fmt.Sprintf(":%s", port)
	} else {
		addr = ":3000"
	}

	http.ListenAndServe(addr, nil)
}

//
//func PrintMemUsage() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
//	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
//	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
//	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
//	fmt.Printf("\tNumGC = %v\n", m.NumGC)
//}
//func bToMb(b uint64) uint64 {
//	return b / 1024 / 1024
//}
