package converter

import (
	"archive"
	"log"
	"path/filepath"
)

type Options struct {
	InputDir  string
	OutputDir string
}

type Converter interface {
	Convert() error
}

func New(opt Options) (Converter, error) {
	cache, err := archive.LoadCache(filepath.Join(opt.InputDir, "cache"))
	if err != nil {
		return nil, err
	}
	cache.SetPath(filepath.Join(opt.OutputDir, "cache"))
	return &converter{
		cache: cache,
		opt:   opt,
	}, nil
}

type converter struct {
	cache *archive.Cache
	opt   Options
}

func (c *converter) Convert() error {
	if n := c.cache.FillLinks(); n != 0 {
		log.Printf("fill links: %d unknown links", n)
	}
	c.cache.ConvertSelfLinks()
	if n := c.cache.FillReferences(); n != 0 {
		log.Printf("fill references: %d unknown ids", n)
	}
	if n := c.cache.CheckReferences(); n != 0 {
		log.Printf("check references: %d invalid links", n)
	}
	return c.cache.Save()
}
