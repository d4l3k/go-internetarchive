package archive

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SearchPrefixURL is the search endpoint.
const SearchPrefixURL = "https://web.archive.org/cdx/search/cdx?matchType=prefix&output=json&showDupeCount=true&url=%s"

// SearchResult is a single search result.
type SearchResult struct {
	URLKey      string
	Timestamp   string
	OriginalURL string
	MIMEType    string
	StatusCode  string
	Digest      string
	Length      string
	DupCount    string
	Err         error
}

// SearchPrefix returns all results for a specific url prefix.
func SearchPrefix(url string) <-chan SearchResult {
	results := make(chan SearchResult, 1)

	finalURL := fmt.Sprintf(SearchPrefixURL, url)
	req, err := http.Get(finalURL)
	if err != nil {
		results <- SearchResult{Err: err}
		close(results)
		return results
	}

	d := json.NewDecoder(req.Body)

	go func() {
		defer req.Body.Close()
		defer close(results)

		if _, err := d.Token(); err != nil {
			results <- SearchResult{Err: err}
			return
		}
		for d.More() {
			var result []string
			if err := d.Decode(&result); err != nil {
				results <- SearchResult{Err: err}
				return
			}
			results <- SearchResult{
				URLKey:      result[0],
				Timestamp:   result[1],
				OriginalURL: result[2],
				MIMEType:    result[3],
				StatusCode:  result[4],
				Digest:      result[5],
				Length:      result[6],
				DupCount:    result[7],
			}
		}
	}()

	return results
}
