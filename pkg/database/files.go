package database

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
)

func (d *Database) GetFileByID(id uint64) (*File, error) {
	var file File

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(filesBucket)
		key := []byte(fmt.Sprintf("%d", id))
		v := b.Get(key)

		if v == nil {
			return fmt.Errorf("файл с id = %d не найден", id)
		}

		err := json.Unmarshal(v, &file)
		if err != nil {
			return fmt.Errorf("failed to unmarshal file: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view failed: %v", err)
	}

	return &file, nil
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
