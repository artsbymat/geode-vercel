package types

type FileTree struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Title    string      `json:"title,omitempty"`
	Link     string      `json:"permalink,omitempty"`
	Children []*FileTree `json:"children,omitempty"`
}
