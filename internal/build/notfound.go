package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"geode/internal/config"
	"geode/internal/types"
)

type NotFoundData struct {
	Name       template.HTML
	Explorer   template.HTML
	Socials    template.HTML
	LiveReload bool
	HasTwitter bool
	HasMermaid bool
}

func Build404(cfg *config.Config, liveReload bool, fileTree *types.FileTree) error {
	templatePath := filepath.Join("themes", cfg.Theme, "templates", "404.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parse 404 template: %w", err)
	}

	data := NotFoundData{
		Name:       template.HTML(cfg.Site.Name),
		Explorer:   template.HTML(RenderExplorer(fileTree)),
		Socials:    template.HTML(RenderSocials(cfg.Socials)),
		LiveReload: liveReload,
		HasTwitter: false,
		HasMermaid: false,
	}

	outPath := filepath.Join("public", "404.html")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
