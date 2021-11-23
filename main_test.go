package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	// a URL to start test with
	startURL = "https://telegram.org"

	// application config
	cfg = Config{
		MaxDepth:     2,
		MaxResults:   20,
		MaxErrors:    20,
		Url:          startURL,
		ReqTimeout:   5,
		CrawlTimeout: 5,
	}

	// test webpage to parse and use in http.RoundTripper
	testWebPage = `<!DOCTYPE html>
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
)

// Page interface, GetTitle()
func TestPageGetTitle(t *testing.T) {
	testPage, _ := NewPage(strings.NewReader(testWebPage))
	got := testPage.GetTitle()
	want := "TestDocument"
	if got != want {
		t.Errorf("titles not equal: got %s, want %s", got, want)
	}
	t.Log("page.GetTitle() - OK ")
}

// Page interface, GetLinks()
func TestPageGetLinks(t *testing.T) {
	testPage, _ := NewPage(strings.NewReader(testWebPage))
	got := testPage.GetLinks()
	want := []string{
		"http://yandex.com/two",
		"http://google.com/one",
		"http://yahoo.com/three",
		"http://rambler.com/four",
		"http://bing.com/five",
	}
	assert.ElementsMatch(t, want, got)
	t.Log("page.GetLinks() - OK ")
}

// Requester interface Get()
func TestRequesterGet(t *testing.T) {
	// requester uses test http.Client with RoundTrip function
	requester := NewRequester(time.Second, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(testWebPage)),
		}, nil
	}))
	got, _ := requester.Get(context.Background(), startURL)

	// check Titles
	wantTitle := "TestDocument"
	if got.GetTitle() != wantTitle {
		t.Errorf("page mismatch: titles are not equal. want %s, got %s", wantTitle, got.GetTitle())
	}
	t.Log("Titles are equal")

	// check URLs
	wantLinks := []string{
		"http://yandex.com/two",
		"http://google.com/one",
		"http://yahoo.com/three",
		"http://rambler.com/four",
		"http://bing.com/five",
	}
	assert.ElementsMatch(t, wantLinks, got.GetLinks())
	t.Log("URLs are equal")
}

// describe RoundTripper interface to pass into http.Client
type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

// Crawler interface Scan()
func TestCrawlerScan(t *testing.T) {
	want := []string{
		"url: " + startURL + ", title: TestDocument",
		"url: http://google.com/one, title: TestDocument",
		"url: http://yandex.com/two, title: TestDocument",
		"url: http://yahoo.com/three, title: TestDocument",
		"url: http://rambler.com/four, title: TestDocument",
		"url: http://bing.com/five, title: TestDocument",
	}
	// requester uses test http.Client with RoundTrip function
	requester := NewRequester(time.Second, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(testWebPage)),
		}, nil
	}))

	crawler := NewCrawler(requester, cfg)
	ctx := context.Background()

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

	assert.ElementsMatch(t, want, res)
	t.Log("crawler.Scan() - OK ")
}
