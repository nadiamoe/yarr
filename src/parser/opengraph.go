package parser

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/content/htmlutil"
	"golang.org/x/net/html"
)

type OpenGraph map[string]string

func ParseOpenGraph(url string) (OpenGraph, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building http request: %w", err)
	}

	req.Header.Add("user-agent", "Yarr/1.0 (Yarr Fallback OpenGraph parser; +https://github.com/nadiamoe/yarr)")

	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   3 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making http request: %w", err)
	}

	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	og := OpenGraph{}

	for _, meta := range htmlutil.FindNodes(root, isOpenGraph) {
		propIndex := slices.IndexFunc(meta.Attr, func(a html.Attribute) bool { return a.Key == "property" })
		if propIndex == -1 {
			continue
		}
		prop := meta.Attr[propIndex]
		if !strings.HasPrefix(prop.Val, "og:") {
			continue
		}

		contentIndex := slices.IndexFunc(meta.Attr, func(a html.Attribute) bool { return a.Key == "content" })
		if contentIndex == -1 {
			continue
		}

		og[strings.TrimPrefix(prop.Val, "og:")] = meta.Attr[contentIndex].Val
	}

	return og, nil
}

func isOpenGraph(n *html.Node) bool {
	if n.Data != "meta" {
		return false
	}

	if n.Parent.Data != "head" {
		return false
	}

	if n.Parent.Parent.Data != "html" {
		return false
	}

	return true
}

func FallbackOpenGraph(thumbnail, url, what string) string {
	if thumbnail != "" {
		return thumbnail
	}

	og, err := ParseOpenGraph(url)
	if err != nil {
		return ""
	}

	return og[what]
}
