package build

import (
	"fmt"
	"geode/internal/config"
	"html"
	"strings"
)

func RenderSocials(socials []config.Social) string {
	if len(socials) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString(`<ul class="social-links">`)

	for _, s := range socials {
		name := html.EscapeString(s.Title)
		url := html.EscapeString(s.Link)

		fmt.Fprintf(
			&b,
			`<li><a href="%s" target="_blank" rel="noopener">%s</a></li>`,
			url, name,
		)
	}

	b.WriteString(`</ul>`)

	return b.String()
}
