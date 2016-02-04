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

func TestDirHash(t *testing.T) {

	pwd, err := os.Getwd()
	Equal(t, err, nil)

	// ignoreRegexp, err = regexp.Compile(ignore)

	filenames, err := BundleDir(pwd+"/testfiles/test1", "", "include(", ")", nil)
	Equal(t, err, nil)
	Equal(t, len(filenames), 3)

	b, err := ioutil.ReadFile(filenames[0])
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filenames[0])
	Equal(t, err, nil)

	b, err = ioutil.ReadFile(filenames[1])
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), "- File 2")

	err = os.Remove(filenames[1])
	Equal(t, err, nil)

	b, err = ioutil.ReadFile(filenames[2])
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), "- File 3")

	err = os.Remove(filenames[2])
	Equal(t, err, nil)
}

func TestDirSuffix(t *testing.T) {

	pwd, err := os.Getwd()
	Equal(t, err, nil)

	// ignoreRegexp, err = regexp.Compile(ignore)

	filenames, err := BundleDir(pwd+"/testfiles/test1", "testsuffix", "include(", ")", nil)
	Equal(t, err, nil)
	Equal(t, len(filenames), 3)

	b, err := ioutil.ReadFile(filenames[0])
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)
	Equal(t, filepath.Base(filenames[0]), "file1-testsuffix.txt")

	Equal(t, string(b), final1)

	err = os.Remove(filenames[0])
	Equal(t, err, nil)

	Equal(t, filepath.Base(filenames[1]), "file2-testsuffix.txt")

	b, err = ioutil.ReadFile(filenames[1])
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), "- File 2")

	err = os.Remove(filenames[1])
	Equal(t, err, nil)

	Equal(t, filepath.Base(filenames[2]), "file3-testsuffix.txt")

	b, err = ioutil.ReadFile(filenames[2])
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), "- File 3")

	err = os.Remove(filenames[2])
	Equal(t, err, nil)
}

func TestSingleFileHash(t *testing.T) {

	pwd, err := os.Getwd()
	Equal(t, err, nil)

	filename, err := BundleFile(pwd+"/testfiles/test1/file1.txt", "", "include(", ")")
	Equal(t, err, nil)
	NotEqual(t, filepath.Base(filename), "file1.txt")

	b, err := ioutil.ReadFile(filename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filename)
	Equal(t, err, nil)
}

func TestSingleFileWithOutput(t *testing.T) {

	pwd, err := os.Getwd()
	Equal(t, err, nil)

	filename, err := BundleFile(pwd+"/testfiles/test1/file1.txt", "test.txt", "include(", ")")
	Equal(t, err, nil)
	Equal(t, filepath.Base(filename), "test.txt")

	b, err := ioutil.ReadFile(filename)
	Equal(t, err, nil)
	NotEqual(t, len(b), 0)

	Equal(t, string(b), final1)

	err = os.Remove(filename)
	Equal(t, err, nil)
}
