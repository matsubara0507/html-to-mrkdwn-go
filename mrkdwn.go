package mrkdwn

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
)

type Mrkdwn struct {
	Text  string
	Image string
}

func FromHTML(html string) (*Mrkdwn, error) {
	converter := md.NewConverter("", false, nil)
	converter.Use(SlackPlugin())
	converter.Use(plugin.Strikethrough("~"))
	converter.Remove("table")

	markdown, err := converter.ConvertString(html)
	if err != nil {
		return nil, err
	}

	image, err := FirstImage(html)
	if err != nil {
		return nil, err
	}

	return &Mrkdwn{markdown, image}, nil
}

func SlackPlugin() md.Plugin {
	return func(conv *md.Converter) (rules []md.Rule) {
		rules = CommonmarkRules()
		rules = append(rules, SlackLinkRule())
		rules = append(rules, SlackHeadingRule())
		rules = append(rules, SlackCheckbox())
		rules = append(rules, SlackListItemRule())
		rules = append(rules, SlackImagesRule())
		return
	}
}

// ref: html-to-markdown/commonmark.go
func CommonmarkRules() (rules []md.Rule) {
	rules = append(rules, md.Rule{
		Filter: []string{"#text"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			text := SlackEscape(selec.Text())
			if trimmed := strings.TrimSpace(text); trimmed == "" {
				return md.String("")
			}

			return &text
		},
	})
	rules = append(rules, md.Rule{
		Filter: []string{"p", "div"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			parent := goquery.NodeName(selec.Parent())
			if md.IsInlineElement(parent) || parent == "li" {
				content = "\n" + content + "\n"
				return &content
			}

			// remove unnecessary spaces to have clean markdown
			content = md.TrimpLeadingSpaces(content)

			content = "\n" + content + "\n"
			return &content
		},
	})
	rules = append(rules, md.Rule{
		Filter: []string{"blockquote"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			content = strings.TrimSpace(content)
			if content == "" {
				return nil
			}
			var beginningR = regexp.MustCompile(`(?m)^`)
			content = beginningR.ReplaceAllString(content, "> ")

			text := "\n\n" + content + "\n\n"
			return &text
		},
	})
	rules = append(rules, md.Rule{
		Filter: []string{"strong", "b"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			// only use one bold tag if they are nested
			parent := selec.Parent()
			if parent.Is("strong") || parent.Is("b") {
				return &content
			}

			trimmed := strings.TrimSpace(content)
			if trimmed == "" {
				return &trimmed
			}
			trimmed = "*" + trimmed + "*" // without StrongDelimiter for validateOptions

			// always have a space to the side to recognize the delimiter
			trimmed = md.AddSpaceIfNessesary(selec, trimmed)

			return &trimmed
		},
	})
	rules = append(rules, md.Rule{
		Filter: []string{"i", "em"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			// only use one italic tag if they are nested
			parent := selec.Parent()
			if parent.Is("i") || parent.Is("em") {
				return &content
			}

			trimmed := strings.TrimSpace(content)
			if trimmed == "" {
				return &trimmed
			}
			trimmed = opt.EmDelimiter + trimmed + opt.EmDelimiter

			// always have a space to the side to recognize the delimiter
			trimmed = md.AddSpaceIfNessesary(selec, trimmed)

			return &trimmed
		},
	})
	rules = append(rules, md.Rule{
		Filter: []string{"pre"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			codeElement := selec.Find("code")
			language := codeElement.AttrOr("class", "")
			language = strings.Replace(language, "language-", "", 1)

			code := codeElement.Text()
			if codeElement.Length() == 0 {
				code = selec.Text()
			}

			fenceChar, _ := utf8.DecodeRuneInString(opt.Fence)
			fence := md.CalculateCodeFence(fenceChar, code)

			text := "\n\n" + fence + language + "\n" +
				code +
				"\n" + fence + "\n\n"
			return &text
		},
	})
	rules = append(rules, md.Rule{
		Filter: []string{"code"},
		Replacement: func(_ string, selec *goquery.Selection, opt *md.Options) *string {
			code := selec.Text()
			// code block contains a backtick as first character
			if strings.HasPrefix(code, "`") {
				code = " " + code
			}
			// code block contains a backtick as last character
			if strings.HasSuffix(code, "`") {
				code = code + " "
			}

			// TODO: configure delimeter in options?
			text := "`" + code + "`"
			return &text
		},
	})
	return
}

func SlackEscape(text string) string {
	text = strings.Replace(text, "&", "&amp;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	return text
}

func SlackLinkRule() md.Rule {
	return md.Rule{
		Filter: []string{"a"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			href, ok := selec.Attr("href")
			if !ok {
				return &content
			}
			text := fmt.Sprintf("<%s|%s> ", href, content)
			return &text
		},
	}
}

func SlackHeadingRule() md.Rule {
	return md.Rule{
		Filter: []string{"h1", "h2", "h3", "h4", "h5", "h6"},
		Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
			return md.String("*" + content + "*\n")
		},
	}
}

func SlackCheckbox() md.Rule {
	return md.Rule{
		Filter: []string{"input"},
		Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
			if selec.HasClass("task-list-item-checkbox") {
				if _, ok := selec.Attr("checked"); ok {
					return md.String("☑︎")
				} else {
					return md.String("☐")
				}
			}
			return &content
		},
	}
}

func SlackListItemRule() md.Rule {
	return md.Rule{
		Filter: []string{"li"},
		Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
			// ToDo: replace newline
			prefix := "• " // without BulletListMarker for validateOptions
			parent := selec.Parent()
			if selec.HasClass("task-list-item") {
				prefix = ""
			}
			if parent != nil && goquery.NodeName(parent) == "ol" {
				index := selec.Index()
				if start, ok := parent.Attr("start"); ok {
					startIndex, _ := strconv.Atoi(start)
					prefix = strconv.Itoa(startIndex+index) + ". "
				} else {
					prefix = strconv.Itoa(index+1) + ". "
				}
			}
			suffix := "\n"
			if nil == selec.Next() {
				suffix = ""
			}
			return md.String(prefix + strings.Trim(content, " \t\n") + suffix)
		},
	}
}

func SlackImagesRule() md.Rule {
	return md.Rule{
		Filter: []string{"img"},
		Replacement: func(content string, selec *goquery.Selection, options *md.Options) *string {
			parent := selec.Parent()
			if parent != nil && goquery.NodeName(parent) == "a" {
				if alt, ok := selec.Attr("alt"); ok {
					return &alt
				}
				if src, ok := selec.Attr("src"); ok {
					return &src
				}
				return nil
			}
			if alt, ok := selec.Attr("alt"); ok {
				src, _ := selec.Attr("src")
				return md.String(fmt.Sprintf("<%s|%s> ", src, alt))
			}
			if src, ok := selec.Attr("src"); ok {
				return md.String(src + " ")
			}
			return nil
		},
	}
}

func FirstImage(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	if img := doc.Find("img"); img != nil {
		if src, ok := img.Attr("src"); ok {
			return src, nil
		}
	}
	return "", nil
}
