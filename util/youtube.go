package util

import (
	"errors"
	"net/url"
	"path/filepath"
	"regexp"
)

var (
	youtubeRE = regexp.MustCompile(`(?:https?://)?(?:www.)?youtu\.?be(?:\.com)?\/?.*(?:watch|embed|shorts)?(?:.*v=|v\/|\/)([\w\-_]+)\&?`)
)

func IsYoutubeURL(href string) bool {
	return youtubeRE.MatchString(href)
}

func MakeYoutubeURL(href string) (string, error) {
	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	if filepath.Dir(u.Path) == "/shorts" {
		return "https://youtube.com/shorts/" + filepath.Base(u.Path), nil
	}
	var id string
	switch u.Host {
	case "youtube.com", "m.youtube.com", "www.youtube.com":
		id = u.Query().Get("v")
	case "youtu.be", "www.youtu.be":
		id = filepath.Base(u.Path)
	}
	if id == "" {
		return "", errors.New("invalid youtube url")
	}
	return "https://youtube.com/watch?v=" + id, nil
}
