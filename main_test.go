package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
		"",
		"",
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

	crawler := NewCrawler(requester, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go crawler.Scan(ctx, "http://google.com", 1)
	go watchDog(cancel, 30*time.Second, t)
	got := mockProcessResult(ctx, cancel, crawler)

	t.Logf("%+v", got)
	if len(want) != len(got) {
		t.Errorf("slice with different length: got %d, want %d", len(got), len(want))
	}

}

func mockProcessResult(ctx context.Context, cancel func(), cr Crawler) []string {
	// var maxResult, maxErrors = cfg.MaxResults, cfg.MaxErrors
	var res []string
	for {
		select {
		case <-ctx.Done():
			// context has been closed
			return res

		case msg := <-cr.ChanResult():
			// got message in the channel
			res = append(res, fmt.Sprintf("url: %s, title: %s", msg.Url, msg.Title))
		}
	}
}

func watchDog(cancel context.CancelFunc, dur time.Duration, t *testing.T) {

	sigInt := make(chan os.Signal)        //Создаем канал для приема сигналов
	signal.Notify(sigInt, syscall.SIGINT) //Подписываемся на сигнал SIGINT

	select {
	case <-time.After(dur):
		t.Log("context has been closed on timeout")
		cancel()
	case <-sigInt:
		t.Log("context has been closed on interrupt signal")
		cancel() //Если пришёл сигнал SigInt - завершаем контекст
	}
}
