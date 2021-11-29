package requester

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/seggga/gb_best_go1/internal/domain"
	"github.com/seggga/gb_best_go1/internal/page"
)

type requester struct {
	timeout time.Duration
	tran    http.RoundTripper
}

func NewRequester(timeout int, tran http.RoundTripper) (*requester, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("cannot create requester, timeout is not valid, %d", timeout)
	}

	return &requester{
		timeout: time.Duration(timeout) * time.Second,
		tran:    tran,
	}, nil
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
		defer body.Body.Close() //nolint:errcheck // error don't have influence on business logic
		page, err := page.NewPage(body.Body)
		if err != nil {
			return nil, err
		}
		return page, nil
	}
	// return nil, nil     // unreachable code
}
