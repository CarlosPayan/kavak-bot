// internal/utils/fetch_kavak.go
package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchKavakInfo(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading info from Kavak: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status no OK: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing HTML: %w", err)
	}

	collect := func(selection *goquery.Selection) string {
		var buf bytes.Buffer
		selection.Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				buf.WriteString(text + "\n")
			}
		})
		return buf.String()
	}

	var buf bytes.Buffer
	doc.Find("article p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			buf.WriteString(text + "\n")
		}
	})

	content := collect(doc.Find("article p"))

	if content == "" {
		content = collect(doc.Find("div.post-content p"))
	}

	if content == "" {
		content = collect(doc.Find("div.blog-post-content p"))
	}

	if content == "" {
		content = collect(doc.Find("p"))
	}

	if content == "" {
		return "", fmt.Errorf("couldn't extract info from page: no se encontraron <p> válidos")
	}

	return content, nil
}
