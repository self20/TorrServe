package utils

import (
	"path/filepath"
	"sort"

	"github.com/anacrolix/torrent"
)

var extVideo = map[string]interface{}{
	".3g2":   nil,
	".3gp":   nil,
	".aaf":   nil,
	".asf":   nil,
	".avchd": nil,
	".avi":   nil,
	".drc":   nil,
	".flv":   nil,
	".m2ts":  nil,
	".ts":    nil,
	".m2v":   nil,
	".m4p":   nil,
	".m4v":   nil,
	".mkv":   nil,
	".mng":   nil,
	".mov":   nil,
	".mp2":   nil,
	".mp4":   nil,
	".mpe":   nil,
	".mpeg":  nil,
	".mpg":   nil,
	".mpv":   nil,
	".mxf":   nil,
	".nsv":   nil,
	".ogg":   nil,
	".ogv":   nil,
	".qt":    nil,
	".rm":    nil,
	".rmvb":  nil,
	".roq":   nil,
	".svi":   nil,
	".vob":   nil,
	".webm":  nil,
	".wmv":   nil,
	".yuv":   nil,
}

var extAudio = map[string]interface{}{
	".aac":  nil,
	".aiff": nil,
	".ape":  nil,
	".au":   nil,
	".flac": nil,
	".gsm":  nil,
	".it":   nil,
	".m3u":  nil,
	".m4a":  nil,
	".mid":  nil,
	".mod":  nil,
	".mp3":  nil,
	".mpa":  nil,
	".pls":  nil,
	".ra":   nil,
	".s3m":  nil,
	".sid":  nil,
	".wav":  nil,
	".wma":  nil,
	".xm":   nil,
}

func GetMimeType(filename string) string {
	ext := filepath.Ext(filename)
	if _, ok := extVideo[ext]; ok {
		return "video/*"
	}
	if _, ok := extAudio[ext]; ok {
		return "audio/*"
	}
	return "*/*"
}

func GetPlayableFiles(torr *torrent.Torrent) []*torrent.File {
	if torr != nil && torr.Info() != nil && len(torr.Files()) > 0 {
		files := make([]*torrent.File, 0)
		for _, f := range torr.Files() {
			if GetMimeType(f.Path()) != "*/*" {
				files = append(files, f)
			}
		}
		if len(files) > 1 {
			sort.Slice(files, func(i, j int) bool {
				return files[i].Path() < files[j].Path()
			})
		}
		return files
	}
	return nil
}
