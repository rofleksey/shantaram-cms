package database

type Element struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type Page struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Elements []Element `json:"elements"`
}

type File struct {
	ID    uint64 `json:"id"`
	Path  string `json:"path"`
	Title string `json:"title"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
}
