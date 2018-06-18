package torr

import (
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
)

func (bt *BTServer) addQueue(tor *torrent.Torrent, onAdd func(*TorrentState)) {
	go func() {
		mi := tor.Metainfo()
		fmt.Println("Geting torrent info:", mi.Magnet(tor.Name(), tor.InfoHash()))
		st := NewState(tor)
		st.IsGettingInfo = true
		bt.Watching(st)

		select {
		case <-tor.GotInfo():
			//get all info
			count := 0
			for tor.Info() != nil && len(tor.Files()) == 0 && count < 60 {
				<-tor.GotInfo()
				time.Sleep(time.Millisecond * 200)
				count++
			}
			st.updateTorrentState()
			st.IsGettingInfo = false

			fmt.Println("Torrent received info:", tor.Name())
			go onAdd(st)
		case <-tor.Closed():
			bt.removeState(st.Hash)
			fmt.Println("Torrent closed:", tor.Name())
		}
	}()
}
