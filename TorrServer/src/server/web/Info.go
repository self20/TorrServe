package server

import (
	"fmt"
	"net/http"
	"sort"

	"server/utils"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/bytes"
)

func initInfo(e *echo.Echo) {
	server.GET("/stat", statePage)
}

func statePage(c echo.Context) error {
	state := bts.BTState()

	msg := ""

	msg += fmt.Sprintf("Listen port: %d<br>\n", state.LocalPort)
	msg += fmt.Sprintf("Peer ID: %+q<br>\n", state.PeerID)
	msg += fmt.Sprintf("Banned IPs: %d<br>\n", state.BannedIPs)

	for _, dht := range state.DHTs {
		msg += fmt.Sprintf("%s DHT server at %s:<br>\n", dht.Addr().Network(), dht.Addr().String())
		dhtStats := dht.Stats()
		msg += fmt.Sprintf("\t&emsp;# Nodes: %d (%d good, %d banned)<br>\n", dhtStats.Nodes, dhtStats.GoodNodes, dhtStats.BadNodes)
		msg += fmt.Sprintf("\t&emsp;Server ID: %x<br>\n", dht.ID())
		msg += fmt.Sprintf("\t&emsp;Announces: %d<br>\n", dhtStats.SuccessfulOutboundAnnouncePeerQueries)
		msg += fmt.Sprintf("\t&emsp;Outstanding transactions: %d<br>\n", dhtStats.OutstandingTransactions)
	}

	sort.Slice(state.Torrents, func(i, j int) bool {
		return state.Torrents[i].Hash < state.Torrents[j].Hash
	})
	msg += "Torrents:<br>\n"
	for _, st := range state.Torrents {
		msg += fmt.Sprintf("Name: %v<br>\n", st.Name)
		msg += fmt.Sprintf("Hash: %v<br>\n", st.Hash)

		msg += fmt.Sprintf("\t&emsp;TotalPeers:   	 %v<br>\n", st.TotalPeers)
		msg += fmt.Sprintf("\t&emsp;PendingPeers: 	 %v<br>\n", st.PendingPeers)
		msg += fmt.Sprintf("\t&emsp;ActivePeers:      %v<br>\n", st.ActivePeers)
		msg += fmt.Sprintf("\t&emsp;ConnectedSeeders: %v<br>\n", st.ConnectedSeeders)
		msg += fmt.Sprintf("\t&emsp;HalfOpenPeers: 	 %v<br>\n", st.HalfOpenPeers)
		msg += fmt.Sprintf("\t&emsp;BytesWritten:     %v<br>\n", bytes.Format(st.BytesWritten))
		msg += fmt.Sprintf("\t&emsp;BytesWrittenData: %v<br>\n", bytes.Format(st.BytesWrittenData))
		msg += fmt.Sprintf("\t&emsp;BytesRead: 			%v<br>\n", bytes.Format(st.BytesRead))
		msg += fmt.Sprintf("\t&emsp;BytesReadData: 		%v<br>\n", bytes.Format(st.BytesReadData))
		msg += fmt.Sprintf("\t&emsp;BytesReadUsefulData: %v<br>\n", bytes.Format(st.BytesReadUsefulData))
		msg += fmt.Sprintf("\t&emsp;ChunksWritten:      %v<br>\n", st.ChunksWritten)
		msg += fmt.Sprintf("\t&emsp;ChunksRead: 	       %v<br>\n", st.ChunksRead)
		msg += fmt.Sprintf("\t&emsp;ChunksReadUseful:   %v<br>\n", st.ChunksReadUseful)
		msg += fmt.Sprintf("\t&emsp;ChunksReadUnwanted: %v<br>\n", st.ChunksReadUnwanted)
		msg += fmt.Sprintf("\t&emsp;PiecesDirtiedGood: %v<br>\n", st.PiecesDirtiedGood)
		msg += fmt.Sprintf("\t&emsp;PiecesDirtiedBad:  %v<br>\n", st.PiecesDirtiedBad)

		msg += fmt.Sprintf("\t&emsp;Download Speed: %v/Sec<br>\n", utils.Format(st.DownloadSpeed))
		msg += fmt.Sprintf("\t&emsp;Upload Speed:   %v/Sec<br>\n", utils.Format(st.UploadSpeed))

		msg += fmt.Sprintf("Cache:<br>\n")
		msg += fmt.Sprintf("Capacity: %v<br>\n", bytes.Format(st.CacheState.Capacity))
		msg += fmt.Sprintf("Filled: %v<br>\n", bytes.Format(st.CacheState.Filled))
		msg += fmt.Sprintf("PiecesLength: %v<br>\n", bytes.Format(st.CacheState.PiecesLength))
		msg += fmt.Sprintf("PiecesCount: %v<br>\n", st.CacheState.PiecesCount)
		for _, p := range st.CacheState.Pieces {
			msg += fmt.Sprintf("\t&emsp;Piece: %v\t&emsp; Access: %s\t&emsp; Buffer size: %d(%s)\t&emsp; Complete: %v\t&emsp; Hash: %s\n<br>", p.Id, p.Accessed.Format("15:04:05.000"), p.BufferSize, bytes.Format(int64(p.BufferSize)), p.Completed, p.Hash)
		}

		msg += "<hr><br><br>\n\n"
	}
	return c.HTML(http.StatusOK, msg)
}
