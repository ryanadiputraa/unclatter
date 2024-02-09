package scrapper

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

const (
	heading      = "h1, h2, h3, h4, h5, h6"
	paragraphs   = "p, blockquote, pre, div"
	styled       = "em, strong, i, b, mark, small, sub, sup, u, s, strike"
	links        = "a"
	lists        = "ul, ol, li"
	containers   = "span, section, article, main"
	quotations   = "q, cite"
	textGrouping = "br, hr"
	code         = "code, var"
)

type Scrapper interface {
	ScrapeTextContent(url string) (string, error)
}

type scrapper struct {
	c *colly.Collector
}

func NewScrapper() Scrapper {
	return &scrapper{
		c: colly.NewCollector(),
	}
}

func (s *scrapper) ScrapeTextContent(url string) (content string, err error) {
	s.c.OnHTML("body", func(h *colly.HTMLElement) {
		selectors := heading + paragraphs + styled + links + containers +
			quotations + textGrouping + code
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
