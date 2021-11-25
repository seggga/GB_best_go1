package crawler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/seggga/gb_best_go1/internal/domain"
	"github.com/seggga/gb_best_go1/internal/requester"
	"github.com/stretchr/testify/assert"
)

var (
	// a URL to start test with
	startURL = "https://telegram.org"

	// application config
	cfg = domain.Config{
		MaxDepth:     2,
		URL:          startURL,
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
	requester, _ := requester.NewRequester(1, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(testWebPage)),
		}, nil
	}))

	crawler, _ := NewCrawler(requester, cfg.MaxDepth)
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

/////////////////////////
//
// mocked test
//
/////////////////////////

type mockPage struct{}

func (mp mockPage) GetTitle() string {
	return "TestDocument"
}

func (mp mockPage) GetLinks() []string {
	return []string{
		"http://yandex.com/two",
		"http://google.com/one",
		"http://yahoo.com/three",
		"http://rambler.com/four",
		"http://bing.com/five",
	}
}

type mockRequester struct{}

func (mock *mockRequester) Get(ctx context.Context, url string) (domain.Page, error) {
	var mp mockPage
	return mp, nil
}

// Crawler interface Scan()
func TestMockCrawlerScan(t *testing.T) {
	want := []string{
		"url: " + startURL + ", title: TestDocument",
		"url: http://google.com/one, title: TestDocument",
		"url: http://yandex.com/two, title: TestDocument",
		"url: http://yahoo.com/three, title: TestDocument",
		"url: http://rambler.com/four, title: TestDocument",
		"url: http://bing.com/five, title: TestDocument",
	}
	// requester uses test http.Client with RoundTrip function
	requester := new(mockRequester)
	crawler, _ := NewCrawler(requester, cfg.MaxDepth)
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
