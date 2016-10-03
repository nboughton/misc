package bq

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
		doc, err := goquery.NewDocument(fmt.Sprintf("%v?q=%v&pg=%v", bQuoteURL, q, i))
		if err != nil {
			return quotes, err
		}

		doc.Find(".bqQt").Each(func(idx int, item *goquery.Selection) {
			qt := Quote{
				Text:   item.Find(".bqQuoteLink").Children().First().Text(),
				Author: item.Find(".bq-aut").Children().First().Text(),
			}
			if len(qt.Text) > 0 && len(qt.Author) > 0 {
				quotes = append(quotes, qt)
			}
		})

	}

	if len(quotes) > 0 {
		return quotes, nil
	}

	return quotes, fmt.Errorf("No quotes returned")
}
