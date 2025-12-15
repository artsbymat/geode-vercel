package render

import (
	"geode/internal/types"
	"strings"
)

func BuildFileTree(pages []types.MetaMarkdown) *types.FileTree {
	root := &types.FileTree{Name: "root", Children: []*types.FileTree{}}

	for _, p := range pages {
		segments := strings.Split(p.RelativePath, "/")
		insertIntoTree(root, segments, p)
	}

	return root
}

func insertIntoTree(node *types.FileTree, segments []string, page types.MetaMarkdown) {
	if len(segments) == 0 {
		return
	}

	current := segments[0]
	var child *types.FileTree

	for _, c := range node.Children {
		if c.Name == current {
			child = c
			break
		}
	}

	if child == nil {
		child = &types.FileTree{
			Name:     current,
			Children: []*types.FileTree{},
		}
		node.Children = append(node.Children, child)
	}

	if len(segments) == 1 {
		child.Path = page.RelativePath
		child.Title = page.Title
		child.Link = page.Link
		return
	}

	insertIntoTree(child, segments[1:], page)
}
