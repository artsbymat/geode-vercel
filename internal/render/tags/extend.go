package hashtag

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Extender struct {
	Resolver   Resolver
	Collector  *Collector
	Attributes []Attribute
}

var _ goldmark.Extender = (*Extender)(nil)

func (e *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&Parser{}, 999),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{
				Resolver:   e.Resolver,
				Collector:  e.Collector,
				Attributes: e.Attributes,
			}, 999),
		),
	)
}
