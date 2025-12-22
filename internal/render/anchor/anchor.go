package anchor

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type Extender struct{}

func (e *Extender) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&HeadingRenderer{}, 100),
		),
	)
}

type HeadingRenderer struct {
	html.Config
}

func (r *HeadingRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHeading, r.renderHeading)
}

func (r *HeadingRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)

	if entering {
		_, _ = w.WriteString("<h")
		_ = w.WriteByte("0123456"[n.Level])

		if n.Attributes() != nil {
			id, ok := n.AttributeString("id")
			if ok {
				_, _ = w.WriteString(" id=\"")
				if idStr, ok := id.([]byte); ok {
					_, _ = w.Write(idStr)
				} else if idStr, ok := id.(string); ok {
					_, _ = w.WriteString(idStr)
				}
				_, _ = w.WriteString("\"")
			}
		}

		_, _ = w.WriteString(">")

		if n.Attributes() != nil {
			id, ok := n.AttributeString("id")
			if ok {
				_, _ = w.WriteString("<a href=\"#")
				if idStr, ok := id.([]byte); ok {
					_, _ = w.Write(idStr)
				} else if idStr, ok := id.(string); ok {
					_, _ = w.WriteString(idStr)
				}
				_, _ = w.WriteString("\" class=\"anchor-heading\" aria-hidden=\"true\">#</a>")
			}
		}
	} else {
		_, _ = w.WriteString("</h")
		_ = w.WriteByte("0123456"[n.Level])
		_, _ = w.WriteString(">\n")
	}

	return ast.WalkContinue, nil
}
