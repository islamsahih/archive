package main

import (
	"archive/fixer"
	"flag"
	"log"
	"strings"
)

var (
	inputDirFlag = flag.String("input-dir", "", "")
	indexesFlag  = flag.String("indexes", "", "")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	fix, err := fixer.New(
		fixer.Options{
			InputDir: *inputDirFlag,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if *indexesFlag != "" {
		for _, category := range strings.Split(*indexesFlag, ",") {
			if err = fix.Indexes(category); err != nil {
				log.Fatal(err)
			}
		}
	}

	//fix.Tags()

	fix.Save()
}
