package database

import "time"

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
	Items        []OrderItem `json:"items"`
}

type DataPage[T any] struct {
	Data       []T
	TotalCount int
}
