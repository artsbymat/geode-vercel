package build

import (
	"geode/internal/types"
	"html/template"
	"strconv"
	"strings"
)

func RenderTOC(items []types.TocItem) string {
	if len(items) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<ul class="toc-list">`)
	for _, it := range items {
		if strings.TrimSpace(it.ID) == "" {
			continue
		}
		b.WriteString(`<li class="toc-level-`)
		b.WriteString(strconv.Itoa(it.Level))
		b.WriteString(`"><a href="#`)
		b.WriteString(template.HTMLEscapeString(it.ID))
		b.WriteString(`">`)
		b.WriteString(template.HTMLEscapeString(it.Text))
		b.WriteString(`</a></li>`)
	}
	b.WriteString(`</ul>`)

	return b.String()
}
