package templates

import (
	"server/version"
)

var mainPage = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link href="` + faviconB64 + `" rel="icon" type="image/x-icon">
	<link rel="stylesheet" href="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css">
	<script src="http://code.jquery.com/jquery-1.11.3.min.js"></script>
	<script src="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js"></script>
	<script src="/js/api.js"></script>
	<title>Torrent server</title>
</head>
<body>
<style>
.ui-btn {
   word-wrap: break-word !important;
   white-space: normal !important;
}
</style>

<div data-role="page">

	<div data-role="header"><h3>TorrServer ` + version.Version + `</h3></div>

	<div data-role="content">
		<h3>Add torrent: </h3>
		<input id="magnet" autocomplete="off">
		<div class="ui-grid-a">
			<div class="ui-block-a"><button id="buttonAdd" data-icon="plus" onclick="addTorr()">Add</button></div>
			<div class="ui-block-b"><button id="buttonUpload" data-icon="plus">Upload</button></div>
		</div>
		
		<br>
		<a href="/search" rel="external" target="_blank" data-role="button" data-icon="search">Search</a>
		<a href="/torrent/playlist.m3u" rel="external" data-role="button" data-icon="bullets">Playlist</a>
		<br>
		<h3>Torrents: </h3>
		<div id="torrents"></div>
		<hr>

		<div class="ui-grid-a">
			<div class="ui-block-a"><button id="buttonShutdown" data-icon="power" onclick="shutdown()">Shutdown</button></div>
			<div class="ui-block-b"><a href="/settings" rel="external" data-role="button" data-icon="gear" id="buttonSettings">Settings</a></div>
		</div>
	</div>

	<form id="uploadForm" style="display:none" action="/torrent/upload" method="post">
		<input type="file" id="filesUpload" style="display:none" multiple onchange="uploadTorrent()" name="files"/> 
	</form>
	
	<div data-role="footer">
		<center><p><a rel="external" style="text-decoration: none;" href="/about">About</a></p></center>
	</div>
</div> 

<script>
	function addTorr(){
		var magnet = $("#magnet").val();
		$("#magnet").val("");
		if(magnet!=""){
			addTorrent(magnet,true,
			function( data ) {
				loadTorrents();
			},
			function( data ) {
				alert(data.responseJSON.message);
			});
		}
	}
	
	function removeTorr(hash){
		if(hash!=""){
			removeTorrent(hash,
			function( data ) {
				loadTorrents();
			},
			function( data ) {
				alert(data.responseJSON.message);
			});
		}
	};
	
	function shutdown(){
		shutdownServer(function( data ) {
				alert(data.responseJSON.message);
		});
	}

	$( document ).ready(function() {
		watchInfo();
	});

	$('#buttonUpload').click(function() {
   		$('#filesUpload').click();
	});
	
	function uploadTorrent() {
		var form = $("#uploadForm");
		var formData = new FormData(document.getElementById("uploadForm"));
		var data = new FormData();
		$.each($('#filesUpload')[0].files, function(i, file) {
    		data.append('file-'+i, file);
		});
		$.ajax({
				cache: false,
				processData: false,
				contentType: false,
				type: form.attr('method'),
				url: form.attr('action'),
				data: data
				}).done(function(data) {
					loadTorrents();
				}).fail(function(data) {
					alert(data.responseJSON.message);
				});
	}
	
	$('#uploadForm').submit(function(event) {
		event.preventDefault();
		var form = $(this);
		$.ajax({
			type: form.attr('method'),
			url: form.attr('action'),
			data: form.serialize()
			}).done(function(data) {
				loadTorrents();
			});
	});

	function loadTorrents() {
		listTorrent(
			function( data ) {
				var torrents = $("#torrents");
				torrents.empty();
				var html = "";
				var queueInfo = [];
				for(var key in data) {
					var tor = data[key];
					if (tor.IsGettingInfo){
						queueInfo.push(tor);
						continue;
					}
					html += tor2Html(tor);
				}
				if (queueInfo.length>0){
					html += "<br><hr><h3>Got info: </h3>";
					for(var key in queueInfo) {
						var tor = queueInfo[key];
						html += tor2Html(tor);
					}
				}
				$(html).appendTo(torrents);
				torrents.enhanceWithin();
			},
			function( data ) {
				alert(data.responseJSON.message);
			});
	}
	
	function tor2Html(tor){
		var html = '<hr>';
		html += '<div id="'+tor.Hash+'" data-role="collapsible">';
		if (tor.IsGettingInfo)
			html += '<h3>'+tor.Name+' '+humanizeSize(tor.Size)+' '+tor.Hash+'</h3>';
		else
			html += '<h3>'+tor.Name+' '+humanizeSize(tor.Size)+'</h3>';
		html += '<button data-icon="delete" onclick="removeTorrent(\''+tor.Hash+'\');">Remove ['+tor.Name+']</button>';
		if (typeof tor.Files != 'undefined' && tor.Files != 0){
			html += '<br>';
			html += '<a data-role="button" data-icon="bullets" target="_blank" href="'+tor.Playlist+'">Playlist</a>';
			for(var i in tor.Files){
				var file = tor.Files[i];
				var ico = "";
				if (file.Viewed)
					ico = 'data-icon="check"';
				html += '<a '+ico+' data-role="button" target="_blank" onClick="loadTorrents()" href="'+file.Link+'">'+file.Name+" "+humanizeSize(file.Size)+'</a>';
			}
		}
		html += '</div>';
		return html;
	}
	
	function watchInfo(){
		var lastTorrentCount = 0;
		var lastGettingInfo = 0;
		timer = setInterval(function() {
			listTorrent(
			function( data ) {
				var gettingInfo = 0;
				for(var key in data) {
					var tor = data[key];
					if (tor.IsGettingInfo)
						gettingInfo++;
				}
	
				if (lastTorrentCount!=data.length || gettingInfo!=lastGettingInfo){
					loadTorrents();
					lastTorrentCount = data.length;
					lastGettingInfo = gettingInfo;
				}
			});
		}, 1000);
	}
</script>
</body>
</html>
`

func (t *Template) parseMainPage() {
	parsePage(t, "mainPage", mainPage)
}
