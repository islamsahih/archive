package main

import (
	"archive/util"
	"flag"
	"log"
)

var (
	srcFlag    = flag.String("src", "", "")
	dstFlag    = flag.String("dst", "", "")
	appendFlag = flag.Bool("append", false, "")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if srcFlag == nil || dstFlag == nil {
		log.Fatal("invalid params")
	}

	html, err := util.Fetch(*srcFlag)
	if err != nil {
		log.Fatal(err)
	}

	if err = util.SaveHTML(html, *dstFlag, *appendFlag); err != nil {
		log.Fatal(err)
	}
}
