package service

import (
	"fmt"
	"shantaram-cms/pkg/database"
)

type General struct {
	db *database.Database
}

func NewGeneral(db *database.Database) *General {
	return &General{
		db: db,
	}
}

func (s *General) GetByID(id string) (any, error) {
	page, err := s.db.GetGeneralByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get general by id %s: %v", id, err)
	}

	return page, nil
}

func (s *General) Upsert(id string, data any) error {
	if id != "menu" && id != "background" {
		return fmt.Errorf("invalid id %s", id)
	}

	if err := s.db.UpsertGeneral(id, data); err != nil {
		return fmt.Errorf("failed to upsert general by id %s: %v", id, err)
	}

	return nil
}
