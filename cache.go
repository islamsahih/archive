package archive

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type Cache struct {
	path        string
	byID        map[string]*Item
	byLink      map[string]*Item
	byCategory  map[string][]*Item
	byReference map[string][]*Item
}

func LoadCache(path string) (*Cache, error) {
	cache := &Cache{
		path:        path,
		byID:        make(map[string]*Item),
		byLink:      make(map[string]*Item),
		byCategory:  make(map[string][]*Item),
		byReference: make(map[string][]*Item),
	}
	if f, err := os.OpenFile(path, os.O_RDONLY, 0660); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		defer f.Close()
		dec := json.NewDecoder(f)
		n := 0
		for {
			var item Item
			if err = dec.Decode(&item); err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("decode failed, decoded %d items", len(cache.byID))
				return nil, err
			}
			n++
			if err = cache.Add(&item); err != nil {
				log.Fatalf("read cache: %s", err)
			}
		}
		for _, items := range cache.byCategory {
			sort.Slice(items, func(i, j int) bool {
				return items[i].Index < items[j].Index
			})
			//for i, item := range items {
			//	if item.Index != i+1 {
			//		log.Printf("category %s: %s: corrupted index: %d\n", item.Category, item.ID, item.Index)
			//	}
			//}
		}
	}
	return cache, nil
}

func (cache *Cache) Save() error {
	f, err := os.OpenFile(cache.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	categories := make([]string, 0, len(cache.byCategory))
	for category := range cache.byCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	for _, category := range categories {
		items := cache.byCategory[category]
		for _, item := range items {
			if err = enc.Encode(&item); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cache *Cache) Size() int {
	return len(cache.byID)
}

func (cache *Cache) Path() string {
	return cache.path
}

func (cache *Cache) SetPath(path string) {
	cache.path = path
}

func (cache *Cache) GetID(id string) (*Item, bool) {
	item, ok := cache.byID[id]
	return item, ok
}

func (cache *Cache) GetLink(link string) (*Item, bool) {
	item, ok := cache.byLink[link]
	return item, ok
}

func (cache *Cache) FixLink(item *Item, link string) {
	delete(cache.byLink, item.Link)
	log.Printf("link %s deleted\n", item.Link)
	item.Link = link
	cache.byLink[link] = item
	log.Printf("link %s assigned to index %d\n", item.Link, item.Index)
}

func (cache *Cache) ListCategories() []string {
	cats := make([]string, 0, len(cache.byCategory))
	for cat := range cache.byCategory {
		cats = append(cats, cat)
	}
	return cats
}

func (cache *Cache) GetCategory(category string) ([]*Item, bool) {
	items, ok := cache.byCategory[category]
	return items, ok
}

func (cache *Cache) GetReferences(link string) ([]*Item, bool) {
	items, ok := cache.byReference[link]
	return items, ok
}

func (cache *Cache) AddReference(item *Item, link string) {
	for _, ref := range item.References {
		if ref == link {
			return
		}
	}
	item.References = append(item.References, link)
	cache.byReference[link] = append(cache.byReference[link], item)
}

func (cache *Cache) RemoveAllReferences(link string) {
	for {
		refItems, ok := cache.byReference[link]
		if !ok {
			return
		}
		if len(refItems) == 0 {
			delete(cache.byReference, link)
			return
		}
		cache.RemoveReference(refItems[0], link)
	}
}

func (cache *Cache) RemoveReference(src *Item, link string) {
	var items []*Item
	refItems, _ := cache.byReference[link]
	for _, item := range refItems {
		if item.ID != src.ID {
			items = append(items, item)
		} else {
			var refs []string
			for _, ref := range item.References {
				if ref != link {
					refs = append(refs, ref)
				}
			}
			item.References = refs
		}
	}
	if len(items) == 0 {
		delete(cache.byReference, link)
	} else {
		cache.byReference[link] = items
	}
}

func (cache *Cache) Add(item *Item) error {
	if err := item.Err(); err != nil {
		return fmt.Errorf("cache add: invalid item: %v\n", err)
	}
	if _, ok := cache.byID[item.ID]; ok {
		return fmt.Errorf("cache add: duplicated item ID: %s\n", item.ID)
	}
	if dup, ok := cache.byLink[item.Link]; ok {
		return fmt.Errorf("cache add: duplicated item link: %s (%s)\n", item.Link, dup.ID)
	}
	cache.byID[item.ID] = item
	cache.byLink[item.Link] = item
	cache.byCategory[item.Category] = append(cache.byCategory[item.Category], item)
	for _, ref := range item.References {
		cache.byReference[ref] = append(cache.byReference[ref], item)
	}
	return nil
}

func (cache *Cache) FillLinks() int {
	n := 0
	for _, item := range cache.byID {
		for _, links := range item.Links {
			for _, link := range links {
				if link.ID == "" {
					target, ok := cache.GetLink(link.URL)
					if !ok {
						log.Printf("fill links: %s: unknown link: %s", item.ID, link.URL)
						n++
					} else {
						link.ID = target.ID
					}
				}
			}
		}
	}
	return n
}

func (cache *Cache) ConvertSelfLinks() {
	for _, item := range cache.byID {
		delete(cache.byLink, item.Link)
		item.Link = item.InternalLink()
		cache.byLink[item.Link] = item
	}
}

func (cache *Cache) FillReferences() int {
	n := 0
	for _, item := range cache.byID {
		for _, links := range item.Links {
			for _, link := range links {
				if target, ok := cache.GetID(link.ID); !ok {
					log.Printf("fill references: %s: unknown id: %s", item.ID, link.ID)
					n++
				} else {
					item.References = append(item.References, target.Link)
					cache.byReference[target.Link] = append(cache.byReference[target.Link], item)
				}
			}
		}
		item.Links = nil
	}
	return n
}

func (cache *Cache) CheckReferences() int {
	var n int
	for _, item := range cache.byID {
	ItemLinksLoop:
		for _, link := range item.References {
			if target, ok := cache.byLink[link]; !ok {
				log.Printf("check references: %s: unknown link: %s", item.ID, link)
				n++
			} else {
				for _, link := range target.References {
					if link == item.Link {
						continue ItemLinksLoop
					}
				}
				log.Printf("check references: %s: back link not exists: %s", target.ID, item.Link)
				n++
			}
		}
	}
	return n
}
