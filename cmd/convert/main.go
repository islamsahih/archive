package main

import (
	"archive/converter"
	"flag"
	"log"
	"os"
)

var (
	inputDirFlag  = flag.String("input-dir", "", "")
	outputDirFlag = flag.String("output-dir", "", "")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *outputDirFlag != "" {
		if err := os.MkdirAll(*outputDirFlag, 0775); err != nil {
			log.Fatal(err)
		}
	}

	c, err := converter.New(
		converter.Options{
			InputDir:  *inputDirFlag,
			OutputDir: *outputDirFlag,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = c.Convert(); err != nil {
		log.Fatal(err)
	}
}
