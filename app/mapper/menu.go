package mapper

import (
	"shantaram/app/api"
	"shantaram/pkg/database"
)

func MapMenu(m database.Menu) api.Menu {
	return api.Menu{
		Groups: []api.ProductGroup{},
		Id:     m.ID,
		Title:  m.Title,
	}
}

func MapProductGroup(g database.ProductGroup) api.ProductGroup {
	return api.ProductGroup{
		Created:  g.Created,
		Id:       g.ID,
		Products: []api.Product{},
		Title:    g.Title,
		Updated:  g.Updated,
	}
}

func MapProduct(p database.Product) api.Product {
	return api.Product{
		Available:   p.Available,
		Created:     p.Created,
		Description: p.Description,
		Id:          p.ID,
		Index:       int(p.Index),
		Price:       p.Price,
		Title:       p.Title,
		Updated:     p.Updated,
	}
}
