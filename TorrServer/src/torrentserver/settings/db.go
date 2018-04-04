package settings

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var (
	db             *bolt.DB
	dbRootName     = []byte("Root")
	dbTorrentsName = []byte("Torrents")
	dbTorrentsTime = []byte("Times")
	dbSettingsName = []byte("Settings")
	dbFileViewName = []byte("FileView")
)

func openDB() error {
	if db != nil {
		return nil
	}

	var err error
	db, err = bolt.Open(filepath.Join(sets.SettingPath, "torrents.db"), 0666, nil)
	if err != nil {
		fmt.Print(err)
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists(dbRootName)
		if err != nil {
			return fmt.Errorf("could not create Root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists(dbSettingsName)
		if err != nil {
			return fmt.Errorf("could not create Settings bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists(dbTorrentsName)
		if err != nil {
			return fmt.Errorf("could not create Torrents bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		CloseDB()
	}
	return err
}

func CloseDB() {
	if db != nil {
		db.Close()
		db = nil
	}
}

func ReadSettingsDB() error {
	err := openDB()
	if err != nil {
		return err
	}

	err = db.View(func(tx *bolt.Tx) error {
		jsbuf := tx.Bucket(dbRootName).Bucket(dbSettingsName).Get([]byte("json"))
		err = json.Unmarshal(jsbuf, sets)
		if sets.PreloadBufferSize > sets.CacheSize {
			sets.PreloadBufferSize = sets.CacheSize
		}
		return err
	})
	return err
}

func SaveSettingsDB() error {
	err := openDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		setsDB, err := tx.Bucket(dbRootName).CreateBucketIfNotExists(dbSettingsName)
		if err != nil {
			return err
		}

		buf, err := json.Marshal(sets)
		if err != nil {
			return err
		}
		return setsDB.Put([]byte("json"), buf)
	})
}

func ReadTorrentsDB() ([]string, error) {
	err := openDB()
	if err != nil {
		return nil, err
	}

	torrs := make([]string, 0)
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(dbRootName).Bucket(dbTorrentsName).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			torrs = append(torrs, string(k))
		}
		return nil
	})
	return torrs, err
}

func SaveTorrentsDB(torrs []string) error {
	err := openDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(dbRootName).Bucket(dbTorrentsName) != nil {
			err := tx.Bucket(dbRootName).DeleteBucket(dbTorrentsName)
			if err != nil {
				return fmt.Errorf("could not clean Torrents bucket: %v", err)
			}
		}
		torrDB, err := tx.Bucket(dbRootName).CreateBucketIfNotExists(dbTorrentsName)
		if err != nil {
			return fmt.Errorf("could not create Torrents bucket after cleaning: %v", err)
		}
		for _, torr := range torrs {
			fmt.Println("Save torrent:", torr)
			err := torrDB.Put([]byte(torr), nil)
			if err != nil {
				return fmt.Errorf("could not insert torrent to db: %v %v", err, torr)
			}
		}
		return nil
	})
}

func SaveTorrentTime(hash, time string) error {
	err := openDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		timeDB, err := tx.Bucket(dbRootName).CreateBucketIfNotExists(dbTorrentsTime)
		if err != nil {
			return err
		}
		return timeDB.Put([]byte(hash), []byte(time))
	})
}

func RemTorrentTime(hash string) error {
	err := openDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		timeDB, err := tx.Bucket(dbRootName).CreateBucketIfNotExists(dbTorrentsTime)
		if err != nil {
			return err
		}
		return timeDB.Delete([]byte(hash))
	})
}

func GetTorrentTime(hash string) (string, error) {
	err := openDB()
	if err != nil {
		return "", err
	}
	time := ""
	err = db.View(func(tx *bolt.Tx) error {
		timeDB := tx.Bucket(dbRootName).Bucket(dbTorrentsTime)
		if timeDB == nil {
			return nil
		}
		v := timeDB.Get([]byte(hash))
		if v != nil {
			time = string(v)
		}
		return nil
	})
	return time, err
}

func SaveTorrentView(torrHash, fileHash string) error {
	err := openDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		viewDB, err := tx.Bucket(dbRootName).CreateBucketIfNotExists(dbFileViewName)
		if err != nil {
			return err
		}
		return viewDB.Put([]byte(torrHash+"/"+fileHash), nil)
	})
}

func ExistTorrView(torrHash, fileHash string) (exists bool, err error) {
	err = openDB()
	if err != nil {
		return false, err
	}
	exists = false
	err = db.View(func(tx *bolt.Tx) error {
		view := tx.Bucket(dbRootName).Bucket(dbFileViewName)
		if view == nil {
			return nil
		}
		k, _ := view.Cursor().Seek([]byte(torrHash + "/" + fileHash))
		if string(k) == torrHash+"/"+fileHash {
			exists = true
		}
		return nil
	})
	return
}
