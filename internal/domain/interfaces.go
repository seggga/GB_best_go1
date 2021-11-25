package domain

import (
	"context"
)

// Config is a structure to setup the application
type Config struct {
	MaxDepth     uint64 `yaml:"maxdepth"`
	MaxResults   int    `yaml:"maxresults"`
	MaxErrors    int    `yaml:"maxerrors"`
	URL          string `yaml:"url"`
	ReqTimeout   int    `yaml:"reqtimeout"`
	CrawlTimeout int    `yaml:"crawltimeout"`
}

// CrawlResult is a structure that represents certain status on given page
type CrawlResult struct {
	Err   error
	Title string
	Url   string
}

// Page represents a parsed web-page
type Page interface {
	GetTitle() string   // obtain title of the 'page'
	GetLinks() []string // collects a list of links found on the given 'page'
}

// Requester sends queries to obtain Pages
type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

// Crawler uses Requester to collect Pages
type Crawler interface {
	Scan(ctx context.Context, url string, depth uint64)
	ChanResult() <-chan CrawlResult
	IncreaseDepth()
}
