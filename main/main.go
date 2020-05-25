package main

import (
	"dmhyCrawler/conf"
	"dmhyCrawler/crawler"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	start := time.Now()

	serviceConf := conf.GetServiceConfig()

	for _, bangumiConfig := range serviceConf.BangumiConfigs {
		if bangumiConfig.SearchExp != "" && bangumiConfig.Title != "" {
			wg.Add(1)
			go func(title string, searchExp string) {
				url := "https://share.dmhy.org/topics/list?keyword=" + searchExp
				myCrawler := crawler.NewCrawler(title, url)
				myCrawler.FetchMagLinks()
				defer wg.Done()
			}(bangumiConfig.Title, bangumiConfig.SearchExp)
		}
	}
	wg.Wait()
	timeUsed := time.Now().Sub(start)
	fmt.Printf("job finished, time used %v\n", timeUsed)
}
