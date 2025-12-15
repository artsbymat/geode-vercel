package build

import (
	"html"
	"strings"

	"geode/internal/types"
	"geode/internal/utils"
)

func RenderExplorer(tree *types.FileTree) string {
	var b strings.Builder
	b.WriteString(`<ul class="file-explorer">`)
	for _, child := range tree.Children {
		renderNode(&b, child, "")
	}
	b.WriteString(`</ul>`)
	return b.String()
}

func normalizeExplorerLink(raw string) string {
	url := strings.TrimPrefix(raw, "/")
	url = strings.TrimSuffix(url, "/")

	if url == "" || url == "index" {
		return "/"
	}

	if strings.HasSuffix(url, "/index") {
		trimmed := strings.TrimSuffix(url, "/index")
		if trimmed == "" {
			return "/"
		}
		return "/" + trimmed
	}

	return "/" + url
}

func renderNode(b *strings.Builder, node *types.FileTree, parentKey string) {
	isFile := node.Link != "" || node.Path != ""

	key := node.Name
	if parentKey != "" {
		key = parentKey + "/" + node.Name
	}

	if isFile {
		b.WriteString("<li>")
	} else {
		b.WriteString(`<li data-node-key="` + html.EscapeString(key) + `">`)
	}

	if isFile {
		var link string
		if node.Link != "" {
			link = normalizeExplorerLink(node.Link)
		} else {
			url := utils.PathToSlug(node.Path)
			url = strings.TrimSuffix(url, ".md")
			link = normalizeExplorerLink(url)
		}

		title := node.Title
		if title == "" {
			title = node.Name
		}

		b.WriteString(`<a href="` + html.EscapeString(link) + `">` +
			html.EscapeString(title) + `</a>`)
	} else {
		b.WriteString(`<span class="folder">` +
			FolderChevronIcon +
			`<span class="folder-name">` +
			html.EscapeString(node.Name) +
			`</span></span>`)
	}

	if len(node.Children) > 0 {
		b.WriteString("<ul>")
		for _, c := range node.Children {
			renderNode(b, c, key)
		}
		b.WriteString("</ul>")
	}

	b.WriteString("</li>")
}
