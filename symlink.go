package symlink

import (
	"os"
	"path/filepath"
)

type SL struct {
	dst  string
	path string
}

func New(dst string) (*SL, error) {
	dst = filepath.FromSlash(dst)
	var err error
	dst, err = filepath.Abs(dst)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(dst)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
