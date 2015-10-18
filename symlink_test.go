package symlink

import (
	"strings"
	"testing"
)

// Making a symlink means:
// - making sure the destination exists
// - making sure the source don't exist and the parent folder is created
// - if the source does exist and
//   - is a symlink pointing to a different destination: rmdir
//   - is a symlink pointing to the same
//   - is a folder x, rename it to x.1 (if x.1 exists, x.2, ...)

type test struct {
	dst string
	err string
	sl  *SL
}

func TestDestination(t *testing.T) {

	// only a nil bit will make filepath.Abs() fail:
	// https://github.com/golang/go/blob/d16c7f8004bd1c9f896367af7ea86f5530596b39/src/syscall/syscall_windows.go#L41
	// from UTF16FromString (https://github.com/golang/go/blob/d16c7f8004bd1c9f896367af7ea86f5530596b39/src/syscall/syscall_windows.go#L71)
	// from FullPath (https://github.com/golang/go/blob/d16c7f8004bd1c9f896367af7ea86f5530596b39/src/syscall/exec_windows.go#L134)
	// from abs (https://github.com/golang/go/blob/d16c7f8004bd1c9f896367af7ea86f5530596b39/src/path/filepath/path_windows.go#L109)
	// from Abs (https://github.com/golang/go/blob/d16c7f8004bd1c9f896367af7ea86f5530596b39/src/path/filepath/path.go#L235)

	tests := []*test{
		&test{dst: "unknown/dst", err: "The system cannot find the path specified"},
		&test{dst: string([]byte{0}), err: "invalid argument"},
	}
	var sl *SL
	var err error
	for _, test := range tests {
		sl, err = New(test.dst)
		if err == nil || strings.Contains(err.Error(), test.err) == false {
			t.Errorf("Err '%v', expected '%s'", err, test.err)
		}
		if sl != nil {
			t.Errorf("SL '%v', expected <nil>", sl)
		}
	}
	New(".")
}
