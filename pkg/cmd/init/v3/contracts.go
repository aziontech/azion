package init

type Template struct {
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

type Item struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Preset  string `json:"preset"`
	Path    string `json:"path"`
	Link    string `json:"link"`
}
