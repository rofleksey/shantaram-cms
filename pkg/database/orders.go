package database

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"time"
)

func (d *Database) IterateOrders(callback func(Order) error) error {
	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(ordersBucket)

		c := bucket.Cursor()
		for key, value := c.Last(); key != nil; key, value = c.Prev() {
			var order Order

			if err := json.Unmarshal(value, &order); err != nil {
				return fmt.Errorf("failed to unmarshal order: %v", err)
			}

			if err := callback(order); err != nil {
				return fmt.Errorf("callback error: %v", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("view failed: %v", err)
	}

	return nil
}

func (d *Database) GetOrdersPaginated(offset, limit int) (*DataPage[Order], error) {
	var result DataPage[Order]

	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(ordersBucket)
		indexBucket := tx.Bucket(ordersIndexBucket)

		c := indexBucket.Cursor()
		curPos := 0
		for indexKey, key := c.Last(); indexKey != nil; indexKey, key = c.Prev() {
			if curPos < offset {
				curPos++
				continue
			}

			if curPos >= offset+limit {
				break
			}

			v := bucket.Get(key)

			var order Order

			err := json.Unmarshal(v, &order)
			if err != nil {
				return fmt.Errorf("failed to unmarshal order: %v", err)
			}

			result.Data = append(result.Data, order)

			curPos++
		}

		result.TotalCount = bucket.Stats().KeyN

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view failed: %v", err)
	}

	return &result, nil
}

func (d *Database) GetOrderByID(id uint64) (*Order, error) {
	var order Order

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(ordersBucket)
		key := []byte(fmt.Sprintf("%d", id))
		v := b.Get(key)

		if v == nil {
			return fmt.Errorf("заказ с id = %d не найден", id)
		}

		err := json.Unmarshal(v, &order)
		if err != nil {
			return fmt.Errorf("failed to unmarshal order: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view failed: %v", err)
	}

	return &order, nil
}

func (d *Database) InsertOrder(order *Order) error {
	indexKey := []byte(order.Created.Format(time.RFC3339))

	err := d.db.Batch(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(ordersBucket)
		indexBucket := tx.Bucket(ordersIndexBucket)

		newId, err := bucket.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to get next id: %v", err)
		}

		order.ID = newId
		key := []byte(fmt.Sprintf("%d", newId))

		value, err := json.Marshal(order)
		if err != nil {
			return fmt.Errorf("failed to marshal: %v", err)
		}

		if bucket.Get(key) != nil {
			return fmt.Errorf("заказ с id = %d уже существует", order.ID)
		}

		if err := bucket.Put(key, value); err != nil {
			return fmt.Errorf("failed to insert order: %v", err)
		}

		if err := indexBucket.Put(indexKey, key); err != nil {
			return fmt.Errorf("failed to insert order index: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	return nil
}

func (d *Database) UpdateOrder(order *Order) error {
	key := []byte(fmt.Sprintf("%d", order.ID))

	value, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	err = d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(ordersBucket)

		if b.Get(key) == nil {
			return fmt.Errorf("заказ с id = %s не найден", order.ID)
		}

		return b.Put(key, value)
	})
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	return nil
}

func (d *Database) DeleteOrder(id uint64) error {
	key := []byte(fmt.Sprintf("%d", id))

	err := d.db.Batch(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(ordersBucket)
		indexBucket := tx.Bucket(ordersIndexBucket)

		v := bucket.Get(key)
		if v == nil {
			return fmt.Errorf("заказ с ID = %d не найден", id)
		}

		var order Order

		err := json.Unmarshal(v, &order)
		if err != nil {
			return fmt.Errorf("failed to unmarshal order: %v", err)
		}

		indexKey := []byte(order.Created.Format(time.RFC3339))

		if err := bucket.Delete(key); err != nil {
			return fmt.Errorf("failed to delete order: %v", err)
		}

		if err := indexBucket.Delete(indexKey); err != nil {
			return fmt.Errorf("failed to delete order index: %v", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	return nil
}
