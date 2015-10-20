package symlink

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SL struct {
	dst  string
	path string
}

func New(dst string) (*SL, error) {
	var err error
	if dst, err = dirAbsPath(dst); err != nil {
		return nil, err
	}
	fi, err := os.Stat(dst)
	fmt.Printf("%+v %s\n", fi, dst)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// a/b => c:\path\to\a\b\
func dirAbsPath(path string) (string, error) {
	path = filepath.FromSlash(path)
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	sep := string(filepath.Separator)
	if strings.HasSuffix(path, sep) == false {
		path = path + sep
	}
	return path, nil
}

func dirExists(path string) (bool, string, error) {
	fi, err := os.Stat(path)
	if fi == nil {
		return false, "", err
	}
	return false, "", nil
}
