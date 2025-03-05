package main

import (
	"archive/parser"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	inputFlag          = flag.String("input", "", "")
	inputDirFlag       = flag.String("input-dir", "", "")
	outputFlag         = flag.String("output", "", "")
	outputDirFlag      = flag.String("output-dir", "", "")
	skipIndexFlag      = flag.Int("skip-index", 0, "")
	skipItemsFlag      = flag.String("skip-items", "", "")
	audioOutputDirFlag = flag.String("audio-output-dir", "", "")
	textOutputDirFlag  = flag.String("text-output-dir", "", "")
	categoryFlag       = flag.String("category", "", "")
	selfLinkMatchFlag  = flag.String("self-link-match", "", "")
	videoMatchFlag     = flag.String("video-match", "", "")
	textMatchFlag      = flag.String("text-match", "", "")
	tagMatchFlag       = flag.String("tag-match", "", "")
	linkTypeMatchFlag  = flag.String("link-type-match", "", "")
	skipLinkMatchFlag  = flag.String("skip-link-match", "", "")
	saveTextFlag       = flag.Bool("save-text", false, "")
	audioRequiredFlag  = flag.Bool("audio-required", false, "")
	titleRequiredFlag  = flag.Bool("title-required", true, "")
	processAudioFlag   = flag.Bool("process-audio", false, "")
	processVideoFlag   = flag.Bool("process-video", false, "")
	processTextFlag    = flag.Bool("process-text", false, "")
	noAlertsFlag       = flag.Bool("no-alerts", false, "")
	noCheckFlag        = flag.Bool("no-check", false, "")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *categoryFlag == "" {
		log.Fatal("--category is required")
	}
	if *processAudioFlag {
		if outputFlag == nil || *outputDirFlag == "" {
			log.Fatal("--output or --output-dir is required")
		}
	}

	var err error
	var r io.ReadCloser
	var w io.WriteCloser

	if *inputDirFlag != "" {
		*inputFlag = filepath.Join(*inputDirFlag, "messages.html")
	}
	if *inputFlag != "" {
		if *inputDirFlag == "" {
			*inputDirFlag = filepath.Dir(*inputFlag)
		}
		if r, err = os.OpenFile(*inputFlag, os.O_RDONLY, 0); err != nil {
			log.Fatal(err)
		}
		defer r.Close()
	} else {
		r = os.Stdin
	}

	if *outputDirFlag != "" {
		//*outputDirFlag = filepath.Join(*outputDirFlag, *categoryFlag)
		//if err = os.RemoveAll(*outputDirFlag); err != nil {
		//	log.Fatal(err)
		//}
		*outputFlag = filepath.Join(*outputDirFlag, strings.ReplaceAll(*categoryFlag, "/", "_")+".json")
	}
	if *outputFlag != "" {
		if *outputDirFlag == "" {
			*outputDirFlag = filepath.Dir(*outputFlag)
		}
		if err = os.MkdirAll(*outputDirFlag, 0775); err != nil {
			log.Fatal(err)
		}
		if w, err = os.OpenFile(*outputFlag, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660); err != nil {
			log.Fatal(err)
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}
	if *audioOutputDirFlag == "" && *outputDirFlag != "" {
		*audioOutputDirFlag = filepath.Join(*outputDirFlag, "audio")
	}
	if *textOutputDirFlag == "" && *outputDirFlag != "" {
		*textOutputDirFlag = filepath.Join(*outputDirFlag, "pages")
	}

	linkTypeMatch := make(map[string][]string)
	if *linkTypeMatchFlag != "" {
		for _, ltm := range strings.Split(*linkTypeMatchFlag, ";") {
			ltm := strings.Split(ltm, ":")
			if len(ltm) != 2 {
				log.Fatalln("invalid link type match parameter")
			}
			linkTypeMatch[ltm[0]] = strings.Split(ltm[1], ",")
		}
	}

	var selfLinkMatch, videoMatch, textMatch, tagMatch, skipLinkMatch []string
	if *videoMatchFlag != "" {
		videoMatch = strings.Split(*videoMatchFlag, ",")
	}
	if *textMatchFlag != "" {
		textMatch = strings.Split(*textMatchFlag, ",")
	}
	if *tagMatchFlag != "" {
		tagMatch = strings.Split(*tagMatchFlag, ",")
	}
	if *selfLinkMatchFlag != "" {
		selfLinkMatch = strings.Split(*selfLinkMatchFlag, ",")
	}
	if *skipLinkMatchFlag != "" {
		skipLinkMatch = strings.Split(*skipLinkMatchFlag, ",")
	}

	p, err := parser.New(r, w, parser.Options{
		Category:       *categoryFlag,
		InputDir:       *inputDirFlag,
		OutputDir:      *outputDirFlag,
		AudioOutputDir: *audioOutputDirFlag,
		TextOutputDir:  *textOutputDirFlag,
		SkipIndex:      *skipIndexFlag,
		SkipItems:      *skipItemsFlag,
		VideoMatch:     videoMatch,
		TextMatch:      textMatch,
		TagMatch:       tagMatch,
		SelfLinkMatch:  selfLinkMatch,
		LinkTypeMatch:  linkTypeMatch,
		SkipLinkMatch:  skipLinkMatch,
		SaveText:       *saveTextFlag,
		AudioRequired:  *audioRequiredFlag,
		TitleRequired:  *titleRequiredFlag,
		ProcessAudio:   *processAudioFlag,
		ProcessVideo:   *processVideoFlag,
		ProcessText:    *processTextFlag,
		NoAlerts:       *noAlertsFlag,
		NoCheck:        *noCheckFlag,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = p.Parse(); err != nil {
		log.Fatal(err)
	}
}
