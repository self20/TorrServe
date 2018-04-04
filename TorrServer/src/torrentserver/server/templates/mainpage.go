package templates

import "torrentserver/utils"

var mainPage = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
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

	<div data-role="header"></div>

	<div data-role="content">
		<h3>Add torrent: </h3>
		<input id="magnet" autocomplete="off">
		<button id="buttonAdd" data-icon="plus" onclick="addTorrent()">Add</button>
		<br>

		<h3>Torrents: </h3>
		<div id="torrents"></div>
		<hr>
	</div>

	<div data-role="footer">
	<center><p>TorrServer ` + utils.Version + `</p></center>
	</div>
</div> 

<script>
	function addTorrent(){
		var magnet = $("#magnet").val();
		$("#magnet").val("");
		if(magnet!=""){
			var magJS = JSON.stringify({ Magnet: magnet });
			$.post('/torrent/add',magJS)
			.done(function( data ) {
				loadTorrents();
			})
			.fail(function( data ) {
				alert(data.responseJSON.message);
			});
		}
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

	function loadTorrents() {
		$.post('/torrent/list')
			.done(function( data ) {
				var torrents = $("#torrents");
				torrents.empty();
				for(var key in data) {
					var tor = data[key];
					$("<hr>").appendTo(torrents);
					var divColl = $('<div id="'+tor.Hash+'" data-role="collapsible"></div>')
					$("<h3>"+tor.Name+"</h3>").appendTo(divColl);
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
