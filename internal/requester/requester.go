package requester

import (
	"context"
	"net/http"
	"time"

	"github.com/seggga/gb_best_go1/internal/domain"
	"github.com/seggga/gb_best_go1/internal/page"
)

type requester struct {
	timeout time.Duration
	tran    http.RoundTripper
}

func NewRequester(timeout time.Duration, tran http.RoundTripper) requester {
	return requester{
		timeout: timeout,
		tran:    tran,
	}
}

// Get searches and returns a webpage on a given URL
func (r requester) Get(ctx context.Context, url string) (domain.Page, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		cl := &http.Client{
			Timeout:   r.timeout,
			Transport: r.tran,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		body, err := cl.Do(req)
		if err != nil {
			return nil, err
		}
		defer body.Body.Close()
		page, err := page.NewPage(body.Body)
		if err != nil {
			return nil, err
		}
		return page, nil
	}
	// return nil, nil     // unreachable code
}
