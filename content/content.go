package content

import (
	"archive"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Options struct {
	Pack          bool
	Unpack        bool
	Repack        string
	ItemDir       string
	ItemFile      string
	TextDir       string
	TextFile      string
	FieldsDir     string
	FieldsFile    string
	NumFirst      bool
	NumOnly       bool
	Skip          int
	TitleTemplate string
}

type Operation interface {
	Exec() error
}

type operation struct {
	opt Options
}

func NewOperation(opt Options) (Operation, error) {
	nact := 0
	if opt.Pack {
		nact++
	}
	if opt.Unpack {
		nact++
	}
	if opt.Repack != "" {
		nact++
	}
	if nact != 1 {
		return nil, errors.New("invalid options: operation")
	}
	if (opt.ItemFile != "" && opt.ItemDir != "") || (opt.ItemFile == "" && opt.ItemDir == "") {
		return nil, errors.New("invalid options: item")
	}
	if opt.FieldsDir == "" {
		opt.FieldsDir = filepath.Dir(opt.FieldsFile)
	}
	if opt.TextDir == "" {
		opt.TextDir = filepath.Dir(opt.TextFile)
	}
	return &operation{
		opt: opt,
	}, nil
}

func (a *operation) Exec() error {
	if err := os.MkdirAll(a.opt.FieldsDir, 0775); err != nil {
		return err
	}
	if err := os.MkdirAll(a.opt.TextDir, 0775); err != nil {
		return err
	}

	if a.opt.Pack {
		return a.pack()
	}
	if a.opt.Unpack {
		return a.unpack()
	}
	if a.opt.Repack != "" {
		return a.repack()
	}
	return errors.New("invalid operation")
}

func (a *operation) pack() error {
	if a.opt.ItemFile != "" {
		fieldsFile := a.opt.FieldsFile
		if fieldsFile == "" {
			fieldsFile = a.fieldsFilePath(a.opt.ItemFile)
		}
		textFile := a.opt.TextFile
		if textFile == "" {
			textFile = a.textFilePath(a.opt.ItemFile)
		}
		return packFile(a.opt.ItemFile, fieldsFile, textFile)
	} else if a.opt.ItemDir != "" {
		files, err := a.itemDirFiles()
		if err != nil {
			return err
		}
		for _, file := range files {
			if err = packFile(filepath.Join(a.opt.ItemDir, file), a.fieldsFilePath(file), a.textFilePath(file)); err != nil {
				log.Print(err)
			}
		}
	}
	return errors.New("invalid operation")
}

func (a *operation) unpack() error {
	if a.opt.ItemFile != "" {
		fieldsFile := a.opt.FieldsFile
		if fieldsFile == "" {
			fieldsFile = a.fieldsFilePath(a.opt.ItemFile)
		}
		textFile := a.opt.TextFile
		if textFile == "" {
			textFile = a.textFilePath(a.opt.ItemFile)
		}
		return unpackFile(a.opt.ItemFile, fieldsFile, textFile)
	} else if a.opt.ItemDir != "" {
		files, err := a.itemDirFiles()
		if err != nil {
			return err
		}
		for _, file := range files {
			if err = unpackFile(filepath.Join(a.opt.ItemDir, file), a.fieldsFilePath(file), a.textFilePath(file)); err != nil {
				log.Print(err)
			}
		}
		return nil
	}
	return errors.New("invalid operation")
}

const (
	RepackModeID       = "id"
	RepackModeIndex    = "index"
	RepackModeDirIndex = "dir_index"
	RepackModeTitle    = "title"
	RepackModeMeta     = "meta"
)

type repackContext struct {
	mode          string
	index         int
	titleTemplate string
}

func (a *operation) repack() error {
	ctx := &repackContext{
		mode:          a.opt.Repack,
		index:         a.opt.Skip,
		titleTemplate: a.opt.TitleTemplate,
	}

	if a.opt.ItemFile != "" {
		return repackFile(ctx, a.opt.ItemFile)
	} else if a.opt.ItemDir != "" {
		files, err := a.itemDirFiles()
		if err != nil {
			return err
		}
		for _, file := range files {
			if err = repackFile(ctx, filepath.Join(a.opt.ItemDir, file)); err != nil {
				log.Print(err)
			}
		}
		return nil
	}
	return errors.New("invalid operation")
}

func (a *operation) itemDirFiles() ([]string, error) {
	dirents, err := os.ReadDir(a.opt.ItemDir)
	if err != nil {
		return nil, err
	}
	var numeric []int
	var literal []string
	for _, dirent := range dirents {
		if dirent.IsDir() {
			continue
		}
		name := filepath.Base(dirent.Name())
		if filepath.Ext(name) != ".json" {
			continue
		}
		index, err := strconv.Atoi(strings.TrimSuffix(name, ".json"))
		if err == nil {
			numeric = append(numeric, index)
		} else if !a.opt.NumOnly {
			literal = append(literal, name)
		}
	}
	sort.Ints(numeric)
	sort.Strings(literal)
	files := make([]string, len(numeric))
	for i, index := range numeric {
		files[i] = strconv.Itoa(index) + ".json"
	}
	if a.opt.NumFirst {
		return append(files, literal...), nil
	} else {
		return append(literal, files...), nil
	}
}

func (a *operation) fieldsFilePath(name string) string {
	return filepath.Join(a.opt.FieldsDir, strings.TrimSuffix(filepath.Base(name), ".json")+".fields.json")
}

func (a *operation) textFilePath(name string) string {
	return filepath.Join(a.opt.TextDir, strings.TrimSuffix(filepath.Base(name), ".json")+".body.md")
}

func repackFile(ctx *repackContext, itemFile string) error {
	item, err := loadItem(itemFile)
	if err != nil {
		return err
	}
	return repackItem(ctx, item, itemFile)
}

func repackItem(ctx *repackContext, item *archive.Item, itemFile string) error {
	switch ctx.mode {
	case RepackModeID:
		if item.ID == "" {
			item.ID = uuid.New().String()
		}
	case RepackModeIndex:
		ctx.index++
		item.Index = ctx.index
	case RepackModeDirIndex:
		ctx.index++
		item.DirIndex = ctx.index
	case RepackModeTitle:
		ctx.index++
		item.Title = strings.ReplaceAll(ctx.titleTemplate, "{{index}}", strconv.Itoa(ctx.index))
	case RepackModeMeta:
		if sourceList, ok := item.Meta["source"].([]any); ok {
			for i, src := range sourceList {
				sourceList[i] = strings.ReplaceAll(src.(string), ",", "")
			}
			item.Meta["source"] = map[string]any{
				"ru": sourceList,
			}
		}
	}

	return saveItem(item, itemFile)
}

func packFile(itemFile, fieldsFile, textFile string) error {
	item, err := loadFields(fieldsFile)
	if err != nil {
		return err
	}
	body, err := loadText(textFile)
	if err != nil {
		return err
	}
	item.Body = body
	return saveItem(item, itemFile)
}

func unpackFile(itemFile, fieldsFile, textFile string) error {
	item, err := loadItem(itemFile)
	if err != nil {
		return err
	}
	if err := saveFields(item, fieldsFile); err != nil {
		return err
	}
	if err := saveText(item.Body, textFile); err != nil {
		return err
	}
	return nil
}

func loadItem(path string) (*archive.Item, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0775)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var item archive.Item
	if err = json.NewDecoder(f).Decode(&item); err != nil {
		return nil, err
	}
	return &item, err
}

func saveItem(item *archive.Item, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0775)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "    ")
	return e.Encode(&item)
}

func saveFields(item *archive.Item, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0775)
	if err != nil {
		return err
	}
	defer f.Close()

	body := item.Body
	defer func() {
		item.Body = body
	}()
	item.Body = ""
	e := json.NewEncoder(f)
	e.SetIndent("", "    ")
	return e.Encode(&item)
}

func loadFields(path string) (*archive.Item, error) {
	item, err := loadItem(path)
	if err != nil {
		return nil, err
	}
	item.Body = ""
	return item, err
}

func loadText(path string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("pandoc", path, "-f", "gfm", "-t", "html")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func saveText(text, path string) error {
	cmd := exec.Command("pandoc", "-f", "html", "-t", "gfm", "-o", path)
	cmd.Stdin = bytes.NewReader([]byte(text))
	return cmd.Run()
}
