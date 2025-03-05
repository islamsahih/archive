package util

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var (
	telegramRE = regexp.MustCompile(`t\.me\/[-a-zA-Z0-9.]+(\/\S*)?`)
)

func IsTelegramURL(href string) bool {
	return telegramRE.MatchString(href)
}

func GetYoutubeURLFromTelegramPost(href string) (string, error) {
	r, err := http.Get(href)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	if url := youtubeRE.Find(body); len(url) != 0 {
		return MakeYoutubeURL(string(url))
	}
	return "", fmt.Errorf("no youtube post found: %s", href)
}
