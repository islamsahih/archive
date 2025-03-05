package fixer

import (
	"archive"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Options struct {
	InputDir string
}

type Fixer interface {
	Indexes(category string) error
	Save() error
	Tags() error
}

func New(opt Options) (Fixer, error) {
	cache, err := archive.LoadCache(filepath.Join(opt.InputDir, "cache"))
	if err != nil {
		return nil, err
	}
	cache.SetPath(filepath.Join(opt.InputDir, "cache.fixed"))
	return &fixer{
		cache: cache,
		opt:   opt,
	}, nil
}

type fixer struct {
	cache *archive.Cache
	opt   Options
}

func (f *fixer) Save() error {
	log.Printf("save cache: %s\n", f.cache.Path())
	return f.cache.Save()
}

func (f *fixer) Indexes(category string) error {
	items, ok := f.cache.GetCategory(category)
	if !ok || len(items) == 0 {
		return fmt.Errorf("fix indexes: unknown category: %s", category)
	}
	for i, item := range items {
		idx := i + 1
		item.Link = "/" + item.Category + "/" + strconv.Itoa(item.Index)
		if item.Index != idx {
			log.Printf("index: %d must be %d -----------------------------------------------\n", item.Index, idx)
			link := archive.CategoryItemLink(item.Category, idx)
			oldLink := item.Link

			if item, ok := f.cache.GetLink(link); ok {
				return fmt.Errorf("fix indexes: target link already exists: %s (index %d)", link, item.Index)
			}
			var text, audio string
			if item.Text != "" {
				text = strings.Replace(item.Text, oldLink, link, 1)
				if _, err := os.Stat(filepath.Join(f.opt.InputDir, text)); !os.IsNotExist(err) {
					return fmt.Errorf("fix indexes: target text already exists: %s", text)
				}
			}
			if item.Audio != "" {
				audio = strings.Replace(item.Audio, oldLink, link, 1)
				if _, err := os.Stat(filepath.Join(f.opt.InputDir, audio)); !os.IsNotExist(err) {
					return fmt.Errorf("fix indexes: target audio already exists: %s", audio)
				}
			}
			if refs, _ := f.cache.GetReferences(link); len(refs) != 0 {
				return fmt.Errorf("fix indexes: target references already exist: %s", link)
			}

			log.Printf("index: %d -> %d\n", item.Index, idx)
			item.Index = idx

			log.Printf("link: %s -> %s\n", oldLink, link)
			f.cache.FixLink(item, link)

			if item.Text != "" {
				if err := os.Rename(filepath.Join(f.opt.InputDir, item.Text), filepath.Join(f.opt.InputDir, text)); err != nil {
					return err
				}

				log.Printf("text: %s -> %s\n", item.Text, text)
				item.Text = text
			}
			if item.Audio != "" {
				if err := os.Rename(filepath.Join(f.opt.InputDir, item.Audio), filepath.Join(f.opt.InputDir, audio)); err != nil {
					return err
				}

				log.Printf("audio: %s -> %s\n", item.Audio, audio)
				item.Audio = audio
			}

			refItems, _ := f.cache.GetReferences(oldLink)
			for _, refItem := range refItems {
				f.cache.RemoveReference(refItem, oldLink)
				f.cache.AddReference(refItem, link)
			}
		}
	}

	log.Printf("category %s: index fixed", category)

	return nil
}

func (f *fixer) Tags() error {
	for _, category := range f.cache.ListCategories() {
		items, _ := f.cache.GetCategory(category)
		for _, item := range items {
			var links []*archive.Link
			for _, link := range item.UnhandledLinks {
				if link.URL == "messages.backup.html" {
					item.Tags = append(item.Tags, link.Title)
				} else {
					links = append(links, link)
				}
			}
			if len(links) == 0 {
				item.UnhandledLinks = nil
			} else {
				item.UnhandledLinks = links
			}
		}
	}
	return nil
}
