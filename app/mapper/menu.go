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
		Price:       p.Price,
		Title:       p.Title,
		Updated:     p.Updated,
	}
}

func MapOrder(o database.Order) api.Order {
	return api.Order{
		ClientComment: o.ClientComment,
		ClientName:    o.ClientName,
		ClientPhone:   o.ClientPhone,
		Created:       o.Created,
		Id:            o.ID,
		Index:         int(o.Index),
		Items:         o.Items,
		Seen:          o.Seen,
		Status:        o.Status,
		TableID:       o.TableID,
	}
}
