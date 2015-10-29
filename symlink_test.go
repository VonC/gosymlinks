package symlink

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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

	osStat = testOsStat
	execRun = testExecRun
	tests := []*test{
		&test{dst: "unknown/dst", err: "The system cannot find the path specified"},
		&test{dst: string([]byte{0}), err: "invalid argument"},
		&test{dst: "err", err: "Test error on os.Stat with non-nil fi"},
	}
	var sl *SL
	var err error
	for _, test := range tests {
		sl, err = New(".", test.dst)
		if err == nil || strings.Contains(err.Error(), test.err) == false {
			t.Errorf("Err '%v', expected '%s'", err, test.err)
		}
		if sl != nil {
			t.Errorf("SL '%v', expected <nil>", sl)
		}
	}
	_, err = New(`.`, `symlink`)
	fmt.Printf("%+v\n", err)
}

func testOsStat(name string) (os.FileInfo, error) {
	if strings.HasSuffix(name, `prj\symlink\err\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("Test error on os.Stat with non-nil fi")
	}
	if strings.HasSuffix(name, `prj\symlink\symlink\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("readlink for symlink")
	}
	return os.Stat(name)
}

func testExecRun(cmd *exec.Cmd) error {
	fmt.Printf("cmd='%+v'\n", cmd.Dir)
	if strings.HasSuffix(cmd.Dir, `\symlink`) {
		out := `
 Répertoire de C:\Users\VonC\prog\git\ggb\deps\src\github.com\VonC

22/06/2015  11:03    <REP>          .
22/06/2015  11:03    <REP>          ..
22/06/2015  11:03    <JONCTION>     symlink [C:\Users\VonC\prog\git\ggb\]`
		io.WriteString(cmd.Stdout, out)
		// fmt.Printf("i='%d', e='%+v'\n", i, e)
		return nil
	} else {
		return cmd.Run()
	}
}
