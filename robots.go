package main

import (
	"fmt"
	"os"
	"reisen/config"
	"strings"
)

//Generates robots.txt
func robots(conf config.Config) {
	//We build the robots.txt as an in-memory string beforehand
	//so we can minimize IO errors we need to handle
	var b strings.Builder
	b.WriteString("User-agent: *\n")

	for _, board := range conf.Boards {
		name := board.Name

		fmt.Fprintf(&b, "Disallow: /%s/search\n", name)
		fmt.Fprintf(&b, "Disallow: /%s/view-same\n", name)
		fmt.Fprintf(&b, "Disallow: /%s/post\n", name)
		fmt.Fprintf(&b, "Disallow: /%s?rkeyset=\n", name)
		fmt.Fprintf(&b, "Disallow: /%s/report", name)
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
