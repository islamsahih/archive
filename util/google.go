package util

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

var (
	googleDocsRE  = regexp.MustCompile(`https://docs\.google\.com/document/d/.*/`)
	googleDriveRE = regexp.MustCompile(`https://drive\.google\.com/file/d/.*/`)
)

func IsGoogleDocsURL(href string) bool {
	return googleDocsRE.MatchString(href)
}

func IsGoogleDriveURL(href string) bool {
	return googleDriveRE.MatchString(href)
}

func fetchGoogleDocsArticle(href string) (string, error) {
	href += "export?exportFormat=html"
	//fmt.Printf("saving google docs article %s -> %s\n", href, path)
	r, err := http.Get(href)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return "", err
	}
	article := doc.Find("body")
	html, err := article.Html()
	if err != nil {
		return "", err
	}
	if len(html) == 0 {
		return "", errors.New("empty article")
	}
	return html, nil
}

func FetchGoogleDocsArticle(href string) (string, error) {
	return fetchGoogleDocsArticle(googleDocsRE.FindString(href))
}

func FetchGoogleDriveArticle(href string) (string, error) {
	return fetchGoogleDocsArticle(
		strings.Replace(googleDriveRE.FindString(href), "drive.google.com/file", "docs.google.com/document", 1),
	)
}
