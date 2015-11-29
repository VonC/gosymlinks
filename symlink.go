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

var osMkdirAll = os.MkdirAll

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

	if link, err = dirAbsPath(link); err != nil {
		return nil, err
	}
	linkdir := filepath.Dir(filepath.Dir(link)) + string(filepath.Separator)
	// fmt.Printf("link='%+v\nlinkdir=%+v\ndst=%+v\n", link, linkdir, dst)
	var hasLinkDir, hasLink bool
	var linkDirTarget, linkTarget string
	if hasLinkDir, linkDirTarget, err = dirExists(linkdir); err != nil {
		if strings.Contains(err.Error(), "The system cannot find the") == false {
			return nil, fmt.Errorf("Impossible to check/access link parent folder '%s':\n'%+v'", linkdir, err)
		}
	}
	if linkDirTarget != "" {
		// move folder to x.1 (or error)
		if err = moveToDotX(linkdir); err != nil {
			return nil, err
		}
		hasLinkDir = false
	}
	if !hasLinkDir {
		if err = osMkdirAll(linkdir, os.ModeDir); err != nil {
			return nil, fmt.Errorf("Impossible to create link parent folder '%s':\n'%+v'", linkdir, err)
		}
	}
	// fmt.Printf("==== link '%s'\n", link)
	if hasLink, linkTarget, err = dirExists(link); err != nil {
		if strings.Contains(err.Error(), "The system cannot find the") == false {
			return nil, fmt.Errorf("Impossible to check/access link'%s':\n'%+v'", link, err)
		}
	}
	if linkTarget == dst {
		return &SL{path: link, dst: dst}, nil
	}
	if hasLink {
		fmt.Printf("=> linkTarget='%s' vs. dst='%s'\n", linkTarget, dst)
		if err = moveToDotX(link); err != nil {
			return nil, err
		}
	}
	if _, err = execcmd("mklink", fmt.Sprintf("/J %s %s", link, dst), linkdir); err != nil {
		return nil, fmt.Errorf("Impossible to create junction between '%s' and '%s':\n'%+v'", link, dst, err)
	}
	res := &SL{path: link, dst: dst}
	return res, nil
}

var osRename = os.Rename

func moveToDotX(path string) error {
	path = filepath.Dir(path)
	i := 1
	newpath := fmt.Sprintf("%s.%d", path, i)
	for {

		exist, _, err := dirExists(newpath)
		if err != nil {
			if strings.Contains(err.Error(), "The system cannot find the") == false {
				return err
			}
		}
		if !exist {
			break
		}
		i = i + 1
		newpath = fmt.Sprintf("%s.%d", path, i)
	}
	fmt.Printf("Move '%s' to '%s'\n", path, newpath)
	if err := osRename(path, newpath); err != nil {
		return err
	}
	return nil
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
	// fmt.Printf("==== fi='%+v', err='%+v'\n", fi, err)
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
		linkTarget := n[0][1]
		if strings.HasSuffix(linkTarget, string(filepath.Separator)) == false {
			linkTarget = linkTarget + string(filepath.Separator)
		}
		return true, linkTarget, nil
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
