package scopes

import (
	"bytes"
	"github.com/quarnster/completion/content"
	"github.com/quarnster/parser"
)

// Returns the RangeSet visible from the given source code location.
// Note that it only uses '{' and '}' to deduce scopes, and further
// processing is likely needed depending on language specifics.
func Visibility(loc content.SourceLocation) (ret parser.RangeSet) {
	var s SCOPES
	s.Parse(loc.File.Contents)
	var rec func(node *parser.Node)
	pos := int(loc.Offset())
	rec = func(node *parser.Node) {
		for _, child := range node.Children {
			if !child.Range.Inside(pos) || child.Name == "TextScope" || child.Name == "CommentScope" {
				r := child.Range
				if child.Name == "BracketScope" || child.Name == "TextScope" {
					// We want the brackets (or "'s) as part of the result to make the resulting
					// source code parseable.
					r.Start += 1
					r.End -= 1
				}
				ret = ret.Cut(r)
			} else {
				rec(child)
			}
		}
	}
	root := s.RootNode()
	ret.Add(root.Range)
	rec(root)
	return
}

func Substr(data string, visible parser.RangeSet) string {
	b := bytes.NewBuffer(nil)
	for _, r := range visible {
		b.WriteString(data[r.Start:r.End])
	}
	return b.String()
}
