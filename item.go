package archive

import (
	"errors"
	"strconv"
)

type Link struct {
	ID      string `json:"id,omitempty"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	OnClick string `json:"onclick,omitempty"`
}

type Item struct {
	ID             string             `json:"id,omitempty"              yaml:"-"`
	Index          int                `json:"index,omitempty"           yaml:"index"`
	Category       string             `json:"category,omitempty"        yaml:"-"`
	DirIndex       int                `json:"dir_index,omitempty"       yaml:"-"`
	Date           string             `json:"date,omitempty"            yaml:"date"`
	Title          string             `json:"title,omitempty"           yaml:"title"`
	Description    string             `json:"description,omitempty"     yaml:"description"`
	Icon           string             `json:"icon,omitempty"            yaml:"icon"`
	Link           string             `json:"link,omitempty"            yaml:"-"`
	Origin         string             `json:"origin,omitempty"          yaml:"-"`
	Text           string             `json:"text,omitempty"            yaml:"-"`
	Body           string             `json:"body,omitempty"            yaml:"-"`
	Audio          string             `json:"audio,omitempty"           yaml:"audio"`
	Video          string             `json:"video,omitempty"           yaml:"video"`
	Tags           []string           `json:"tags,omitempty"            yaml:"tags,omitempty"`
	Links          map[string][]*Link `json:"links,omitempty"           yaml:"-"`
	UnhandledLinks []*Link            `json:"unhandled_links,omitempty" yaml:"-"`
	References     []string           `json:"references,omitempty"      yaml:"references,omitempty"`
	Meta           map[string]any     `json:"meta,omitempty"            yaml:"meta,omitempty"`
	Settings       map[string]any     `json:"settings,omitempty"        yaml:"settings,omitempty"`
}

func (item *Item) Err() error {
	if item.ID == "" {
		return errors.New("no ID")
	}
	if item.Link == "" {
		return errors.New("no link")
	}
	if item.Category == "" {
		return errors.New("uncategorized")
	}
	if item.Index == 0 {
		return errors.New("no category index")
	}
	return nil
}

func (item *Item) InternalLink() string {
	return CategoryItemLink(item.Category, item.Index)
}

func CategoryItemLink(category string, index int) string {
	return "/" + category + "/" + strconv.Itoa(index)
}
