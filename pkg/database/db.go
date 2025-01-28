package database

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
)

var pageBucket = []byte("pages")
var filesBucket = []byte("files")

type Database struct {
	db *bbolt.DB
}

func New() (*Database, error) {
	db, err := bbolt.Open(filepath.Join("data", "db.db"), 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	return &Database{
		db: db,
	}, nil
}

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
			return fmt.Errorf("page with id = %s not found", id)
		}

		return b.Delete(key)
	})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	return nil
}

func (d *Database) InsertFile(file *File) error {
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(filesBucket)

		newId, err := b.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to get next id: %v", err)
		}

		file.ID = newId

		value, err := json.Marshal(file)
		if err != nil {
			return fmt.Errorf("failed to marshal: %v", err)
		}

		key := []byte(fmt.Sprintf("%d", newId))

		if b.Get(key) != nil {
			return fmt.Errorf("файл с id = %d уже существует", file.ID)
		}

		return b.Put(key, value)
	})
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	return nil
}

func (d *Database) DeleteFile(id uint64) error {
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(filesBucket)

		key := []byte(fmt.Sprintf("%d", id))
		v := b.Get(key)
		if v == nil {
			return fmt.Errorf("файл с id = %d не найден", id)
		}

		return b.Delete(key)
	})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	return nil
}

func (d *Database) GetAllFiles() ([]File, error) {
	var result []File

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(filesBucket)

		err := b.ForEach(func(k, v []byte) error {
			var file File

			err := json.Unmarshal(v, &file)
			if err != nil {
				return fmt.Errorf("failed to unmarshal file: %v", err)
			}

			result = append(result, file)

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
		result = []File{}
	}

	return result, nil
}

func (d *Database) Init() error {
	if err := d.db.Batch(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(pageBucket); err != nil {
			return fmt.Errorf("failed to create pages bucket: %v", err)
		}

		if _, err := tx.CreateBucketIfNotExists(filesBucket); err != nil {
			return fmt.Errorf("failed to create files bucket: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to init buckets: %v", err)
	}

	if _, err := d.GetPageByID("index"); err != nil {
		_ = d.InsertPage(&Page{
			ID:    "index",
			Title: "Главная",
		})
	}

	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
