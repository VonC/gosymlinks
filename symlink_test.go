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
	src string
	dst string
	err string
	sl  *SL
}

func TestDestination(t *testing.T) {
	t.Skip("Skip TestDestination")
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
		&test{dst: "badsymlink/dir", err: "unreadable dir on symlink"},
		&test{dst: "nojunction/dir", err: "Unable to find junction symlink in parent dir"},
		&test{dst: "cmdRun/dir", err: "The system cannot find the file specified"},
		&test{dst: "WarningOnDir/dir", err: "Warning on run"},
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
	// destination is a symlink
	_, err = New(`.`, `symlink`)
	// destination exists
	_, err = New(`x`, `.`)
	// fmt.Printf("%+v\n", err)
}

func TestSource(t *testing.T) {
	osStat = testOsStat
	execRun = testExecRun
	osMkdirAll = testOsMkdirAll
	osRename = testOsRename

	tests := []*test{
		// &test{src: "parentNotYetCreated/newlink"},
		// &test{src: "badSrcParent/newlink", err: "Test error badSrcParent on os.Stat with non-nil fi"},
		// &test{src: "badSrcParentMdirAll/newlink", err: "Error on mkDirAll for"},
		&test{src: "symlinkdir/newlink", err: "Error on ss"},
	}
	var sl *SL
	var err error
	for _, test := range tests {
		sl, err = New(test.src, ".")
		if err != nil && strings.Contains(err.Error(), test.err) == false {
			t.Errorf("Err '%v', expected '%s'", err, test.err)
		}
		if err == nil && test.err != "" {
			t.Errorf("Err nil, expected '%s'", test.err)
		}
		if sl == nil && err == nil {
			t.Errorf("SL '%v', expected NOT <nil>", sl)
		}
	}
}

func testOsStat(name string) (os.FileInfo, error) {
	fmt.Printf("testOsStat name='%+v'\n", name)
	if strings.HasSuffix(name, `prj\symlink\err\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("Test error on os.Stat with non-nil fi")
	}
	if strings.HasSuffix(name, `prj\symlink\symlink\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("readlink for symlink")
	}
	if strings.HasSuffix(name, `prj\symlink\badsymlink\dir\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("readlink for bad symlink")
	}
	if strings.HasSuffix(name, `prj\symlink\nojunction\dir\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("readlink for no junction")
	}
	if strings.HasSuffix(name, `prj\symlink\cmdRun\dir\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("readlink for no junction")
	}
	if strings.HasSuffix(name, `prj\symlink\WarningOnDir\dir\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("readlink for warning on dir")
	}
	if strings.HasSuffix(name, `prj\symlink\badSrcParent\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("Test error badSrcParent on os.Stat with non-nil fi")
	}
	if strings.HasSuffix(name, `prj\symlink\symlinkdir\`) {
		fi, _ := os.Stat(".")
		return fi, fmt.Errorf("readlink for symlinkdir")
	}
	if strings.HasSuffix(name, `.1`) {
		fi, _ := os.Stat(".")
		return fi, nil
	}
	return os.Stat(name)
}

func testOsMkdirAll(path string, perm os.FileMode) error {
	fmt.Printf("testOsMkdirAll path='%+v'\n", path)
	if strings.HasSuffix(path, `badSrcParentMdirAll\`) {
		return fmt.Errorf("Error on mkDirAll for '%s'", path)
	}
	return nil
}

var junctionOut = `
 Répertoire de C:\Users\VonC\prog\git\ggb\deps\src\github.com\VonC

22/06/2015  11:03    <REP>          .
22/06/2015  11:03    <REP>          ..
22/06/2015  11:03    <JONCTION>     symlink [C:\Users\VonC\prog\git\ggb\]
22/06/2015  11:03    <JONCTION>     symlinkdir [C:\Users\VonC\prog\git\ggb\]`

func testExecRun(cmd *exec.Cmd) error {
	tmsg := fmt.Sprintf("testExecRun cmd='%v' in '%s'", cmd.Args, cmd.Dir)
	fmt.Println(tmsg)
	if strings.Contains(tmsg, "/J") {
		return nil
	}
	if strings.HasSuffix(cmd.Dir, `\WarningOnDir`) {
		io.WriteString(cmd.Stdout, "dummy content")
		io.WriteString(cmd.Stderr, "Some warning on dir")
		return nil
	}
	if strings.HasSuffix(cmd.Dir, `\nojunction`) {
		io.WriteString(cmd.Stdout, "dummy content without any junction")
		return nil
	}
	if strings.HasSuffix(cmd.Dir, `\badsymlink`) {
		return fmt.Errorf("unreadable dir on symlink")
	}
	if strings.HasSuffix(cmd.Dir, `\symlink`) {
		io.WriteString(cmd.Stdout, junctionOut)
		return nil
	}
	return cmdRun(cmd)
}

func testOsRename(oldpath, newpath string) error {
	fmt.Printf("testOsRename oldpath='%v', newpath '%s'\n", oldpath, newpath)
	return nil
}
