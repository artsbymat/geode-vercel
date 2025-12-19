package render

import (
	"bytes"
	"geode/internal/types"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type headingShiftAndTocTransformer struct {
	shift int
	toc   *[]types.TocItem
}

func (t *headingShiftAndTocTransformer) Transform(node *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		if t.shift != 0 {
			newLevel := h.Level + t.shift
			if newLevel < 1 {
				newLevel = 1
			} else if newLevel > 6 {
				newLevel = 6
			}
			h.Level = newLevel
		}

		if t.toc != nil {
			text := headingText(h, source)
			id := headingID(h)
			if strings.TrimSpace(text) != "" && strings.TrimSpace(id) != "" {
				*t.toc = append(*t.toc, types.TocItem{Level: h.Level, Text: text, ID: id})
			}
		}

		return ast.WalkContinue, nil
	})
}

func withHeadingShiftAndTOC(shift int, toc *[]types.TocItem) parser.Option {
	return parser.WithASTTransformers(
		util.Prioritized(&headingShiftAndTocTransformer{shift: shift, toc: toc}, 100),
	)
}

func headingID(h *ast.Heading) string {
	if v, ok := h.AttributeString("id"); ok {
		switch vv := v.(type) {
		case []byte:
			return string(vv)
		case string:
			return vv
		}
	}
	return ""
}

func headingText(h *ast.Heading, source []byte) string {
	var b bytes.Buffer
	ast.Walk(h, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch nn := n.(type) {
		case *ast.Text:
			b.Write(nn.Value(source))
		case *ast.String:
			b.Write(nn.Value)
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(b.String())
}
