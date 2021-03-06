package mrkdwn

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestFromHTML(t *testing.T) {
	fixturePaths := make([]string, 0)
	for _, fixturesDir := range []string{"html-to-mrkdwn/test/fixtures/*.mrkdwn", "fixtures/*.mrkdwn"} {
		paths, err := filepath.Glob(fixturesDir)
		if err != nil {
			panic(err)
		}
		fixturePaths = append(fixturePaths, paths...)
	}
	for _, path := range fixturePaths {
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
			md, err := FromHTML(fs[0])
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			actual := strings.Trim(md.Text, " \t\n")
			expect := strings.Trim(fs[1], " \t\n")
			if !strings.EqualFold(actual, expect) {
				t.Errorf("\nexpected:\n%s\nbut actual:\n%s", expect, actual)
			}
		})
	}

	html := `
		<p><strong>Hello</strong> <a href="https://example.com">cruel</a> <em>world</em>!</p>
		<p><img src="https://example.com/first.gif"></p>
		<p><img src="https://example.com/secon.gif"></p>
	`
	t.Run("return images", func(t *testing.T) {
		md, err := FromHTML(html)
		if err != nil {
			t.Errorf("%v", err)
			return
		}
		expect := "https://example.com/first.gif"
		if !strings.EqualFold(md.Image, expect) {
			t.Errorf("expected %s, but acutual %s", expect, md.Image)
		}
	})
}
