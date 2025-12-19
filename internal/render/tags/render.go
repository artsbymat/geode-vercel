package hashtag

import (
	"fmt"
	"sync"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Resolver interface {
	ResolveHashtag(*Node) (destination []byte, err error)
}

type Renderer struct {
	Resolver   Resolver
	Collector  *Collector
	Attributes []Attribute
	hasDest    sync.Map
}

type Attribute struct {
	Name  string
	Value string
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)
}

func (r *Renderer) Render(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, ok := node.(*Node)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *Node", node)
	}

	if entering {
		if err := r.enter(w, n); err != nil {
			return ast.WalkStop, err
		}
	} else {
		r.exit(w, n)
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) enter(w util.BufWriter, n *Node) error {
	_, _ = w.WriteString(`<span class="hashtag">`)

	if r.Collector != nil {
		r.Collector.Add(n.Tag)
	}

	var dest []byte
	if res := r.Resolver; res != nil {
		var err error
		dest, err = res.ResolveHashtag(n)
		if err != nil {
			return fmt.Errorf("resolve hashtag %q: %w", n.Tag, err)
		}
	}

	if len(dest) == 0 {
		return nil
	}

	r.hasDest.Store(n, struct{}{})
	_, _ = w.WriteString(`<a `)
	for _, attr := range r.Attributes {
		_, _ = w.WriteString(attr.Name)
		_, _ = w.WriteString(`="`)
		_, _ = w.Write(util.EscapeHTML([]byte(attr.Value)))
		_, _ = w.WriteString(`" `)
	}
	_, _ = w.WriteString(`href="`)
	_, _ = w.Write(util.URLEscape(dest, true /* resolve references */))
	_, _ = w.WriteString(`">`)
	return nil
}

func (r *Renderer) exit(w util.BufWriter, n *Node) {
	if _, ok := r.hasDest.LoadAndDelete(n); ok {
		_, _ = w.WriteString("</a>")
	}
	_, _ = w.WriteString("</span>")
}
