package settings

import "time"

type Settings struct {
	CacheSize         int // in byte, def 200 mb
	PreloadBufferSize int // buffer for readahead

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
	sets.ConnectionsLimit = 100
	StartTime = time.Now()
}

func Get() *Settings {
	return sets
}

func LoadFile(path string) error {
	sets.SettingPath = path
	return ReadSettingsDB()
}

func SaveFile(path string) error {
	if path == "" {
		path = sets.SettingPath
	}
	return SaveSettingsDB()
}
