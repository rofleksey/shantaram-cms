package database

import (
	_ "embed"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
)

var generalBucket = []byte("general")
var pageBucket = []byte("pages")
var filesBucket = []byte("files")
var ordersBucket = []byte("orders")
var ordersIndexBucket = []byte("ordersIndex")

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

func (d *Database) Init() error {
	if err := d.db.Batch(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(generalBucket); err != nil {
			return fmt.Errorf("failed to create general bucket: %v", err)
		}

		if _, err := tx.CreateBucketIfNotExists(pageBucket); err != nil {
			return fmt.Errorf("failed to create pages bucket: %v", err)
		}

		if _, err := tx.CreateBucketIfNotExists(filesBucket); err != nil {
			return fmt.Errorf("failed to create files bucket: %v", err)
		}

		if _, err := tx.CreateBucketIfNotExists(ordersBucket); err != nil {
			return fmt.Errorf("failed to create orders bucket: %v", err)
		}

		if _, err := tx.CreateBucketIfNotExists(ordersIndexBucket); err != nil {
			return fmt.Errorf("failed to create orders index bucket: %v", err)
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
