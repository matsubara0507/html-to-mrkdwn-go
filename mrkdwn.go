package mrkdwn

import (
	"bytes"

	"github.com/mattn/godown"
)

type Mrkdwn struct {
	text []byte
	image []string
}

func FromHTML(html []byte) (*Mrkdwn, error) {
	var buf bytes.Buffer
	err := godown.Convert(&buf, bytes.NewBuffer(html), nil)
	if err != nil {
		return nil, err
	}
	return &Mrkdwn{buf.Bytes(), nil}, nil
}