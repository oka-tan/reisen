package main

import (
	"fmt"
	"os"
	"reisen/config"
	"strings"
)

func robots(conf config.Config) {
	var b strings.Builder
	b.WriteString("User-agent: *\n")

	for _, board := range conf.Boards {
		name := board.Name

		fmt.Fprintf(&b, "/%s/search\n", name)
		fmt.Fprintf(&b, "/%s/view-same\n", name)
		fmt.Fprintf(&b, "/%s/post\n", name)
		fmt.Fprintf(&b, "/%s?rkeyset=\n", name)
	}

	robotsTxt, err := os.Create("static/robots.txt")
	if err != nil {
		panic(err)
	}

	if _, err := robotsTxt.WriteString(b.String()); err != nil {
		panic(err)
	}

	robotsTxt.Close()
}
