package builder

import (
	"archive"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	InputDir  string
	OutputDir string
	Format    string
	Category  string
	Index     int
}

type Builder interface {
	Build() error
}

type builder struct {
	cache *archive.Cache
	opt   Options
}

func New(opt Options) (Builder, error) {
	cache, err := archive.LoadCache(filepath.Join(opt.InputDir, "cache"))
	if err != nil {
		return nil, err
	}
	return &builder{
		cache: cache,
		opt:   opt,
	}, nil
}

func (b *builder) Build() error {
	if err := os.MkdirAll(b.opt.OutputDir, 0775); err != nil {
		return err
	}
	if b.opt.Category != "" && b.opt.Index != 0 {
		return b.buildCategoryItem(b.opt.Category, b.opt.Index)
	} else {
		for _, category := range b.cache.ListCategories() {
			if err := b.buildCategory(category); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *builder) buildCategory(category string) error {
	items, ok := b.cache.GetCategory(category)
	if !ok {
		return fmt.Errorf("category not found: %s", category)
	}
	if err := os.MkdirAll(filepath.Join(b.opt.OutputDir, category), 0775); err != nil {
		return err
	}
	for _, item := range items {
		if err := b.buildItem(item); err != nil {
			return err
		}
	}
	return nil
}

func (b *builder) buildCategoryItem(category string, index int) error {
	items, ok := b.cache.GetCategory(category)
	if !ok {
		return fmt.Errorf("category not found: %s", category)
	}
	if err := os.MkdirAll(filepath.Join(b.opt.OutputDir, category), 0775); err != nil {
		return err
	}
	for _, item := range items {
		if item.Index == index {
			return b.buildItem(item)
		}
	}
	return fmt.Errorf("item %d not found in category %s", index, category)
}

func (b *builder) buildItem(item *archive.Item) error {
	if b.opt.Format == "html" {
		return b.buildItemHTML(item)
	}
	if b.opt.Format == "md" {
		return b.buildItemMD(item)
	}
	return fmt.Errorf("invalid format: %s", b.opt.Format)
}

func (b *builder) buildItemHTML(item *archive.Item) error {
	link := item.Link
	item.ID = ""
	item.Category = ""
	item.Link = ""
	if item.Text != "" {
		raw, err := os.ReadFile(filepath.Join(b.opt.InputDir, item.Text))
		if err != nil {
			return err
		}
		item.Body = string(raw)
		item.Text = ""
	}
	f, err := os.OpenFile(filepath.Join(b.opt.OutputDir, link+".json"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "    ")
	e.SetEscapeHTML(false)
	return e.Encode(item)
}

func (b *builder) buildItemMD(item *archive.Item) error {
	f, err := os.OpenFile(filepath.Join(b.opt.OutputDir, item.Link+".md"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "---\n")
	if err = yaml.NewEncoder(f).Encode(item); err != nil {
		return err
	}
	fmt.Fprintf(f, "---\n")
	if item.Text != "" {
		cmd := exec.Command("pandoc", filepath.Join(b.opt.InputDir, item.Text), "-f", "html", "-t", "markdown", "-t", "markdown-smart")
		cmd.Stdout = f
		return cmd.Run()
	}
	return nil
}
