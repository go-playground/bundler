package bundler

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

const (
	final1 = `- File 2
- File 1
- File 3`
)

func TestSingleFileHashWithSymlinks(t *testing.T) {

	// pwd, err := os.Getwd()
	// Equal(t, err, nil)

	filename, err := BundleFile("testfiles/test2/file1.txt", "", false, "", "include(", ")")
	Equal(t, err, nil)
	NotEqual(t, filepath.Base(filename.NewFilename), filepath.Base(filename.OriginalFilename))

	b, err := ioutil.ReadFile(filename.NewFilename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filename.NewFilename)
	Equal(t, err, nil)
}

func TestSingleFileWithOutputAndSymlinks(t *testing.T) {

	// pwd, err := os.Getwd()
	// Equal(t, err, nil)

	filename, err := BundleFile("testfiles/test2/file1.txt", "test.txt", false, "", "include(", ")")
	Equal(t, err, nil)
	Equal(t, filepath.Base(filename.NewFilename), "test.txt")

	b, err := ioutil.ReadFile(filename.NewFilename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filename.NewFilename)
	Equal(t, err, nil)
}

func TestDirHash(t *testing.T) {

	// pwd, err := os.Getwd()
	// Equal(t, err, nil)

	filenames, err := BundleDir("testfiles/test1", "", false, "", "include(", ")", nil)
	Equal(t, err, nil)
	Equal(t, len(filenames), 3)

	b, err := ioutil.ReadFile(filenames[0].NewFilename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filenames[0].NewFilename)
	Equal(t, err, nil)

	b, err = ioutil.ReadFile(filenames[1].NewFilename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), "- File 2")

	err = os.Remove(filenames[1].NewFilename)
	Equal(t, err, nil)

	b, err = ioutil.ReadFile(filenames[2].NewFilename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), "- File 3")

	err = os.Remove(filenames[2].NewFilename)
	Equal(t, err, nil)
}

// func TestDirSuffix(t *testing.T) {

// 	pwd, err := os.Getwd()
// 	Equal(t, err, nil)

// 	// ignoreRegexp, err = regexp.Compile(ignore)

// 	filenames, err := BundleDir(pwd+"/testfiles/test1", "testsuffix", false, "", "include(", ")", nil)
// 	Equal(t, err, nil)
// 	Equal(t, len(filenames), 3)

// 	b, err := ioutil.ReadFile(filenames[0].NewFilename)
// 	Equal(t, err, nil)
// 	NotEqual(t, len(b), 0)
// 	Equal(t, filepath.Base(filenames[0].NewFilename), "file1-testsuffix.txt")

// 	Equal(t, string(b), final1)

// 	err = os.Remove(filenames[0].NewFilename)
// 	Equal(t, err, nil)

// 	Equal(t, filepath.Base(filenames[1].NewFilename), "file2-testsuffix.txt")

// 	b, err = ioutil.ReadFile(filenames[1].NewFilename)
// 	Equal(t, err, nil)
// 	NotEqual(t, len(b), 0)

// 	Equal(t, string(b), "- File 2")

// 	err = os.Remove(filenames[1].NewFilename)
// 	Equal(t, err, nil)

// 	Equal(t, filepath.Base(filenames[2].NewFilename), "file3-testsuffix.txt")

// 	b, err = ioutil.ReadFile(filenames[2].NewFilename)
// 	Equal(t, err, nil)
// 	NotEqual(t, len(b), 0)

// 	Equal(t, string(b), "- File 3")

// 	err = os.Remove(filenames[2].NewFilename)
// 	Equal(t, err, nil)
// }

func TestSingleFileHash(t *testing.T) {

	// pwd, err := os.Getwd()
	// Equal(t, err, nil)

	// filename, err := BundleFile(pwd+"/testfiles/test1/file1.txt", "", false, "", "include(", ")")
	filename, err := BundleFile("testfiles/test1/file1.txt", "", false, "", "include(", ")")
	Equal(t, err, nil)
	NotEqual(t, filepath.Base(filename.NewFilename), filepath.Base(filename.OriginalFilename))

	b, err := ioutil.ReadFile(filename.NewFilename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filename.NewFilename)
	Equal(t, err, nil)
}

func TestSingleFileWithOutput(t *testing.T) {

	// pwd, err := os.Getwd()
	// Equal(t, err, nil)

	filename, err := BundleFile("testfiles/test1/file1.txt", "test.txt", false, "", "include(", ")")
	Equal(t, err, nil)
	Equal(t, filepath.Base(filename.NewFilename), "test.txt")

	b, err := ioutil.ReadFile(filename.NewFilename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filename.NewFilename)
	Equal(t, err, nil)
}
