package database

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
)

func (d *Database) GetAllPages() ([]Page, error) {
	var result []Page

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(pageBucket)

		err := b.ForEach(func(k, v []byte) error {
			var page Page

			err := json.Unmarshal(v, &page)
			if err != nil {
				return fmt.Errorf("failed to unmarshal page: %v", err)
			}

			if len(page.Elements) == 0 {
				page.Elements = []Element{}
			}

			result = append(result, page)

			return nil
		})
		if err != nil {
			return fmt.Errorf("bucket for each error: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view failed: %v", err)
	}

	if len(result) == 0 {
		result = []Page{}
	}

	return result, nil
}

func (d *Database) GetPageByID(id string) (*Page, error) {
	var page Page

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(pageBucket)
		key := []byte(id)
		v := b.Get(key)

		if v == nil {
			return fmt.Errorf("страница с id = %s не найдена", id)
		}

		err := json.Unmarshal(v, &page)
		if err != nil {
			return fmt.Errorf("failed to unmarshal page: %v", err)
		}

		if len(page.Elements) == 0 {
			page.Elements = []Element{}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view failed: %v", err)
	}

	return &page, nil
}

func (d *Database) InsertPage(page *Page) error {
	key := []byte(page.ID)

	value, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	err = d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(pageBucket)

		if b.Get(key) != nil {
			return fmt.Errorf("страница с id = %s уже существует", page.ID)
		}

		return b.Put(key, value)
	})
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	return nil
}

func (d *Database) UpdatePage(page *Page) error {
	key := []byte(page.ID)

	value, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	err = d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(pageBucket)

		if b.Get(key) == nil {
			return fmt.Errorf("страница с id = %s не найдена", page.ID)
		}

		return b.Put(key, value)
	})
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	return nil
}

func (d *Database) DeletePage(id string) error {
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(pageBucket)

		key := []byte(id)
		v := b.Get(key)
		if v == nil {
			return fmt.Errorf("страница с ID = %s не найдена", id)
		}

		return b.Delete(key)
	})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	return nil
}
