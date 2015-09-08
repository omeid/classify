package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/omeid/classify"
)

var (
	wl   = flag.String("terms", "", "The file that holds the term list.")
	text = flag.String("text", "", "The file that holds the text. Defaults to os.Stdin")
)

func main() {
	flag.Parse()

	if *wl == "" {
		log.Fatal("Word list is required.")
	}

	wordlist, err := os.Open(*wl)
	if err != nil {
		log.Fatal(err)
	}

	terms, err := classify.FromCSV(wordlist)
	if err != nil {
		log.Fatal(err)
	}

	var content io.Reader
	if *text == "" {
		content = os.Stdin
	} else {
		content, err = os.Open(*text)
		if err != nil {
			log.Fatal(err)
		}
	}

	results, err := terms.Analyse(content)
	if err != nil {
		log.Fatal(err)
	}

	for cat, r := range results {
		fmt.Printf("Category %s: %v\n", cat, r)
	}
}
