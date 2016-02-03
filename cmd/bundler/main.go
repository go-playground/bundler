package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

	useHash   bool
	isDirMode bool

	ignoreRegexp *regexp.Regexp
)

func main() {
	parseFlags()

	fi, err := os.Stat(input)
	if err != nil {
		panic(err)
	}

	//process multiple files
	if fi.IsDir() {
		isDirMode = true
		processDir(input, "", false, "")
		return
	}

	// process file
	processFile(input)
}

func processFile(path string) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	fmt.Println("Writing Temp File:", f.Name()+"-")
	newFile, err := ioutil.TempFile("", f.Name()+"-")
	if err != nil {
		return err
	}

	bundler.Bundle(f, newFile, filepath.Dir(path), leftDelim, rightDelim)

	var newName string

	dirname, filename := filepath.Split(path)
	ext := filepath.Ext(filename)
	filename = filepath.Base(filename)
	// filename = filename[0:strings.LastIndex(filename, ext)]

	newName = dirname + filename[0:strings.LastIndex(filename, ext)]

	if useHash {
		b, err := ioutil.ReadFile(newFile.Name())
		if err != nil {
			return err
		}

		h := md5.New()
		h.Write(b)
		// io.WriteString(h, string(content))
		// hash
		hash := string(h.Sum(nil))

		newName += "-" + hash + ext

	} else {
		if isDirMode {
			newName += "-" + output + ext
		}

		newName = dirname + output
	}

	fmt.Println("Renaming from", newFile.Name(), "to", newName)

	os.Rename(newFile.Name(), newName)

	return nil
}

func processDir(path string, dir string, isSymlinkDir bool, symlinkDir string) error {

	var p string

	f, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range files {

		info := file
		p = path + string(os.PathSeparator) + file.Name()
		fPath := p

		if isSymlinkDir {
			fPath = strings.Replace(p, dir, symlinkDir, 1)
		}

		if ignoreRegexp != nil && ignoreRegexp.MatchString(fPath) {
			continue
		}

		if file.IsDir() {
			processDir(p, p, isSymlinkDir, symlinkDir+string(os.PathSeparator)+info.Name())
			continue
		}

		if file.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err := filepath.EvalSymlinks(p)
			if err != nil {
				log.Panic("Error Resolving Symlink", err)
			}

			fi, err := os.Stat(link)
			if err != nil {
				log.Panic(err)
			}

			info = fi

			if fi.IsDir() {
				processDir(link, link, true, fPath)
				continue
			}
		}
	}

	// call processFile()
	processFile(p)

	return nil
}

func parseFlags() {

	flag.Parse()

	input = *flagFileOrDir
	output = *flagOuputFile
	leftDelim = *flagLeftDelim
	rightDelim = *flagRightDelim
	ignore = *flagIgnore

	input = filepath.Clean(input)

	if len(strings.TrimSpace(input)) == 0 {
		panic("** No File Or Directory specified with -i option")
	}

	if len(output) == 0 {
		useHash = true
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
