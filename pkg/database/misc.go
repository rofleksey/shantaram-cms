package database

type MenuProduct struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price int    `json:"price"`
}

type MenuGroup struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Products []MenuProduct `json:"products"`
}

type MenuSettings struct {
	Groups []MenuGroup `json:"groups"`
}

func (m *MenuSettings) FindProduct(id string) *MenuProduct {
	for _, group := range m.Groups {
		for _, product := range group.Products {
			if product.ID == id {
				return &product
			}
		}
	}

	return nil
}
