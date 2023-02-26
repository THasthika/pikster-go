package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/THasthika/pikster-go"
)

func main() {

	clusterCount := flag.Uint("c", 2, "Number of Clusters")
	maxIters := flag.Uint("iter", 100, "Maximum Iterations")
	outFileName := flag.String("o", "", "Output filename")

	flag.Parse()

	var inFile string

	if flag.NArg() == 0 {
		reader := bufio.NewReader(os.Stdin)
		var err error
		inFile, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln("failed to read input")
		}
		inFile = strings.TrimSpace(inFile)
	} else {
		inFile = flag.Arg(0)
	}

	outFile := *outFileName
	if *outFileName == "" {
		ext := filepath.Ext(inFile)
		beforeExt, _ := strings.CutSuffix(inFile, ext)
		outFile = beforeExt + ".out"
	}

	pImg := pikster.NewPImageFromFile(inFile, *clusterCount)

	for i := 0; i < int(*maxIters); i++ {
		if pImg.RunClusteringStep() == 0 {
			break
		}
	}

	pImg.SaveFile(outFile)

}
