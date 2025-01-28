package service

import (
	"fmt"
	"shantaram-cms/pkg/database"
)

type Page struct {
	db *database.Database
}

func NewPage(db *database.Database) *Page {
	return &Page{
		db: db,
	}
}

func (s *Page) GetByID(id string) (*database.Page, error) {
	page, err := s.db.GetPageByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get page by id %s: %v", id, err)
	}

	return page, nil
}

func (s *Page) GetAll() ([]database.Page, error) {
	pages, err := s.db.GetAllPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get all pages: %v", err)
	}

	return pages, nil
}

func (s *Page) Delete(id string) error {
	if err := s.db.DeletePage(id); err != nil {
		return fmt.Errorf("failed to delete page by id %s: %v", id, err)
	}

	return nil
}

func (s *Page) Insert(page *database.Page) error {
	if err := s.db.InsertPage(page); err != nil {
		return fmt.Errorf("failed to insert page by id %s: %v", page.ID, err)
	}

	return nil
}

func (s *Page) Update(page *database.Page) error {
	if err := s.db.UpdatePage(page); err != nil {
		return fmt.Errorf("failed to update page by id %s: %v", page.ID, err)
	}

	return nil
}
