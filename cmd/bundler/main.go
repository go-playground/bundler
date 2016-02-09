package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-playground/bundler"
)

var (
	flagFileOrDir             = flag.String("i", "", "File or DIR to bundle files for; DIR will bundle all files within the DIR recursivly.")
	flagOuputFile             = flag.String("o", "", "Output filename, hash appended otherwise, if using a DIR in -i option then this will be the output file directory, which keeps the original file structure.")
	flagLeftDelim             = flag.String("ld", bundler.DefaultLeftDelim, "the Left Delimiter for file includes, if not specified default to "+bundler.DefaultLeftDelim+".")
	flagRightDelim            = flag.String("rd", bundler.DefaultRightDelim, "the Right Delimiter for file includes, if not specified default to "+bundler.DefaultRightDelim+".")
	flagIncludesRelativeToDir = flag.Bool("rtd", true, "Specifies if the files included should be treated as relative to the directory, or relative to the files from which they are included.")
	flagIgnore                = flag.String("ignore", "", "Regexp for files/dirs we should ignore i.e. \\.gitignore; only used when -i option is a DIR")

	input      string
	output     string
	leftDelim  string
	rightDelim string
	ignore     string

	relativeToDir bool

	ignoreRegexp *regexp.Regexp

	processed []*bundler.ProcessedFile
)

func main() {
	parseFlags()

	var err error

	fi, err := os.Stat(input)
	if err != nil {
		panic(err)
	}

	//process multiple files
	if fi.IsDir() {
		processed, err = bundler.BundleDir(input, output, relativeToDir, input, leftDelim, rightDelim, ignoreRegexp)
		if err != nil {
			panic(err)
		}
		printResults(processed...)
		return
	}

	// process file
	file, err := bundler.BundleFile(input, output, relativeToDir, input, leftDelim, rightDelim)
	if err != nil {
		panic(err)
	}

	processed = append(processed, file)
	printResults(processed...)
}

func printResults(processed ...*bundler.ProcessedFile) {

	fmt.Printf("The following files were processed:\n\n")

	for _, file := range processed {
		fmt.Println("  " + file.NewFilename)
	}

	fmt.Printf("\n\n")
}

func parseFlags() {

	flag.Parse()

	input = strings.TrimSpace(*flagFileOrDir)
	output = *flagOuputFile
	leftDelim = *flagLeftDelim
	rightDelim = *flagRightDelim
	relativeToDir = *flagIncludesRelativeToDir
	ignore = *flagIgnore

	wasBlank := len(input) == 0

	input = filepath.Clean(input)

	if wasBlank && input == "." {
		panic("** No File Or Directory specified with -i option")
	}

	if len(leftDelim) == 0 {
		leftDelim = bundler.DefaultLeftDelim
	}

	if len(rightDelim) == 0 {
		rightDelim = bundler.DefaultRightDelim
	}

	if len(ignore) > 0 {
		var err error

		ignoreRegexp, err = regexp.Compile(ignore)
		if err != nil {
			panic("**Error Compiling Regex:" + err.Error())
		}
	}
}
