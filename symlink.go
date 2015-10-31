package symlink

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type SL struct {
	dst  string
	path string
}

func New(link, dst string) (*SL, error) {
	var err error
	if dst, err = dirAbsPath(dst); err != nil {
		return nil, err
	}
	exist, _, err := dirExists(dst)
	msgerr := ""
	if err != nil {
		msgerr = fmt.Sprintf("\rError: '%+v'", err)
	}
	if !exist {
		return nil, fmt.Errorf("Unknown destination '%s'%s", dst, msgerr)
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

var osStat = os.Stat

func dirExists(path string) (bool, string, error) {
	fi, err := osStat(path)
	if fi == nil {
		return false, "", err
	}
	// sys := fi.Sys().(*syscall.Win32FileAttributeData)
	if err == nil {
		return true, "", nil
	}
	if strings.HasPrefix(err.Error(), "readlink ") == false {
		return false, "", err
	}
	// This is a symlink (JUNCTION on Windows)
	dir := filepath.Dir(path)
	base := filepath.Base(dir)
	dir = filepath.Dir(dir)
	sdir := ""
	if sdir, err = execcmd("dir", ".", dir); err != nil {
		return false, "", err
	}
	r := regexp.MustCompile(fmt.Sprintf(`(?m)<J[UO]NCTION>\s+%s\s+\[([^\]]+)\]\s*$`, base))
	n := r.FindAllStringSubmatch(sdir, -1)
	// fmt.Printf("n='%+v'\nr='%+v'\n", n, r)
	if len(n) == 1 {
		return true, n[0][1], nil
	}
	return false, "", fmt.Errorf("Unable to find junction symlink in parent dir '%s' for '%s'", dir, base)
}

func cmdRun(cmd *exec.Cmd) error {
	return cmd.Run()
}

var execRun = cmdRun

func execcmd(exe, cmd string, dir string) (string, error) {
	args := strings.Split(cmd, " ")
	args = append([]string{"/c", exe}, args...)
	c := exec.Command("cmd", args...)
	c.Dir = dir
	var bout bytes.Buffer
	c.Stdout = &bout
	var berr bytes.Buffer
	c.Stderr = &berr
	err := execRun(c)
	if err != nil {
		return bout.String(), fmt.Errorf("Unable to run '%s %s' in '%s': err '%s'\n'%s'", exe, cmd, dir, err.Error(), berr.String())
	} else if berr.String() != "" {
		return bout.String(), fmt.Errorf("Warning on run '%s %s' in '%s': '%s'", exe, cmd, dir, berr.String())
	}
	return bout.String(), nil
}
