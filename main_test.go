package main

import (
	"context"
	"log"
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

func TestCrawlerScan(t *testing.T) {
	want := []string{
		"",
		"",
	}

	requester := NewRequester(time.Second)
	crawler := NewCrawler(requester, 1)
	ctx, cancel := context.WithCancel(context.Background())
	url := "http://localhost:8080/homepage"
	go crawler.Scan(ctx, url, 1)
	got := mockProcessResult(ctx, cancel, crawler)

	if len(want) != len(got) {
		t.Errorf("slice with different length: got %d, want %d", len(want), len(got))
	}

	//rawler := NewCrawler()
	t.Skip()
}

func mockProcessResult(ctx context.Context, cancel func(), cr Crawler) []string {
	// var maxResult, maxErrors = cfg.MaxResults, cfg.MaxErrors
	for {
		select {
		case <-ctx.Done():
			return nil

		// got message in the channel
		case msg := <-cr.ChanResult():
			if msg.Err != nil {
				// message contains error
				maxErrors--
				log.Printf("crawler result return err: %s\n", msg.Err.Error())
				if maxErrors <= 0 {
					log.Println("Maximum number of errors occured.")
					cancel()
					return nil
				}
			} else {
				// message contains data
				maxResult--
				log.Printf("crawler result: [url: %s] Title: %s\n", msg.Url, msg.Title)
				if maxResult <= 0 {
					log.Println("Maximum number of results obtained.")
					cancel()
					return nil
				}
			}
		}
	}
}
