package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/a8m/mark"
	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"golang.org/x/net/html"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			printToStderr(r.(error).Error())
			os.Exit(1)
		}
	}()

	args := getArgs()

	bs, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		panic(err)
	}

	n, err := html.Parse(strings.NewReader(mark.Render(string(bs))))

	if err != nil {
		panic(err)
	}

	ok := true

	for s := range extractURLs(n) {
		colored := color.New(color.FgCyan).SprintFunc()(s)

		if _, err := http.Get(s); err != nil {
			printToStderr(color.New(color.FgRed).SprintFunc()("ERROR") + "\t" + colored + "\t" + err.Error())
			ok = false
		} else if err == nil && args["--verbose"].(bool) {
			printToStderr(color.New(color.FgGreen).SprintFunc()("OK") + "\t" + colored)
		}
	}

	if !ok {
		os.Exit(1)
	}
}

func extractURLs(n *html.Node) map[string]bool {
	ss := make(map[string]bool)
	ns := make([]*html.Node, 0, 1024)
	ns = append(ns, n)

	for len(ns) > 0 {
		i := len(ns) - 1
		n := ns[i]
		ns = ns[:i]

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && isURL(a.Val) {
					ss[a.Val] = true
					break
				}
			}
		}

		for n := n.FirstChild; n != nil; n = n.NextSibling {
			ns = append(ns, n)
		}
	}

	return ss
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func getArgs() map[string]interface{} {
	usage := `Link checker for Markdown and HTML

Usage:
	linkcheck [-v] <filename>

Options:
	-v, --verbose  Be verbose`

	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		panic(err)
	}

	return args
}

func printToStderr(xs ...interface{}) {
	fmt.Fprintln(os.Stderr, xs...)
}
