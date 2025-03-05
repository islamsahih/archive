package main

import (
	"archive/builder"
	"flag"
	"log"
)

var (
	inputDirFlag  = flag.String("input-dir", "", "")
	outputDirFlag = flag.String("output-dir", "", "")
	formatFlag    = flag.String("format", "html", "")
	categoryFlag  = flag.String("category", "", "")
	indexFlag     = flag.Int("index", 0, "")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	builder, err := builder.New(
		builder.Options{
			InputDir:  *inputDirFlag,
			OutputDir: *outputDirFlag,
			Format:    *formatFlag,
			Category:  *categoryFlag,
			Index:     *indexFlag,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = builder.Build(); err != nil {
		log.Fatal(err)
	}
}
