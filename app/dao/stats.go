package dao

type StatsData struct {
	Date  int64   `json:"date"`
	Value float64 `json:"value"`
}

type NamedData struct {
	Name  string  `json:"name"`
	Name1 string  `json:"name1"`
	Value float64 `json:"value"`
}

type Stats struct {
	OrderPerDay []StatsData `json:"ordersPerDay"`
	SumPerDay   []StatsData `json:"sumPerDay"`
}
