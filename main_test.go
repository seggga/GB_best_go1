package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestPageGetTitle(t *testing.T) {
	t.Skip()
}

func TestPageGetLinks(t *testing.T) {
	t.Skip()
}
func TestRequesterGet(t *testing.T) {
	t.Skip()
}

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func TestCrawlerScan(t *testing.T) {

	want := []string{
		"url: https://telegram.org Title: TestDocument",
		"url: http://google.com/one Title: TestDocument",
		"url: http://yandex.com/two Title: TestDocument",
		"url: http://yahoo.com/three Title: TestDocument",
		"url: http://rambler.com/four Title: TestDocument",
		"url: http://bing.com/five Title: TestDocument",
	}

	requester := NewRequester(time.Second, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		testWebPage := `<!DOCTYPE html>
						<html lang="en">
							<head>
								<meta charset="UTF-8">
								<meta name="viewport" content="width=device-width, initial-scale=1.0">
								<title>TestDocument</title>
							</head>
							<body>
								<p><a href="http://google.com/one">ONE</a></p>
								<p><a href="http://yandex.com/two">TWO</a></p>
								<p><a href="http://yahoo.com/three">THREE</a></p>
								<p><a href="http://rambler.com/four">FOUR</a></p>
								<p><a href="http://bing.com/five">FIVE</a></p>
							</body>
						</html>`

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(testWebPage)),
		}, nil
	}))
	startURL := "https://telegram.org"
	cfg := Config{
		MaxDepth:     2,
		MaxResults:   20,
		MaxErrors:    20,
		Url:          startURL,
		ReqTimeout:   5,
		CrawlTimeout: 5,
	}
	crawler := NewCrawler(requester, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go crawler.Scan(ctx, startURL, 1)
	var res []string
	var next = true
	for next {
		select {
		case <-time.After(time.Duration(cfg.CrawlTimeout) * time.Second):
			// stop test on timeout
			t.Logf("read results stopped on timeout %d sec", cfg.CrawlTimeout)
			next = false
		case msg := <-crawler.ChanResult():
			// got result in the channel
			fmt.Printf("%+v", msg)
			res = append(res, fmt.Sprintf("url: %s, title: %s", msg.Url, msg.Title))
		}
	}

	t.Logf("%+v", res)
	if len(want) != len(res) {
		t.Errorf("slice with different length: got %d, want %d", len(res), len(want))
	}
}
