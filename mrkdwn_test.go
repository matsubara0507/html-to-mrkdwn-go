package mrkdwn

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestFromHTML(t *testing.T) {
	paths, err := filepath.Glob("html-to-mrkdwn/test/fixtures/*.mrkdwn")
	if err != nil {
		panic(err)
	}
	for _, path := range paths {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			_ = fmt.Errorf("can't read fixture: %s, %v", path, err)
			continue
		}
		fs := strings.Split(string(file), "====")
		if len(fs) != 2 {
			_ = fmt.Errorf("fixture is not formated: %s", path)
			continue
		}
		t.Run(path, func(t *testing.T) {
			md, err := FromHTML([]byte(fs[0]))
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			actual := strings.Trim(string(md.text), " \t\n")
			expect := strings.Trim(fs[1], " \t\n")
			if !strings.EqualFold(actual, expect) {
				t.Errorf("\nexpected:\n%s\nbut actual:\n%s", expect, actual)
			}
		})
	}
}
