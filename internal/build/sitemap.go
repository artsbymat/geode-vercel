package build

import (
	"encoding/xml"
	"fmt"
	"geode/internal/config"
	"geode/internal/types"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

type UrlSet struct {
	XMLName xml.Name     `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	URLs    []SitemapURL `xml:"url"`
}

func BuildSitemap(cfg *config.Config, pages []types.MetaMarkdown) error {
	outputDir := cfg.Build.Output
	if outputDir == "" {
		outputDir = "public"
	}

	var urls []SitemapURL
	baseURL := cfg.Site.BaseURL

	for _, page := range pages {
		if page.Link == "" {
			continue
		}

		lastMod := time.Now().Format("2006-01-02")

		// modified -> created -> now
		if v, ok := page.Frontmatter["modified"]; ok {
			switch t := v.(type) {
			case time.Time:
				lastMod = t.Format("2006-01-02")
			case string:
				lastMod = t
			}
		} else if v, ok := page.Frontmatter["created"]; ok {
			switch t := v.(type) {
			case time.Time:
				lastMod = t.Format("2006-01-02")
			case string:
				lastMod = t
			}
		}

		url := page.Link

		if url == "index" {
			url = ""
		}

		url = strings.TrimSuffix(url, "/index")

		urls = append(urls, SitemapURL{
			Loc:     baseURL + url,
			LastMod: lastMod,
		})
	}

	urlSet := UrlSet{URLs: urls}

	f, err := os.Create(filepath.Join(outputDir, "sitemap.xml"))
	if err != nil {
		return fmt.Errorf("create sitemap.xml: %w", err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("write xml header: %w", err)
	}

	encoder := xml.NewEncoder(f)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlSet); err != nil {
		return fmt.Errorf("encode sitemap: %w", err)
	}

	return nil
}
