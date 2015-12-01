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
	testDir{"nonexistentdir", false, false, "", ""},
}

func TestDir(t *testing.T) {
	for _, tst := range testsDir {
		path := "tests/" + tst.dirpath
		_, err := DirFrom(path)
		checkErrorMsg(t, err, tst.errormsg, tst.dirpath)
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
