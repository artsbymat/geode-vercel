package hashtag

import "github.com/yuin/goldmark/ast"

var Kind = ast.NewNodeKind("Hashtag")

type Node struct {
	ast.BaseInline

	Tag []byte
}

func (*Node) Kind() ast.NodeKind { return Kind }

func (n *Node) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"Tag": string(n.Tag),
	}, nil)
}
