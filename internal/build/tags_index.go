package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"geode/internal/config"
	"geode/internal/types"
	"geode/internal/utils"
)

type TagIndexPage struct {
	Title string
	URL   string
	Tags  []TagLink
}

type TagLink struct {
	Name string
	URL  string
}

type TagIndexGroup struct {
	ID    string
	Tag   string
	Pages []TagIndexPage
}

type TagIndexData struct {
	Name       template.HTML
	Suffix     template.HTML
	Explorer   template.HTML
	Socials    template.HTML
	LiveReload bool

	TotalTags int
	TagGroups []TagIndexGroup
}

func BuildTagsIndex(cfg *config.Config, pages []types.MetaMarkdown, liveReload bool, fileTree *types.FileTree) error {
	templatePath := filepath.Join("themes", cfg.Theme, "templates", "tags.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parse tags template: %w", err)
	}

	byTag := make(map[string][]types.MetaMarkdown)
	for _, p := range pages {
		for _, raw := range p.Tags {
			t := strings.TrimSpace(raw)
			t = strings.TrimPrefix(t, "#")
			if t == "" {
				continue
			}
			byTag[t] = append(byTag[t], p)
		}
	}

	tags := make([]string, 0, len(byTag))
	for t := range byTag {
		tags = append(tags, t)
	}
	sort.Strings(tags)

	groups := make([]TagIndexGroup, 0, len(tags))
	for _, t := range tags {
		ps := byTag[t]
		sort.Slice(ps, func(i, j int) bool {
			return strings.ToLower(ps[i].Title) < strings.ToLower(ps[j].Title)
		})

		items := make([]TagIndexPage, 0, len(ps))
		for _, p := range ps {
			url := p.Link
			if url == "" {
				url = "/" + strings.TrimSuffix(utils.PathToSlug(p.RelativePath), ".md")
			}

			pageTags := make([]string, 0, len(p.Tags))
			for _, raw := range p.Tags {
				tt := strings.TrimSpace(raw)
				tt = strings.TrimPrefix(tt, "#")
				if tt == "" {
					continue
				}
				pageTags = append(pageTags, tt)
			}
			sort.Strings(pageTags)

			tagLinks := make([]TagLink, 0, len(pageTags))
			for _, tt := range pageTags {
				tagLinks = append(tagLinks, TagLink{
					Name: tt,
					URL:  "/tags/" + escapeTagPath(tt),
				})
			}

			items = append(items, TagIndexPage{
				Title: p.Title,
				URL:   url,
				Tags:  tagLinks,
			})
		}

		groups = append(groups, TagIndexGroup{ID: utils.PathToSlug(t), Tag: t, Pages: items})
	}

	data := TagIndexData{
		Name:       template.HTML(cfg.Site.Name),
		Suffix:     template.HTML(cfg.Site.Suffix),
		Explorer:   template.HTML(RenderExplorer(fileTree)),
		Socials:    template.HTML(""),
		LiveReload: liveReload,
		TotalTags:  len(tags),
		TagGroups:  groups,
	}

	outPath := filepath.Join("public", "tags.html")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
