package database

import (
	"strconv"
	"strings"
	"time"
)

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

type OrderStatus string

const (
	OrderStatusOpen      OrderStatus = "open"
	OrderStatusClosed    OrderStatus = "closed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
	Amount int    `json:"amount"`
}

type Order struct {
	ID           uint64      `json:"id"`
	TableID      string      `json:"tableID"`
	Created      time.Time   `json:"created"`
	Status       OrderStatus `json:"status"`
	Name         string      `json:"name"`
	Phone        string      `json:"phone"`
	Comment      string      `json:"comment"`
	AdminComment string      `json:"adminComment"`
	Seen         bool        `json:"seen"`
	Items        []OrderItem `json:"items"`
}

func (o *Order) TelegramString() string {
	var builder strings.Builder

	builder.WriteString("Заказ #")
	builder.WriteString(strconv.FormatUint(o.ID, 10))
	builder.WriteString("\n\n")

	builder.WriteString("Имя: ")
	builder.WriteString(o.Name)
	builder.WriteString("\n")

	builder.WriteString("Телефон: ")
	builder.WriteString(o.Phone)
	builder.WriteString("\n")

	builder.WriteString("Комментарий: ")
	builder.WriteString(o.Comment)
	builder.WriteString("\n\n")

	builder.WriteString("Товары: \n")

	var totalPrice int

	for i, item := range o.Items {
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString(". ")
		builder.WriteString(item.Title)
		builder.WriteString(" x ")
		builder.WriteString(strconv.Itoa(item.Amount))
		builder.WriteString(" - ")
		builder.WriteString(strconv.Itoa(item.Price * item.Amount))
		builder.WriteString(" ₽\n")

		totalPrice += item.Amount * item.Price
	}

	builder.WriteString("\n")
	builder.WriteString("Сумма: ")
	builder.WriteString(strconv.Itoa(totalPrice))
	builder.WriteString(" ₽\n\n")
	builder.WriteString("https://shantaram-spb.ru/admin/order/")
	builder.WriteString(strconv.FormatUint(o.ID, 10))

	return builder.String()
}

type DataPage[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"totalCount"`
}
