# html-to-mrkdwn-go

[![Test](https://github.com/matsubara0507/html-to-mrkdwn-go/actions/workflows/test.yaml/badge.svg)](https://github.com/matsubara0507/html-to-mrkdwn-go/actions/workflows/test.yaml)

Convert HTML to Slack's [mrkdwn](https://api.slack.com/docs/message-formatting) format.

```go
package main

import (
	"fmt"

	mrkdwn "github.com/matsubara0507/html-to-mrkdwn-go"
)

func main() {
	html := `
		<p><strong>Hello</strong> <a href="https://example.com">cruel</a> <em>world</em>!</p>
		<p><img src="https://media.giphy.com/media/5xtDarEbygs3Pu7p3jO/giphy.gif"></p>
	`
	md, err := mrkdwn.FromHTML(html)
	if err != nil {
		panic(err)
	}
	fmt.Println(md.Text)
	// *Hello*<https://example.com|cruel>  _world_!
	//
	// https://media.giphy.com/media/5xtDarEbygs3Pu7p3jO/giphy.gif
}
```

This package is greatly inspired by [html-to-mrkdwn](https://www.npmjs.com/package/html-to-mrkdwn) npm package.
