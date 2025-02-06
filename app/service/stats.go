package service

import (
	"fmt"
	"math"
	"shantaram-cms/app/dao"
	"shantaram-cms/pkg/database"
	"sort"
	"time"
)

type Stats struct {
	db *database.Database
}

func NewStats(db *database.Database) *Stats {
	return &Stats{
		db: db,
	}
}

func (s *Stats) GetStats() (dao.Stats, error) {
	opdMap := make(map[int64]float64)
	sumMap := make(map[int64]float64)
	now := time.Now()

	if err := s.db.IterateOrders(func(order database.Order) error {
		if now.Sub(order.Created).Hours() > 24*31 {
			return nil
		}

		if order.Status != "closed" {
			return nil
		}

		year, month, day := order.Created.Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, time.Local).UnixMilli()

		opdMap[date] = math.Round(opdMap[date] + 1)

		for _, item := range order.Items {
			sumMap[date] = math.Round(sumMap[date] + float64(item.Price*item.Amount))
		}

		return nil
	}); err != nil {
		return dao.Stats{}, fmt.Errorf("failed to iterate orders: %w", err)
	}

	opdSlice := make([]dao.StatsData, 0, len(opdMap))
	for date, value := range opdMap {
		opdSlice = append(opdSlice, dao.StatsData{
			Date:  date,
			Value: value,
		})
	}
	sort.Slice(opdSlice, func(i, j int) bool {
		return opdSlice[i].Date < opdSlice[j].Date
	})

	sumSlice := make([]dao.StatsData, 0, len(sumMap))
	for date, value := range sumMap {
		sumSlice = append(sumSlice, dao.StatsData{
			Date:  date,
			Value: value,
		})
	}
	sort.Slice(sumSlice, func(i, j int) bool {
		return sumSlice[i].Date < sumSlice[j].Date
	})

	return dao.Stats{
		OrderPerDay: opdSlice,
		SumPerDay:   sumSlice,
	}, nil
}
