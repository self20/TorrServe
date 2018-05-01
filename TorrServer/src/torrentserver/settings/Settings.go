package settings

import (
	"encoding/json"
	"time"

	"torrentserver/db"
)

type Settings struct {
	CacheSize         int // in byte, def 200 mb
	PreloadBufferSize int // buffer for readahead

	IsElementumCache bool

	//BT Config
	DisableTCP        bool
	DisableUTP        bool
	DisableUPNP       bool
	DisableDHT        bool
	DisableUpload     bool
	Encryption        int // 0 - Enable, 1 - disable, 2 - force
	DownloadRateLimit int // in kb, 0 - inf
	UploadRateLimit   int // in kb, 0 - inf
	ConnectionsLimit  int

	SettingPath string `json:"-"`
}

var (
	sets      *Settings
	StartTime time.Time
)

func init() {
	sets = new(Settings)
	sets.CacheSize = 200 * 1024 * 1024
	sets.PreloadBufferSize = sets.CacheSize / 2
	sets.ConnectionsLimit = 150
	StartTime = time.Now()
}

func Get() *Settings {
	return sets
}

func LoadFile() error {
	buf, err := db.ReadSettingsDB()
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, sets)
	return err
}

func SaveFile() error {
	buf, err := json.Marshal(sets)
	if err != nil {
		return err
	}
	return db.SaveSettingsDB(buf)
}
