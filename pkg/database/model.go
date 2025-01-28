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
