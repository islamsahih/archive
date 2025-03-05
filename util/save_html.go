package util

import (
	"os"
	"path/filepath"
	"strings"
)

var htmlReplacer = strings.NewReplacer(
	"><", ">\n<",
	"📚", "",
	"🔹", "",
	"➡️", "",
	"🖊", "",
	"📜", "",
	"🖇", "",
	"▶️", "",
	"🔸", "",
)

func SaveHTML(html string, path string, append bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		return err
	}
	flag := os.O_CREATE | os.O_WRONLY
	if append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	f, err := os.OpenFile(path, flag, 0660)
	if err != nil {
		return err
	}
	defer f.Close()

	html = htmlReplacer.Replace(html)
	html = strings.TrimPrefix(html, ":")
	html = strings.TrimPrefix(html, "<br/>\n<br/>\n<br/>\n<br/>")

	_, err = f.WriteString(html)
	return err
}
