//nolint:errcheck
package requester

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// a URL to start test with
	startURL = "https://telegram.org"

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

// Requester interface Get()
func TestRequesterGet(t *testing.T) {
	// requester uses test http.Client with RoundTrip function
	requester, _ := NewRequester(1, roundTripFunc(func(r *http.Request) (*http.Response, error) {
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
