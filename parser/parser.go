package parser

import (
	"archive"
	"archive/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	cp "github.com/otiai10/copy"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	Category       string
	InputDir       string
	OutputDir      string
	AudioOutputDir string
	TextOutputDir  string
	SkipIndex      int
	SkipItems      string
	VideoMatch     []string
	TextMatch      []string
	TagMatch       []string
	SkipLinkMatch  []string
	LinkTypeMatch  map[string][]string
	SelfLinkMatch  []string
	SaveText       bool
	AudioRequired  bool
	TitleRequired  bool
	ProcessAudio   bool
	ProcessVideo   bool
	ProcessText    bool
	NoAlerts       bool
	NoCheck        bool
}

type Parser interface {
	Parse() error
}

const (
	itemSelector   = ".message.default"
	dateSelector   = ".message.service > .body.details"
	dateClass      = "body details"
	idAttribute    = "id"
	dateFormat     = "2006-01-02"
	audioSelector  = ".media_wrap > a.media_audio_file"
	audioAttribute = "href"
	titleSelector  = ".title"
	textSelector   = ".text"
	linkSelector   = "a"
)

func New(r io.Reader, w io.Writer, opt Options) (Parser, error) {
	cache, err := archive.LoadCache(filepath.Join(opt.OutputDir, "cache"))
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	skips := make(map[string]struct{})
	if opt.SkipItems != "" {
		f, err := os.OpenFile(opt.SkipItems, os.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		for {
			var s string
			if _, err := fmt.Fscanln(f, &s); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			skips[s] = struct{}{}
		}
	}
	return &parser{
		cache: cache,
		doc:   doc,
		out:   w,
		skips: skips,
		opt:   opt,
	}, err
}

type parser struct {
	cache *archive.Cache
	doc   *goquery.Document
	out   io.Writer
	skips map[string]struct{}
	opt   Options
}

func (p *parser) Parse() error {
	defer p.cache.Save()

	p.write("[")
	defer p.write("]")

	n := p.opt.SkipIndex
	ts := int64(0)
	p.doc.Find(itemSelector + "," + dateSelector).Each(func(_ int, s *goquery.Selection) {
		if s.HasClass(dateClass) {
			if date, err := dateparse.ParseAny(strings.TrimSpace(s.Text())); err == nil {
				ts = date.Unix()
			}
			return
		}

		var item archive.Item
		item.Links = make(map[string][]*archive.Link)

		// general

		id, _ := s.Attr(idAttribute)
		if id == "" {
			//log.Println("skipping record without id")
			return
		}
		if _, ok := p.skips[id]; ok {
			//log.Printf("skipping id %s\n", id)
			return
		}
		n++
		item.ID = id //uuid.New().String()
		item.Category = p.opt.Category
		item.Index = n
		item.Date = time.Unix(ts, 0).Format(dateFormat)
		item.Title = strings.TrimSpace(s.Find(titleSelector).First().Text())
		if item.Title == "" || item.Title == "Audio file" || item.Title == "Voice message" || strings.Contains(item.Title, "…") {
			text := s.Find(textSelector).Clone()
			//text.Find("a").Remove()
			title := strings.TrimSpace(text.Text())
			item.Title = title
		}
		if item.Title == "" {
			if p.opt.TitleRequired {
				log.Printf("ERROR >> %s >> skipping record without title\n", id)
				n--
				return
			}
			//item.Title = item.Category + "/" + strconv.Itoa(item.Index) + "##" + id
			item.Title = "Без названия"
		}
		item.Origin = p.opt.SelfLinkMatch[0] + "/" + strings.TrimPrefix(id, "message")
		item.Link = item.Category + "/" + strconv.Itoa(item.Index)

		// audio

		audio := s.Find(audioSelector).First().AttrOr(audioAttribute, "")
		if audio == "" {
			if p.opt.AudioRequired {
				log.Printf("ERROR >> %s >> skipping record without audio\n", id)
				n--
				return
			}
		} else {
			p.processAudio(&item, audio)
		}

		// parse text
		text := s.Find(textSelector)
		if p.opt.SaveText {
			p.saveText(&item, text)
		}
		text.Find(linkSelector).Each(func(i int, s *goquery.Selection) {
			if href, _ := s.Attr("href"); href != "" {
				if strings.HasPrefix(href, "stickers") {
					return
				}
				if strings.Contains(href, item.Link) {
					return
				}
				link := archive.Link{
					Title: s.Text(),
					URL:   href,
				}
				//if p.processSelfLink(&item, &link) {
				//	return
				//}
				if p.opt.SaveText {
					if p.processText(&item, &link, true) {
						return
					}
				} else if item.Text == "" {
					if p.processText(&item, &link, false) {
						return
					}
				}
				if p.processVideo(&item, &link) {
					return
				}
				if p.processLinkType(&item, &link) {
					return
				}
				item.UnhandledLinks = append(item.UnhandledLinks, &link)
			} else if onclick, _ := s.Attr("onclick"); onclick != "" {
				link := &archive.Link{
					Title:   s.Text(),
					OnClick: onclick,
				}
				if p.processTag(&item, link) {
					return
				}
			}
		})

		if item.Audio == "" && item.Video == "" && item.Text == "" {
			//log.Printf("ERROR >> %s >> skipping record without content: %s\n", id, item.Title)
			n--
			return
		}

		//log.Printf("%s >> %s\n", id, item.Title)

		if !p.opt.NoCheck {
			if err := item.Err(); err != nil {
				log.Printf("message %s: invalid item: %v\n", id, err)
				if err = p.clear(&item); err != nil {
					log.Fatal(err)
				}
				n--
				return
			}
		}

		if !p.opt.NoAlerts {
			for _, link := range item.UnhandledLinks {
				log.Printf("item %s (%d): unhandled link: %s\n", item.ID, item.Index, link.URL)
			}
		}

		if err := p.cache.Add(&item); err != nil {
			log.Fatalf("ERROR >> %s >> %s >> item %s: %s", id, item.ID, item.Title, err)
		}

		if n != p.opt.SkipIndex+1 {
			p.write(",\n")
		}
		p.write(&item)
		return
	})

	fmt.Printf("done, parsed %d items, last index: %d\n", p.cache.Size(), n)
	return nil
}

func (p *parser) clear(item *archive.Item) error {
	if item.Text != "" {
		if err := os.Remove(filepath.Join(p.opt.OutputDir, item.Text)); err != nil {
			return err
		}
	}
	if item.Audio != "" {
		if err := os.Remove(filepath.Join(p.opt.OutputDir, item.Audio)); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) saveText(item *archive.Item, s *goquery.Selection) {
	path := filepath.Join(p.opt.TextOutputDir, p.opt.Category, strconv.Itoa(item.Index)+".html")
	s = s.Clone()
	s.Find("a").Remove()
	html, _ := s.Html()
	html = strings.TrimSpace(html)
	if err := util.SaveHTML(html, path, false); err != nil {
		log.Fatalf("item %s (%d): save text failed: %v\n", item.ID, item.Index, err)
	}
	item.Text, _ = filepath.Rel(p.opt.OutputDir, path)
}

func (p *parser) processAudio(item *archive.Item, href string) {
	do := func(href string) error {
		item.Audio = href
		if p.opt.ProcessAudio && item.Audio != "" {
			file := strconv.Itoa(item.Index) + ".mp3"
			path := filepath.Join(p.opt.AudioOutputDir, p.opt.Category, file)
			if err := cp.Copy(filepath.Join(p.opt.InputDir, item.Audio), path); err != nil {
				return err
			}
			item.Audio, _ = filepath.Rel(p.opt.OutputDir, path)
		}
		return nil
	}

	err := do(href)
	if err != nil {
		href, err = url.QueryUnescape(href)
		if err == nil {
			err = do(href)
		}
	}
	if err != nil {
		log.Fatalf("item %s (%d): processing audio: %v\n", item.ID, item.Index, err)
	}
}

func (p *parser) processSelfLink(item *archive.Item, link *archive.Link) bool {
	if p.match(link, p.opt.SelfLinkMatch) {
		path := link.URL
		if u, err := url.Parse(link.URL); err == nil {
			path = strings.TrimPrefix(u.Path, "/")
		}
		if item.Link == "" {
			item.Link = path
			return true
		}
		if item.Link == path {
			return true
		}
	}
	return false
}

func (p *parser) processText(item *archive.Item, link *archive.Link, append bool) bool {
	if p.match(link, p.opt.TextMatch) {
		if p.opt.ProcessText {
			path := filepath.Join(p.opt.TextOutputDir, p.opt.Category, strconv.Itoa(item.Index)+".html")
			html, err := util.Fetch(link.URL)
			if err != nil {
				log.Fatalf("item %s (%d): processing text: %v\n", item.ID, item.Index, err)
			}
			if err = util.SaveHTML(html, path, append); err != nil {
				log.Fatalf("item %s (%d): save text failed: %v\n", item.ID, item.Index, err)
			}
			item.Text, _ = filepath.Rel(p.opt.OutputDir, path)
		} else if item.Text == "" {
			item.Text = link.URL
		}
		return true
	}
	return false
}

func (p *parser) processVideo(item *archive.Item, link *archive.Link) bool {
	if p.match(link, p.opt.VideoMatch) {
		var url string
		var err error
		if util.IsYoutubeURL(link.URL) {
			url, err = util.MakeYoutubeURL(link.URL)
		} else if util.IsTelegramURL(link.URL) {
			if p.opt.ProcessVideo {
				url, err = util.GetYoutubeURLFromTelegramPost(link.URL)
			} else {
				url = link.URL
			}
		} else {
			err = errors.New("invalid video url")
		}
		if err != nil {
			log.Printf("item %s (%d): processing video: %v: %s\n", item.ID, item.Index, link.URL, err)
			return false
		}
		if item.Video != "" && item.Video != url {
			return false
		}
		item.Video = url
		return true
	}
	return false
}

func (p *parser) processTag(item *archive.Item, link *archive.Link) bool {
	if p.match(link, p.opt.TagMatch) {
		item.Tags = append(item.Tags, link.Title)
		return true
	}
	return false
}

func (p *parser) processLinkType(item *archive.Item, link *archive.Link) bool {
	for t, m := range p.opt.LinkTypeMatch {
		if p.match(link, m) {
			if u, err := url.Parse(link.URL); err == nil {
				link.URL = strings.TrimPrefix(u.Path, "/")
			}
			item.Links[t] = append(item.Links[t], link)
			return true
		}
	}
	if p.match(link, p.opt.SkipLinkMatch) {
		return true
	}
	return false
}

func (p *parser) match(link *archive.Link, matcher []string) bool {
	for _, m := range matcher {
		if strings.Contains(link.Title, m) || strings.Contains(link.URL, m) || strings.Contains(link.OnClick, m) {
			return true
		}
	}
	return false
}

func (p *parser) write(data any) {
	var b []byte
	switch s := data.(type) {
	case []byte:
		b = s
	case string:
		b = []byte(s)
	default:
		var err error
		if b, err = json.MarshalIndent(data, "", "  "); err != nil {
			log.Fatalln(err)
		}
	}
	if _, err := p.out.Write(b); err != nil {
		log.Fatalln(err)
	}
}
