package util

import "fmt"

func Fetch(src string) (string, error) {
	if IsTelegraphURL(src) {
		return FetchTelegraphArticle(src)
	} else if IsTeletypeURL(src) {
		return FetchTeletypeArticle(src)
	} else if IsGoogleDocsURL(src) {
		return FetchGoogleDocsArticle(src)
	} else if IsGoogleDriveURL(src) {
		return FetchGoogleDriveArticle(src)
	}
	return "", fmt.Errorf("invalid url: %s", src)
}
