package bundler

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// BundleDir bundles an entire directory recursively and returns an array of filenames and if an error occured processing
// suffix will be appended to filenames, if blank a hash of file contents will be added
func BundleDir(dirname string, suffix string, leftDelim string, rightDelim string, ignoreRegexp *regexp.Regexp) ([]string, error) {
	return bundleDir(dirname, "", false, "", ignoreRegexp, suffix, leftDelim, rightDelim)
}

func bundleDir(path string, dir string, isSymlinkDir bool, symlinkDir string, ignoreRegexp *regexp.Regexp, output string, leftDelim string, rightDelim string) ([]string, error) {

	var p string
	var processed []string

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		return nil, err
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

			processedFiles, err := bundleDir(p, p, isSymlinkDir, symlinkDir+string(os.PathSeparator)+info.Name(), ignoreRegexp, output, leftDelim, rightDelim)
			if err != nil {
				return nil, err
			}

			processed = append(processed, processedFiles...)

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

				processedFiles, err := bundleDir(link, link, true, fPath, ignoreRegexp, output, leftDelim, rightDelim)
				if err != nil {
					return nil, err
				}

				processed = append(processed, processedFiles...)

				continue
			}
		}

		// process file
		file, err := bundleFile(p, output, leftDelim, rightDelim, true)
		if err != nil {
			return nil, err
		}

		processed = append(processed, file)
	}

	return processed, nil
}

// BundleFile bundles a single file on disk and returns the filename and if an error occured processing
func BundleFile(path string, output string, leftDelim string, rightDelim string) (string, error) {
	return bundleFile(path, output, leftDelim, rightDelim, false)
}

func bundleFile(path string, output string, leftDelim string, rightDelim string, isDirMode bool) (string, error) {

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	// fmt.Println("Writing Temp File for:", f.Name())
	newFile, err := ioutil.TempFile("", filepath.Base(f.Name()))
	if err != nil {
		return "", err
	}

	if err = Bundle(f, newFile, filepath.Dir(path), leftDelim, rightDelim); err != nil {
		return "", err
	}

	var newName string

	dirname, filename := filepath.Split(path)
	ext := filepath.Ext(filename)
	filename = filepath.Base(filename)

	newName = dirname + filename[0:strings.LastIndex(filename, ext)]

	if output == "" {
		b, err := ioutil.ReadFile(newFile.Name())
		if err != nil {
			return "", err
		}

		h := md5.New()
		h.Write(b)
		hash := string(h.Sum(nil))

		newName += "-" + fmt.Sprintf("%x", hash) + ext

	} else {
		if isDirMode {
			newName += "-" + output + ext
		} else {
			newName = dirname + output
		}
	}

	// fmt.Println("Renaming from", newFile.Name(), "to", newName)

	if err = os.Rename(newFile.Name(), newName); err != nil {
		return "", err
	}

	return newName, nil
}

// Bundle combines the given input and writes it out to the provided writer
// removing delims from the combined files
func Bundle(r io.Reader, w io.Writer, dir string, leftDelim string, rightDelim string) error {
	return bundle(r, w, dir, leftDelim, rightDelim, false)
}

// BundleKeepDelims combines the given input and writes it out to the provided writer
// but unlike Bundle() keeps the delims in the combined data
func BundleKeepDelims(r io.Reader, w io.Writer, dir string, leftDelim string, rightDelim string) error {
	return bundle(r, w, dir, leftDelim, rightDelim, true)
}

func bundle(r io.Reader, w io.Writer, dir string, leftDelim string, rightDelim string, keepDelims bool) error {

	var err error

	if !filepath.IsAbs(dir) {
		if dir, err = filepath.Abs(dir); err != nil {
			return err
		}
	}

	fi, err := os.Lstat(dir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {

		// check if symlink

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {

			link, err := filepath.EvalSymlinks(dir)
			if err != nil {
				return errors.New("Error Resolving Symlink:" + err.Error())
			}

			fi, err = os.Stat(link)
			if err != nil {
				return err
			}

			if !fi.IsDir() {
				return errors.New("dir passed is not a directory")
			}

			dir = link

		} else {
			return errors.New("dir passed is not a directory")
		}
	}

	l, err := NewLexer("bundle", r, leftDelim, rightDelim)
	if err != nil {
		return err
	}

LOOP:
	for {
		itm := l.NextItem()

		switch itm.typ {
		case itemLeftDelim, itemRightDelim:
			if keepDelims {
				w.Write([]byte(itm.val))
			}
		case itemText:
			w.Write([]byte(itm.val))
		case itemFile:
			path := dir + "/" + itm.val

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if err = bundle(file, w, filepath.Dir(path), leftDelim, rightDelim, keepDelims); err != nil {
				return err
			}
		case itemEOF:
			break LOOP
		}
	}

	return nil
}
