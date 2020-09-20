package crawler

import (
	page_getter "dmhyCrawler/page-getter"
	"dmhyCrawler/pipeline"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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
	detailedPageCollector.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second,
	})

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("\\")))
	listPageCollector.WithTransport(t)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	listPageHtmlresponse := page_getter.GetResponse(c.GetListPageUrl())
	saveFile(dir, listPageHtmlresponse, c.GetTitle())

	// Find and visit all links
	listPageCollector.OnHTML("tbody tr td.title a[target='_blank']", func(e *colly.HTMLElement) {
		link := fmt.Sprintf("https://share.dmhy.org%s", e.Attr("href"))
		detailedPageCollector.Visit(link)
	})

	listPageCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err)
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
	urlDir := dir
	urlDir = strings.ReplaceAll(urlDir, "\\", "/")
	url := fmt.Sprintf("file://" + urlDir + "/" + c.GetTitle() + ".html")
	listPageCollector.Visit(url)
	listPageCollector.Wait()
	detailedPageCollector.Wait()

	defer removeFile(dir, c.GetTitle())
	return
}

func saveFile(dir string, htmlRes string, fileName string) {
	f, err := os.Create(dir + "\\" + fileName + ".html")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_, err = f.WriteString(htmlRes)
	if err != nil {
		panic(err)
	}
}

func removeFile(dir string, fileName string) {
	err := os.Remove(dir + "\\" + fileName + ".html")
	if err != nil {
		panic(err)
	}
}
