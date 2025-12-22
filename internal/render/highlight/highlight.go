package highlight

import (
	"bytes"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Extender struct{}

func (e *Extender) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{}, 200),
		),
	)
}

type Renderer struct{}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.Render)
}

func (r *Renderer) Render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.FencedCodeBlock)
	lang := string(n.Language(source))

	var buf bytes.Buffer
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buf.Write(line.Value(source))
	}
	code := buf.String()

	var lexer chroma.Lexer
	if lang != "" {
		lexer = lexers.Get(lang)
	}
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	langAttr := lang
	if langAttr == "" {
		langAttr = lexer.Config().Name
		if langAttr == "" {
			langAttr = "plaintext"
		}
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return ast.WalkStop, err
	}

	formatter := html.New(
		html.WithClasses(true),
		html.PreventSurroundingPre(true),
	)

	var codeBuf bytes.Buffer
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	if err := formatter.Format(&codeBuf, style, iterator); err != nil {
		return ast.WalkStop, err
	}

	w.WriteString(`<div class="code-block" data-lang="` + langAttr + `">`)
	w.WriteString(`<button class="copy-btn" aria-label="Copy code">Copy</button>`)
	w.WriteString(`<pre class="chroma">`)
	w.WriteString(`<code class="language-` + langAttr + `">`)
	w.Write(codeBuf.Bytes())
	w.WriteString(`</code>`)
	w.WriteString(`</pre>`)
	w.WriteString(`</div>`)

	return ast.WalkSkipChildren, nil
}
