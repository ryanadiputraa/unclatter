package scrapper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Scrapper interface {
	ScrapeTextContent(url string) (string, error)
}

type scrapper struct {
	c *colly.Collector
}

func NewScrapper() Scrapper {
	c := colly.NewCollector()
	c.AllowURLRevisit = true
	return &scrapper{
		c: c,
	}
}

func (s *scrapper) ScrapeTextContent(url string) (content string, err error) {
	s.c.OnHTML("body", func(h *colly.HTMLElement) {
		selectors := "p, blockquote, pre, code, var"
		h.ForEach(selectors, func(i int, h *colly.HTMLElement) {
			content += fmt.Sprintf("<%s>%s</%s>", h.Name, h.Text, h.Name)
		})
	})

	s.c.OnError(func(r *colly.Response, e error) {
		err = e
	})

	err = s.c.Visit(url)

	return content, err
}
