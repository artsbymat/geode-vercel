package build

import (
	"encoding/json"
	"fmt"
	"geode/internal/types"
	"geode/internal/utils"
	"html"
)

func BuildGraph(page types.MetaMarkdown) *types.GraphData {
	nodes := make([]types.GraphNode, 0)
	links := make([]types.GraphLink, 0)
	nodeSet := make(map[string]bool)

	ensureNode := func(id, title string) {
		if !nodeSet[id] {
			nodes = append(nodes, types.GraphNode{
				ID:    id,
				Title: title,
				URL:   id,
			})
			nodeSet[id] = true
		}
	}

	currentPageURL := page.Link
	if currentPageURL == "" {
		currentPageURL = "/" + utils.PathToSlug(page.RelativePath)
	}

	ensureNode(currentPageURL, page.Title)

	for _, out := range page.OutgoingLinks {
		ensureNode(out.URL, out.Title)
		links = append(links, types.GraphLink{
			Source: currentPageURL,
			Target: out.URL,
		})
	}

	for _, back := range page.Backlinks {
		ensureNode(back.URL, back.Title)
		links = append(links, types.GraphLink{
			Source: back.URL,
			Target: currentPageURL,
		})
	}

	return &types.GraphData{
		Nodes: nodes,
		Links: links,
	}
}

func RenderGraphView(graphData *types.GraphData, currentPageURL string) string {
	if graphData == nil || len(graphData.Nodes) == 0 {
		return ""
	}

	filteredNodes := make([]types.GraphNode, 0)
	filteredLinks := make([]types.GraphLink, 0)
	nodeSet := make(map[string]bool)

	nodeSet[currentPageURL] = true

	for _, link := range graphData.Links {
		if link.Source == currentPageURL {
			nodeSet[link.Target] = true
		} else if link.Target == currentPageURL {
			nodeSet[link.Source] = true
		}
	}

	for _, node := range graphData.Nodes {
		if nodeSet[node.ID] {
			filteredNodes = append(filteredNodes, node)
		}
	}

	for _, link := range graphData.Links {
		if nodeSet[link.Source] && nodeSet[link.Target] {
			filteredLinks = append(filteredLinks, link)
		}
	}

	filteredGraphData := &types.GraphData{
		Nodes: filteredNodes,
		Links: filteredLinks,
	}

	jsonData, err := json.Marshal(filteredGraphData)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(`<div id="graph-container" data-graph=%q data-current-page=%q></div>`,
		html.EscapeString(string(jsonData)),
		html.EscapeString(currentPageURL))
}
