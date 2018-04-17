package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func ReadSettingsDB() ([]byte, error) {
	err := openDB()
	if err != nil {
		return nil, err
	}
	sets := make([]byte, 0)
	err = db.View(func(tx *bolt.Tx) error {
		sdb := tx.Bucket(dbSettingsName)
		if sdb == nil {
			return fmt.Errorf("error load settings")
		}

		sets = sdb.Get([]byte("json"))
		if sets == nil {
			return fmt.Errorf("error load settings")
		}
		return nil
	})
	return sets, err
}

func SaveSettingsDB(sets []byte) error {
	err := openDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		setsDB, err := tx.CreateBucketIfNotExists(dbSettingsName)
		if err != nil {
			return err
		}
		return setsDB.Put([]byte("json"), []byte(sets))
	})
}
