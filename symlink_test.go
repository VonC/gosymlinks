package gosymlink

import (
	"strings"
	"testing"
)

type testDir struct {
	dirpath   string
	exist     bool
	issymlink bool
	destpath  string
	errormsg  string
}

var testsDir = []testDir{
	// Hollidays() would return 0 if no holidays.txt file
	testDir{dirpath: "nonexistentdir"},
}

func TestDir(t *testing.T) {
	for _, tst := range testsDir {
		path := "tests/" + tst.dirpath
		dir, err := DirFrom(path)
		checkBool(t, dir.exist, tst.exist, "Dir.exist ("+tst.dirpath+")")
		checkBool(t, dir.issymlink, tst.issymlink, "Dir.issymlink ("+tst.dirpath+")")
		checkString(t, dir.destpath, tst.destpath, "Dir.destpath ("+tst.dirpath+")")
		checkErrorMsg(t, err, tst.errormsg, tst.dirpath)
	}
}

func checkBool(t *testing.T, b bool, expected bool, id string) {
	if b != expected {
		t.Errorf("%s:\nExpected:\n%s',\nBUT got\n%s'", id, expected, b)
	}
}

func checkString(t *testing.T, s string, expected string, id string) {
	if strings.Contains(s, expected) == false {
		t.Errorf("%s:\nExpected:\n%s',\nBUT got\n%s'", id, expected, s)
	}
}

func checkErrorMsg(t *testing.T, err error, errormsg string, id string) {
	if errormsg == "" && err != nil {
		t.Errorf("%s:\nDid not Expect an error, but got '%s'", id, err.Error())
	}
	if errormsg != "" {
		if err == nil {
			t.Errorf("%s:\nExpected an error", id)
		} else {
			errmsgs := strings.Split(errormsg, "\n")
			expected := ""
			for _, errmsg := range errmsgs {
				errmsg = strings.TrimSpace(errmsg)
				if strings.Contains(err.Error(), errmsg) == false {
					expected = expected + errmsg + "\n"
				}
			}
			expected = strings.TrimSpace(expected)
			if expected != "" {
				t.Errorf("%s:\nExpected:\n%s',\nBUT got\n%s'", id, expected, err.Error())
			}
		}
	}

}
