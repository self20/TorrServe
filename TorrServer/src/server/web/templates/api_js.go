package templates

import (
	"net/http"

	"server/settings"
	"server/web/helpers"

	"github.com/labstack/echo"
)

var apijs = `
function addTorrent(link, save, done, fail){
	var reqJson = JSON.stringify({ Link: link, DontSave: !save});
	$.post('/torrent/add',reqJson)
	.done(function( data ) {
		done(data);
	})
	.fail(function( data ) {
		fail(data);
	});
}

function getTorrent(hash, done, fail){
	var reqJson = JSON.stringify({ Hash: hash});
	$.post('/torrent/get',reqJson)
	.done(function( data ) {
		done(data);
	})
	.fail(function( data ) {
		fail(data);
	});
}

function removeTorrent(hash, done, fail){
	var reqJson = JSON.stringify({ Hash: hash});
	$.post('/torrent/rem',reqJson)
	.done(function( data ) {
		done(data);
	})
	.fail(function( data ) {
		fail(data);
	});
}

function statTorrent(hash, done, fail){
	var reqJson = JSON.stringify({ Hash: hash});
	$.post('/torrent/rem',reqJson)
	.done(function( data ) {
		done(data);
	})
	.fail(function( data ) {
		fail(data);
	});
}

function listTorrent(done, fail){
	$.post('/torrent/list')
	.done(function( data ) {
		done(data);
	})
	.fail(function( data ) {
		fail(data);
	});
}

function restartService(done, fail){
	$.get('/torrent/restart')
	.done(function( data ) {
		done();
	})
	.fail(function( data ) {
		fail(data);
	});
}

function preloadTorrent(hash, fileLink, done, fail){
	$.get('/torrent/preload/'+hash+'/'+fileLink)
	.done(function( data ) {
		done();
	})
	.fail(function( data ) {
		fail(data);
	});
}

function shutdownServer(fail){
	$.post('/shutdown')
	.fail(function( data ) {
		fail(data);
	});
}

function humanizeSize(size) {
	if (typeof size == 'undefined' || size == 0)
		return "";
	var i = Math.floor( Math.log(size) / Math.log(1024) );
	return ( size / Math.pow(1024, i) ).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
}
`

func Api_JS(c echo.Context) error {
	http.ServeContent(c.Response(), c.Request(), "api.js", settings.StartTime, helpers.NewSeekingBuffer(apijs))
	return c.NoContent(http.StatusOK)
}
