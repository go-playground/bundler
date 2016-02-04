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
	flagFileOrDir  = flag.String("i", "", "File or DIR to bundle files for; DIR will bundle all files within the DIR recursivly.")
	flagOuputFile  = flag.String("o", "", "Output filename, or if using a DIR in -i option the suffix, otherwise will be be the filename with appended hash of file contents.")
	flagLeftDelim  = flag.String("ld", bundler.DefaultLeftDelim, "the Left Delimiter for file includes, if not specified default to "+bundler.DefaultLeftDelim+".")
	flagRightDelim = flag.String("rd", bundler.DefaultRightDelim, "the Right Delimiter for file includes, if not specified default to "+bundler.DefaultRightDelim+".")
	flagIgnore     = flag.String("ignore", "", "Regexp for files/dirs we should ignore i.e. \\.gitignore; only used when -i option is a DIR")

	input      string
	output     string
	leftDelim  string
	rightDelim string
	ignore     string

	ignoreRegexp *regexp.Regexp

	processed []string
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
		processed, err = bundler.BundleDir(input, output, leftDelim, rightDelim, ignoreRegexp)
		if err != nil {
			panic(err)
		}
		PrintResults(processed)
		return
	}

	// process file
	file, err := bundler.BundleFile(input, output, leftDelim, rightDelim)
	if err != nil {
		panic(err)
	}

	processed = append(processed, file)
	PrintResults(processed)
}

func PrintResults(processed []string) {
	fmt.Println("The following files were processed:\n")

	for _, file := range processed {
		fmt.Println("  " + file)
	}

	fmt.Println("\n")
}

func parseFlags() {

	flag.Parse()

	input = strings.TrimSpace(*flagFileOrDir)
	output = *flagOuputFile
	leftDelim = *flagLeftDelim
	rightDelim = *flagRightDelim
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
