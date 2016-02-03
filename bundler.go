package bundler

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

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
		return errors.New("dir passed is not a directory")
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
