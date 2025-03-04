package domtree

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func elemAttr(name, val string) html.Attribute {
	return html.Attribute{Key: name, Val: val}
}

func classAttr(val string) html.Attribute {
	return html.Attribute{Key: "class", Val: val}
}

func appendAttr(node *html.Node, name, val string) {
	node.Attr = append(node.Attr, elemAttr(name, val))
}

func removeAttr(node *html.Node, name string) {
	attrs := []html.Attribute{}
	for _, attr := range node.Attr {
		if attr.Key != name {
			attrs = append(attrs, attr)
		}
	}
	node.Attr = attrs
}

func setId(node *html.Node, id string) {
	removeAttr(node, "id")
	appendAttr(node, "id", id)
}

func appendClass(node *html.Node, val string) {
	node.Attr = append(node.Attr, classAttr(val))
}

func isHeadingElem(node *html.Node) bool {
	return node.Data == "h1" || node.Data == "h2" || node.Data == "h3" || node.Data == "h4" || node.Data == "h5" || node.Data == "h6"
}

func findElements(node *html.Node, tags []string) []*html.Node {
	var elements []*html.Node
	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, tag := range tags {
				if n.Data == tag {
					elements = append(elements, n)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(node)
	return elements
}

func getAttribute(node *html.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

func getTextContent(node *html.Node) string {
	var buf strings.Builder
	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(node)
	return buf.String()
}

func isBlankNode(n *html.Node) bool {
	t := getTextContent(n)
	return len(strings.TrimSpace(t)) < 1
}

func newTextNode(data string) *html.Node {
	n := &html.Node{
		Type: html.TextNode,
		Data: data,
	}
	return n
}

func newElementNode(data string, dataAtom atom.Atom) *html.Node {
	n := &html.Node{
		Type:     html.ElementNode,
		Data:     data,
		DataAtom: dataAtom,
	}
	return n
}

func newHeadNode() *html.Node {
	return newElementNode("head", atom.Head)
}

func newDivNode() *html.Node {
	return newElementNode("div", atom.Div)
}

func newANode() *html.Node {
	return newElementNode("a", atom.A)
}

func newUlNode() *html.Node {
	return newElementNode("ul", atom.Ul)
}

func newLiNode() *html.Node {
	return newElementNode("li", atom.Li)
}

func newLinkNode() *html.Node {
	return newElementNode("link", atom.Link)
}

func newInputNode() *html.Node {
	return newElementNode("input", atom.Input)
}

func newCheckboxInputNode(disabled bool, filled bool) *html.Node {
	n := newInputNode()
	appendAttr(n, "type", "checkbox")
	if disabled {
		appendAttr(n, "disabled", "true")
	}
	if filled {
		appendAttr(n, "checked", "true")
	}
	return n
}

func newScriptNode() *html.Node {
	return newElementNode("script", atom.Script)
}

func newStyleNode() *html.Node {
	return newElementNode("style", atom.Style)
}

func NewHtmlNode(lang string) *html.Node {
	n := newElementNode("html", atom.Html)
	appendAttr(n, "lang", lang)
	return n
}

func Decode(node *html.Node) string {
	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>")
	html.Render(&buf, node)
	return buf.String()
}

func getLastModTime(src string) string {
	fi, err := os.Stat(src)
	if err != nil {
		return ""
	}
	return fi.ModTime().Format("2006-01-02")
}

func newTimestampNode(src string) *html.Node {
	d := newDivNode()
	setId(d, "timestamp")
	m := getLastModTime(src)
	d.AppendChild(newTextNode(fmt.Sprintf("update: %s", m)))
	return d
}
