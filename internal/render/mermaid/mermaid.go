package mermaid

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
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{}, 2000),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{}, 1000),
		),
	)
}

type Transformer struct{}

func (t *Transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindFencedCodeBlock {
			fenced := n.(*ast.FencedCodeBlock)
			lang := string(fenced.Language(reader.Source()))
			if lang == "mermaid" {
				pc.Set(mermaidKey, true)

				mermaidBlock := &MermaidBlock{
					BlockLine: fenced,
				}
				parent := n.Parent()
				parent.ReplaceChild(parent, n, mermaidBlock)
			}
		}
		return ast.WalkContinue, nil
	})
}

var mermaidKey = parser.NewContextKey()

func GetHasMermaid(pc parser.Context) bool {
	return pc.Get(mermaidKey) == true
}

type MermaidBlock struct {
	ast.BaseBlock
	BlockLine *ast.FencedCodeBlock
}

func (n *MermaidBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{}, nil)
}

var KindMermaidBlock = ast.NewNodeKind("MermaidBlock")

func (n *MermaidBlock) Kind() ast.NodeKind {
	return KindMermaidBlock
}

type Renderer struct{}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMermaidBlock, r.Render)
}

func (r *Renderer) Render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("<div class=\"mermaid\">")
		n := node.(*MermaidBlock)
		lines := n.BlockLine.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			w.Write(line.Value(source))
		}
		w.WriteString("</div>")
	}
	return ast.WalkContinue, nil
}
