package database

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
)

func (d *Database) UpsertGeneral(id string, data any) error {
	key := []byte(id)

	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(generalBucket)

		value, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal: %v", err)
		}

		return b.Put(key, value)
	})
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	return nil
}

func (d *Database) GetGeneralByID(id string) (any, error) {
	var data any

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(generalBucket)
		key := []byte(id)
		value := b.Get(key)

		if value == nil {
			return nil
		}

		err := json.Unmarshal(value, &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal general: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view failed: %v", err)
	}

	return data, nil
}

func (d *Database) GetGeneralByIDData(id string) ([]byte, error) {
	var data []byte

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(generalBucket)
		key := []byte(id)
		data = b.Get(key)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view failed: %v", err)
	}

	return data, nil
}
