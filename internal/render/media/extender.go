package media

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Extender struct{}

var _ goldmark.Extender = (*Extender)(nil)

func (e *Extender) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{}, 999),
		),
	)
}
