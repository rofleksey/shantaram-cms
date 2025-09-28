package mapper

import (
	"fmt"
	"shantaram/app/api"
	"shantaram/pkg/database"
	"shantaram/pkg/util"
	"strings"
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

func OrderToNotificationText(o database.Order) string {
	var builder strings.Builder

	builder.WriteString("Заказ #")
	builder.WriteString(fmt.Sprint(o.Index))
	builder.WriteString("\n\n")

	builder.WriteString("Имя: ")
	builder.WriteString(o.ClientName)
	builder.WriteString("\n")

	builder.WriteString("Телефон: ")
	builder.WriteString(o.ClientPhone)
	builder.WriteString("\n")

	if o.ClientComment != nil {
		builder.WriteString("Комментарий: ")
		builder.WriteString(*o.ClientComment)
		builder.WriteString("\n")
	}

	builder.WriteString("\nТовары: \n")

	var totalPrice float64

	for i, item := range o.Items {
		builder.WriteString(fmt.Sprint(i + 1))
		builder.WriteString(". ")
		builder.WriteString(item.Title)
		builder.WriteString(" x ")
		builder.WriteString(fmt.Sprint(item.Amount))
		builder.WriteString(" - ")
		builder.WriteString(fmt.Sprintf("%.2f", util.FixPrice(item.Price*float64(item.Amount))))
		builder.WriteString(" ₽\n")

		totalPrice += util.FixPrice(float64(item.Amount) * item.Price)
	}

	totalPrice = util.FixPrice(totalPrice)

	builder.WriteString("\n")
	builder.WriteString("Сумма: ")
	builder.WriteString(fmt.Sprintf("%.2f", totalPrice))
	builder.WriteString(" ₽\n\n")
	builder.WriteString("https://admin.shantaram-spb.ru/#/order/")
	builder.WriteString(o.ID.String())

	return builder.String()
}
