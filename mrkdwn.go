package mrkdwn

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
)

type Mrkdwn struct {
	text string
	image []string
}

func FromHTML(html string) (*Mrkdwn, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return nil, err
	}
	return &Mrkdwn{markdown, nil}, nil
}
