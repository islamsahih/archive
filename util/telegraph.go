package util

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
)

var (
	telegraphRE = regexp.MustCompile(`telegra\.ph\/[-a-zA-Z0-9.]+(\/\S*)?`)
)

func IsTelegraphURL(href string) bool {
	return telegraphRE.MatchString(href)
}

func FetchTelegraphArticle(href string) (string, error) {
	//fmt.Printf("saving telegraph article %s -> %s\n", href, path)
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
	article.Find("h1").Remove()
	article.Find("address").Remove()
	article.Find("figure").Remove()
	html, err := article.Html()
	if err != nil {
		return "", err
	}
	if len(html) == 0 {
		return "", errors.New("empty article")
	}
	return html, nil
}
