package main

import (
	"archive/content"
	"flag"
	"log"
)

var (
	packFlag          = flag.Bool("pack", false, "")
	unpackFlag        = flag.Bool("unpack", false, "")
	repackFlag        = flag.String("repack", "", "")
	itemDirFlag       = flag.String("item-dir", "", "")
	itemFileFlag      = flag.String("item-file", "", "")
	textDirFlag       = flag.String("text-dir", "", "")
	textFileFlag      = flag.String("text-file", "", "")
	fieldsDirFlag     = flag.String("fields-dir", "", "")
	fieldsFileFlag    = flag.String("fields-file", "", "")
	numOnlyFlag       = flag.Bool("numeric-only", true, "")
	numFirstFlag      = flag.Bool("numeric-first", true, "")
	skipFlag          = flag.Int("skip", 0, "")
	titleTemplateFlag = flag.String("title-template", "", "")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	operation, err := content.NewOperation(content.Options{
		Pack:          *packFlag,
		Unpack:        *unpackFlag,
		Repack:        *repackFlag,
		ItemDir:       *itemDirFlag,
		ItemFile:      *itemFileFlag,
		TextDir:       *textDirFlag,
		TextFile:      *textFileFlag,
		FieldsDir:     *fieldsDirFlag,
		FieldsFile:    *fieldsFileFlag,
		NumFirst:      *numFirstFlag,
		NumOnly:       *numOnlyFlag,
		Skip:          *skipFlag,
		TitleTemplate: *titleTemplateFlag,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err = operation.Exec(); err != nil {
		log.Fatal(err)
	}
}
