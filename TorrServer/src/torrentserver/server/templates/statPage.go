package templates

var statPage = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="http://code.jquery.com/jquery-1.11.3.min.js"></script>
	<title>Torrent status</title>
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
		<h3>Torrents: </h3>
		<div id="torrents"/>
	</div>

	<div data-role="footer"></div>
</div> 

<script>
	$( document ).ready(function() {
		setInterval(loadStat(), 1000);
	});

	function loadStat() {
		var torrents = $("#torrents");
		torrents.empty();
		$.post('/stat/get')
			.done(function( data ) {
				for(var key in data.Torrents) {
					var tor = data.Torrents[key];
					$("<p>"+tor.Name+"</p>").appendTo(torrents);
					$("<p>"+tor.Magnet+"</p>").appendTo(torrents);
					for(var i in tor.Files){
						var file = tor.Files[i];
						$("<p>"+file.Path+" "+file.PreloadOffset+" / "+" "+file.Size+" "+file.Priority+"</p>").appendTo(torrents);
					}
					$("<br>").appendTo(torrents);
					for(var i in tor.Pieces){
						var piece = tor.Pieces[i];
						$("<p>"+piece.Hash+" "+piece.Name+" "+piece.Priority+" "+piece.Completed+" "+piece.Size+"</p>").appendTo(torrents);
					}
					$("<br>").appendTo(torrents);
					$("<hr>").appendTo(torrents);
				}
				torrents.enhanceWithin();
			})
			.fail(function( data ) {
				$("<p>"+data.responseJSON.message+"</p>").appendTo(torrents);
				alert();
			});
	}
</script>
</body>
</html>
`

func (t *Template) parseStatPage() {
	parsePage(t, "statPage", statPage)
}
