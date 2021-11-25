package page

import (
	"strings"
	"testing"

	"github.com/seggga/gb_best_go1/internal/domain"
	"github.com/stretchr/testify/assert"
)

var (
	// a URL to start test with
	startURL = "https://telegram.org"

	// application config
	cfg = domain.Config{
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
