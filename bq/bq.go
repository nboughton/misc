package bq

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Quote contains the extracted quote data
type Quote struct {
	Author, Text string
}

// Search page URL for BrainyQuote.com
var bQuoteURL = "http://www.brainyquote.com/search_results.html"

// Search returns an array of quotes and can be constrained to only retrieve a set number
// of pages worth of results, if a page returns no quotes the function breaks out of the
// loop and returns all quotes found so far.
func Search(words []string, pages int) ([]Quote, error) {
	q, quotes := strings.Join(words, "+"), []Quote{}

	// get the search results
	for i := 1; i <= pages; i++ {
		qFound := 0
		resp, err := http.Get(fmt.Sprintf("%v?q=%v&pg=%v", bQuoteURL, q, i))
		if err != nil {
			return quotes, err
		}
		defer resp.Body.Close()

		doc, err := html.Parse(resp.Body)
		if err != nil {
			return quotes, err
		}

		var (
			getTextAndAuthor func(*html.Node)
			findQuotes       func(*html.Node)
			text, author     string
		)

		findQuotes = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "div" {
				for _, attr := range n.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "bqQt") {
						getTextAndAuthor(n)
						if len(text) > 0 && len(author) > 0 {
							quotes = append(quotes, Quote{Author: author, Text: text})
							text, author = "", ""
						}
						break
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				findQuotes(c)
			}
		}

		getTextAndAuthor = func(n *html.Node) {
			if n.Type == html.ElementNode && (n.Data == "span" || n.Data == "div") {
				for _, attr := range n.Attr {
					if attr.Key == "class" && attr.Val == "bqQuoteLink" {
						text = n.FirstChild.FirstChild.Data
					} else if attr.Key == "class" && attr.Val == "bq-aut" {
						author = n.FirstChild.FirstChild.Data
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				getTextAndAuthor(c)
			}
		}

		findQuotes(doc)
		if qFound == 0 {
			break
		}
	}

	if len(quotes) > 0 {
		return quotes, nil
	}

	return quotes, fmt.Errorf("No quotes returned")
}
