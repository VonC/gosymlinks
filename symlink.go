package gosymlink

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Dir struct {
	dirpath   string
	exist     bool
	issymlink bool
	destpath  string
}

func DirFrom(path string) (*Dir, error) {
	res := &Dir{}
	fi, err := os.Stat(path)
	fmt.Printf("%s=> %+v err %+v\n", path, fi, err)
	if fi == nil {
		if strings.Contains(err.Error(), "The system cannot find the") {
			return res, nil
		}
		return nil, err
	}
	if fi.IsDir() == false {
		return nil, fmt.Errorf("%s is not a folder, it is a file", path)
	}
	return res, err
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
