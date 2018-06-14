package mods

/*
	e.POST("/torrent/add", torrentAdd)
	e.POST("/torrent/upload", torrentUpload)
	e.POST("/torrent/get", torrentGet)
	e.POST("/torrent/rem", torrentRem)
	e.POST("/torrent/list", torrentList)
	e.POST("/torrent/stat", torrentStat)

	e.POST("/torrent/cleancache", torrentCleanCache)
	e.GET("/torrent/restart", torrentRestart)

	e.GET("/torrent/playlist/:hash/*", torrentPlayList)
	e.GET("/torrent/playlist.m3u", torrentPlayListAll)

	e.GET("/torrent/view/:hash/:file", torrentView)
	e.HEAD("/torrent/view/:hash/:file", torrentView)
	e.GET("/torrent/preload/:hash/:file", torrentPreload)
*/

var corejs = `

function addTorrent(link, save){
	var magnet = $("#magnet").val();
	$("#magnet").val("");
	if(magnet!=""){
		var magJS = JSON.stringify({ Link: magnet });
		$.post('/torrent/add',magJS)
		.done(function( data ) {
			loadTorrents();
		})
		.fail(function( data ) {
			alert(data.responseJSON.message);
		});
	}
}

function getTorrent(hash){
}

function removeTorrent(hash){
}

function statTorrent(hash){
}

function listTorrent(){
}


function cleanCache(hash){
}

function restartService(){
}

function preloadTorrent(hash, fileLink){
}

`
