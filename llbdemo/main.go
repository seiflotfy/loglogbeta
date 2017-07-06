package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/seiflotfy/loglogbeta"

	"log"
)

func main() {
	files, err := filepath.Glob(fmt.Sprintf("%s/*", os.Args[1]))
	if err != nil {
		log.Fatalln(err)
		return
	}

	totalUnique := map[string]bool{}
	tllb := loglogbeta.New()

	for _, f := range files {
		f, err := os.Open(f)
		if err != nil {
			log.Fatalln(err)
			return
		}
		reader := bufio.NewReader(f)
		unique := map[string]bool{}

		llb := loglogbeta.New()

		for {
			text, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			unique[string(text)] = true
			totalUnique[string(text)] = true
			llb.Add([]byte(text))
			tllb.Add([]byte(text))
		}

		est := llb.Cardinality()
		ratio := fmt.Sprintf("%2f%%", 100*(1-float64(len(unique))/float64(est)))
		log.Println("\n\tfile: ", f.Name(), "\n\texact:", len(unique), "\n\testimate:", est, "\n\tratio:", ratio)
	}

	est := tllb.Cardinality()
	ratio := fmt.Sprintf("%2f%%", 100*(1-float64(len(totalUnique))/float64(est)))
	log.Println("\n\ttotal\n\texact:", len(totalUnique), "\n\testimate:", est, "\n\tratio:", ratio)
}
