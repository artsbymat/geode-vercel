package mark

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Extender struct{}

func (e *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&MarkParser{}, 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&MarkRenderer{}, 500),
		),
	)
}

var KindMark = ast.NewNodeKind("Mark")

type Mark struct {
	ast.BaseInline
}

func (m *Mark) Dump(source []byte, level int) {
	ast.DumpHelper(m, source, level, nil, nil)
}

func (m *Mark) Kind() ast.NodeKind {
	return KindMark
}

func NewMark() *Mark {
	return &Mark{}
}

type MarkParser struct{}

func (s *MarkParser) Trigger() []byte {
	return []byte{'='}
}

func (s *MarkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()

	if !shouldParseMark(parent) {
		return nil
	}

	if len(line) < 4 {
		return nil
	}

	if line[0] != '=' || line[1] != '=' {
		return nil
	}

	i := 2
	for i < len(line)-1 {
		if line[i] == '=' && line[i+1] == '=' {

			if i == 2 {
				return nil
			}

			node := NewMark()

			contentSegment := text.NewSegment(segment.Start+2, segment.Start+i)
			node.AppendChild(node, ast.NewTextSegment(contentSegment))

			block.Advance(i + 2)
			return node
		}
		i++
	}

	return nil
}

func shouldParseMark(parent ast.Node) bool {
	current := parent
	for current != nil {
		kind := current.Kind()

		if kind == ast.KindCodeSpan ||
			kind == ast.KindCodeBlock ||
			kind == ast.KindFencedCodeBlock ||
			kind == ast.KindHTMLBlock ||
			kind == ast.KindRawHTML {
			return false
		}

		current = current.Parent()
	}

	return true
}

type MarkRenderer struct{}

func (r *MarkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMark, r.renderMark)
}

func (r *MarkRenderer) renderMark(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<mark>")
	} else {
		_, _ = w.WriteString("</mark>")
	}
	return ast.WalkContinue, nil
}
