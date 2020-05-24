package crawler

import (
	"dmhyCrawler/pipeline"
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

type Crawler struct {
	sync.RWMutex
	title       string
	listPageUrl string
}

func NewCrawler(title string, listPageUrl string) *Crawler {
	return &Crawler{
		title:       title,
		listPageUrl: listPageUrl,
	}
}

func (c *Crawler) GetListPageUrl() string {
	c.RLock()
	defer c.RUnlock()
	return c.listPageUrl
}

func (c *Crawler) GetTitle() string {
	c.RLock()
	defer c.RUnlock()
	return c.title
}

func (c *Crawler) FetchMagLinks() {
	listPageCollector := colly.NewCollector(
		colly.Async(true),
	)
	detailedPageCollector := listPageCollector.Clone()

	// Find and visit all links
	listPageCollector.OnHTML("tbody tr td.title a[target='_blank']", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		detailedPageCollector.Visit(link)
	})

	listPageCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	detailedPageCollector.OnHTML("#a_magnet", func(e *colly.HTMLElement) {
		item := pipeline.Item{
			Title:   c.title,
			MagLink: e.Attr("href"),
		}
		e.Response.Ctx.Put("pipelineItem", item)
	})

	detailedPageCollector.OnScraped(func(response *colly.Response) {
		item := response.Ctx.GetAny("pipelineItem").(pipeline.Item)

		pipeline.RunPipeline(item)
	})

	listPageCollector.Visit(c.GetListPageUrl())
	listPageCollector.Wait()
	detailedPageCollector.Wait()
	return
}
