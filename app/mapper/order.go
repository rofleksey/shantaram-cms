package mapper

import (
	"fmt"
	"shantaram/app/api"
	"shantaram/pkg/database"
	"strings"

	"github.com/rofleksey/meg"
)

func MapOrder(o database.Order) api.Order {
	return api.Order{
		ClientComment: o.ClientComment,
		ClientName:    o.ClientName,
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
		builder.WriteString(fmt.Sprintf("%.2f", meg.FixPrice(item.Price*float64(item.Amount))))
		builder.WriteString(" ₽\n")

		totalPrice += meg.FixPrice(float64(item.Amount) * item.Price)
	}

	totalPrice = meg.FixPrice(totalPrice)

	builder.WriteString("\n")
	builder.WriteString("Сумма: ")
	builder.WriteString(fmt.Sprintf("%.2f", totalPrice))
	builder.WriteString(" ₽\n\n")
	builder.WriteString("https://admin.shantaram-spb.ru/#/order/")
	builder.WriteString(o.ID.String())

	return builder.String()
}
