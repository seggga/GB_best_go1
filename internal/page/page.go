package page

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// page holds a webpage
type page struct {
	doc *goquery.Document
}

// NewPage reads web-page's body
func NewPage(raw io.Reader) (*page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		return nil, err
	}
	return &page{
		doc: doc,
	}, nil
}

// GetTitle gets title of the 'page'
func (p *page) GetTitle() string {
	return p.doc.Find("title").First().Text()
}

// GetLinks collects a list of links found on the given 'page'
func (p *page) GetLinks() []string {
	var urls []string
	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if ok {
			// link validation
			parsedLink, err := url.Parse(link)
			if err != nil {
				return
			}
			if !parsedLink.IsAbs() {
				if !strings.HasPrefix(link, "//") {
					return
				}
				link = "http:" + link
			}
			urls = append(urls, link)
		}
	})
	return urls
}
