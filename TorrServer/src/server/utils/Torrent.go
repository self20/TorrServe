package utils

import (
	"encoding/base32"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

var trackers = []string{
	"udp://tracker.zer0day.to:1337/announce",
	"udp://tracker.trackerfix.com:83/announce",
	"udp://tracker.thepiratebay.org:80/announce",
	"udp://tracker.publicbt.com:80/announce",
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://tracker.openbittorrent.com:80/announce",
	"udp://tracker.leechers-paradise.org:6969/announce",
	"udp://tracker.leechers-paradise.org:6969",
	"udp://tracker.ilibr.org:6969/announce",
	"udp://tracker.coppersurfer.tk:6969/announce",
	"udp://tracker.coppersurfer.tk:6969",
	"udp://shubt.net:2710",
	"udp://shadowshq.yi.org:6969/announce",
	"udp://public.popcorn-tracker.org:6969/announce",
	"udp://opentor.org:2710",
	"udp://open.demonii.com:1337",
	"udp://ipv6.leechers-paradise.org:6969",
	"udp://explodie.org:6969/announce",
	"udp://explodie.org:6969",
	"udp://eddie4.nl:6969/announce",
	"udp://denis.stalker.h3q.com:6969/announce",
	"udp://bt.rutor.org:2710",
	"udp://46.148.18.250:2710",
	"https://open.kickasstracker.com:443/announce",
	"http://www.torrentheaven.de/announce.php",
	"http://www.progressivetorrents.com/announce.php",
	"http://www.filebase.ws:5678/announce",
	"http://www.bittorrent-support.com/announce.php",
	"http://tracker3.torrentino.com/announce?passkey=00000000000000000000000000000000",
	"http://tracker2.wasabii.com.tw:6969/announce",
	"http://tracker2.itzmx.com:6961/announce",
	"http://tracker.torrentbay.to:6969/announce",
	"http://tracker.torrent.to:2710/announce",
	"http://tracker.tiny-vps.com:6969/announce",
	"http://tracker.thepiratebay.org/announce",
	"http://tracker.tfile.me/announce.php?uk=aFnt7k16j6",
	"http://tracker.tfile.co:80/announce",
	"http://tracker.tfile.co/announce.php?uk=aFnt7k16j6",
	"http://tracker.publicbt.com:80/announce",
	"http://tracker.openbittorrent.com/announce",
	"http://tracker.mp3-es.com/announce.php",
	"http://tracker.istole.it:80/announce",
	"http://tracker.internetwarriors.net:1337/announce",
	"http://tracker.grepler.com:6969/announce",
	"http://tracker.filetracker.pl:8089/announce",
	"http://tracker.electro-torrent.pl:80/announce",
	"http://tracker.dler.org:6969/announce",
	"http://tracker.city9x.com:2710/announce",
	"http://tracker.calculate.ru:6969/announce",
	"http://torrentsmd.eu:8080/announce",
	"http://sound-park.ru/announce.php",
	"http://share.camoe.cn:8080/announce",
	"http://retracker.spark-rostov.ru:80/announce",
	"http://retracker.mgts.by:80/announce",
	"http://retracker.local/announce",
	"http://retracker.home/announce",
	"http://retracker.bashtel.ru:80/announce",
	"http://retracker.bashtel.ru/announce.php",
	"http://radioarchive.cc/announce.php",
	"http://piratbit.net/bt/announce.php",
	"http://papaja.v2v.cc:6970/announce",
	"http://medbit.ru/announce.php",
	"http://i.bandito.org/announce.php?uk=aFnt7k16j6",
	"http://exodus.desync.com/announce",
	"http://denis.stalker.h3q.com:6969/announce",
	"http://bt.nnm-club.ru:2710/announce",
	"http://bt.nnm-club.info:2710/announce",
	"http://87.248.186.252:8080/announce",
	"http://80.246.243.18:6969/announce",
	"http://46.4.109.148:6969/announce",
	"http://210.244.71.25:6969/announce",
	"http://182.176.139.129:6969/announce",
	"http://125.227.35.196:6969/announce",
}

func AddRetracker(magnet string) string {
	tr := strings.Join(trackers, "&tr=")
	return magnet + "&tr=" + tr
}

func RemoveRetracker(magnet string) string {
	m, err := metainfo.ParseMagnetURI(magnet)
	if err != nil {
		fmt.Println("Error remove retracker:", err)
		return magnet
	}
	m.Trackers = []string{}
	return m.String()
}

func PeerIDRandom(peer string) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return peer + base32.StdEncoding.EncodeToString(randomBytes)[:20-len(peer)]
}

func GotInfo(t *torrent.Torrent) error {
	select {
	case <-t.GotInfo():
		return nil
	case <-time.Tick(time.Second * 120):
		return errors.New("timeout load torrent info")
	}
}
