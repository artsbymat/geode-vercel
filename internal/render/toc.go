package render

import (
	"geode/internal/types"
	"geode/internal/utils"
	"strings"
)

func TocFromURL(pages []types.MetaMarkdown, currentURL string) []types.TocItem {
	currentURL = strings.TrimSpace(currentURL)
	if currentURL == "" {
		return nil
	}

	for _, p := range pages {
		url := strings.TrimSpace(p.Link)
		if url == "" {
			url = "/" + utils.PathToSlug(p.RelativePath)
		}
		if url == currentURL {
			return p.TableOfContents
		}
	}

	return nil
}
