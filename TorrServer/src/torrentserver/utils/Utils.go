package utils

import (
	"encoding/base32"
	"math/rand"
	"regexp"
	"runtime/debug"
	"strings"
)

func FileToLink(file string) string {
	re := regexp.MustCompile(`[ !\*'\(\);:@&=\+\$,/\?#\[\]~",]`)
	return re.ReplaceAllString(file, `_`)
}

func PeerIDRandom(peer string) string {
	return peer + getToken(20-len(peer))
}

func getToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func ReleaseMemory() {
	debug.FreeOSMemory()
}

var treckers = []string{
	"udp://opentor.org:2710",
	"udp://bt.rutor.org:2710",
	"http://tracker.thepiratebay.org/announce",
	"udp://tracker.thepiratebay.org:80/announce",
	"udp://denis.stalker.h3q.com:6969/announce",
	"http://denis.stalker.h3q.com:6969/announce",
	"udp://open.demonii.com:1337",
	"udp://tracker.coppersurfer.tk:6969",
	"udp://tracker.leechers-paradise.org:6969",
	"http://tracker.openbittorrent.com/announce",
	"udp://tracker.openbittorrent.com:80/announce",
	"http://tracker.torrentbay.to:6969/announce",
	"http://tracker.istole.it:80/announce",
	"http://tracker.torrent.to:2710/announce",
	"http://papaja.v2v.cc:6970/announce",
	"http://i.bandito.org/announce.php?uk=aFnt7k16j6&",
	"udp://tracker.publicbt.com:80/announce",
	"http://tracker.publicbt.com:80/announce",
	"http://tracker.tfile.me/announce.php?uk=aFnt7k16j6&",
	"http://tracker.tfile.co/announce.php?uk=aFnt7k16j6&",
	"http://retracker.home/announce",
	"http://tracker3.torrentino.com/announce?passkey=00000000000000000000000000000000",
	"http://bt.nnm-club.ru:2710/announce",
	"http://bt.nnm-club.info:2710/announce",
	"http://www.filebase.ws:5678/announce",
	"http://exodus.desync.com/announce",
	"http://www.progressivetorrents.com/announce.php",
	"http://retracker.bashtel.ru/announce.php",
	"http://radioarchive.cc/announce.php",
	"http://medbit.ru/announce.php",
	"http://piratbit.net/bt/announce.php",
	"http://sound-park.ru/announce.php",
	"udp://tracker.trackerfix.com:83/announce",
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://tracker.zer0day.to:1337/announce",
	"udp://explodie.org:6969/announce",
	"udp://eddie4.nl:6969/announce",
	"udp://tracker.ilibr.org:6969/announce",
	"udp://shadowshq.yi.org:6969/announce",
	"udp://shubt.net:2710",
	"http://retracker.local/announce",
	"http://tracker.filetracker.pl:8089/announce",
	"http://tracker2.wasabii.com.tw:6969/announce",
	"http://tracker.grepler.com:6969/announce",
	"http://80.246.243.18:6969/announce",
	"http://125.227.35.196:6969/announce",
	"http://tracker.tiny-vps.com:6969/announce",
	"http://87.248.186.252:8080/announce",
	"http://www.torrentheaven.de/announce.php",
	"http://tracker.mp3-es.com/announce.php",
	"http://tracker.calculate.ru:6969/announce",
	"http://210.244.71.25:6969/announce",
	"http://46.4.109.148:6969/announce",
	"http://tracker.dler.org:6969/announce",
}

func AddRetracker(magnet string) string {
	tr := strings.Join(treckers, "&tr=")
	return magnet + "&tr=" + tr
}
