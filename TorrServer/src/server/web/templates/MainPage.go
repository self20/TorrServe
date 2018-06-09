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
			<div class="ui-block-a"><button id="buttonAdd" data-icon="plus" onclick="addTorrent()">Add</button></div>
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
	function addTorrent(){
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
	
	function shutdown(){
		$.post('/shutdown');
	}

	function removeTorrent(id){
		if(id!=""){
			var magJS = JSON.stringify({ Hash: id });
			$.post('/torrent/rem', magJS)
			.done(function( data ) {
				loadTorrents();
			})
			.fail(function( data ) {
				alert(data.responseJSON.message);
			});
		}
	};

	$( document ).ready(function() {
		loadTorrents();
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
    event.preventDefault(); // Prevent the form from submitting via the browser
    var form = $(this);
    $.ajax({
      type: form.attr('method'),
      url: form.attr('action'),
      data: form.serialize()
    }).done(function(data) {
      loadTorrents();
    }).fail(function(data) {
      // Optionally alert the user of an error here...
    });
  });

	function loadTorrents() {
		$.post('/torrent/list')
			.done(function( data ) {
				var torrents = $("#torrents");
				torrents.empty();
				for(var key in data) {
					var tor = data[key];
					$("<hr>").appendTo(torrents);
					var divColl = $('<div id="'+tor.Hash+'" data-role="collapsible"></div>')
					$("<h3>"+tor.Name+" "+humanizeSize(tor.Size)+"</h3>").appendTo(divColl);
					$('<a data-role="button" data-icon="bullets" target="_blank" href="'+tor.Playlist+'">Playlist</a>').appendTo(divColl);
					$('<button data-icon="delete" onclick="removeTorrent(\''+tor.Hash+'\');">Remove ['+tor.Name+']</button>').appendTo(divColl);
					$("<br>").appendTo(divColl);
					for(var i in tor.Files){
						var file = tor.Files[i];
						var btn = $('<a data-role="button" target="_blank" onClick="loadTorrents()" href="'+file.Link+'">'+file.Name+" "+humanizeSize(file.Size)+'</a>');
						if (file.Viewed)
							btn.buttonMarkup({ icon: "check" });
						btn.appendTo(divColl);
					}
					divColl.appendTo(torrents);
				}
				torrents.enhanceWithin();
			})
			.fail(function( data ) {
				alert(data.responseJSON.message);
			});
	}

	function humanizeSize(size) {
		var i = Math.floor( Math.log(size) / Math.log(1024) );
    	return ( size / Math.pow(1024, i) ).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
	};

</script>
</body>
</html>
`

func (t *Template) parseMainPage() {
	parsePage(t, "mainPage", mainPage)
}
