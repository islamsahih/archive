package util

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
)

var (
	teletypeRE = regexp.MustCompile(`teletype\.in/.*$`)
)

func IsTeletypeURL(href string) bool {
	return teletypeRE.MatchString(href)
}

func FetchTeletypeArticle(href string) (string, error) {
	r, err := http.Get(href)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return "", err
	}
	article := doc.Find("article")
	html, err := article.Html()
	if err != nil {
		return "", err
	}
	if len(html) == 0 {
		return "", errors.New("empty article")
	}
	return html, nil
}
