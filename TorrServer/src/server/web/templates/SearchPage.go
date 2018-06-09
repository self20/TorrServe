package templates

import "server/version"

var searchPage = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link href="` + faviconB64 + `" rel="icon" type="image/x-icon">
	<link rel="stylesheet" href="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css">
	<script src="http://code.jquery.com/jquery-1.11.3.min.js"></script>
	<script src="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js"></script>
	<title>Search</title>
	
	<style>
	.movie{
		height: 150px;
		background: #ccc;
		background-repeat: no-repeat;
		background-size: cover;
		background-position: center center;
	}
	.poster{
		height: 100%;
		float: left;
		padding-right: 10px;
		padding-bottom: 10px;
	}
	.description{
		text-shadow: -1px 0 #333, 0 1px #333, 1px 0 #333, 0 -1px #333, #000 0 0 5px; 
		color: #56ffaa;
		margin-right: 10px;
		height: 100%;
	}
	.ui-btn {
   		word-wrap: break-word !important;
		white-space: normal !important;
	}
	.hidden { display: none; }
	</style>
</head>
<body>
<div data-role="page">
	<div data-role="header"><h3>TorrServer ` + version.Version + `</h3></div>
	<div data-role="content">
		<div class="ui-field-contain">
			<label for="sName">Name</label>
			<input type="text" name="search_movie" id="sName" value="">
		</div>
		<div class="ui-field-contain">
			<label>
				<input id="hideWOTorrents" type="checkbox">Hide without torrents
			</label>
		</div>
		<div class="ui-field-contain">
			<label for="sType">Type</label>
			<select id="sType">
				<option value="byName">Search by name</option>
				<option value="byWatching">Now watching</option>
				<option value="byFilter">Search by filter</option>
			</select>
		</div>
		<div class="ui-field-contain">	
			<div data-role="collapsible" id="filter">
				<h3>Filter</h3>
				<div class="ui-field-contain">
					<fieldset data-role="controlgroup" data-type="horizontal" data-mini="true">
						<input name="ascd" id="asc" checked="checked" type="radio">
						<label for="asc">Ascend &lt;</label>
						<input name="ascd" id="desc" type="radio">
						<label for="desc">Descend &gt;</label>
					</fieldset>
				</div>
				<div class="ui-field-contain">
					<div class="ui-grid-a">
						<div class="ui-block-a">
							<fieldset data-role="controlgroup" id="with_genre">
							</fieldset>
						</div>
						<div class="ui-block-b">
							<fieldset data-role="controlgroup" id="without_genre">
							</fieldset>
						</div>
					</div>
				</div>
			</div>
		</div>
		<button id="search" data-icon="search" onclick="searchTorrents()">Search</button>
		<div class="ui-field-contain">
			<div class="ui-grid-b">
				<div class="ui-block-a">
					<button id="PagePrev" data-icon="carat-l" onclick="prevPage()">Prev</button>
				</div>
				<div class="ui-block-b">
					<div style="text-align: center" class="ui-body ui-body-d" id="lPage">Page</div>
				</div>
				<div class="ui-block-c">
					<button id="PageNext" data-icon="carat-r" onclick="nextPage()">Next</button>
				</div>
			</div>
		</div>
		<br>
		<div id="torrents"></div>
	</div>
	<div data-role="footer">
		<p style="text-align:center;"><a rel="external" style="text-decoration: none;" href="/about">About</a></p>
	</div>
</div>
<script>
	$( document ).ready(function() {
		loadConfig();
		updateUI();
	});
	
	$("#sName").keyup(function(event) {
    	if (event.keyCode === 13)
        	$("#search").click();
	});
	
	var page = 1;
	var pages = 1;
	
	function nextPage(){
		if (page>=pages)
			return;
		page++;
		searchTorrents();
	}
	
	function prevPage(){
		if (page<2)
			return;
		page--;
		searchTorrents();
	}
	
	function updateUI(){
		selectSearchType = $('#sType')[0].selectedIndex;
		if (selectSearchType != 2)
			$('#filter').addClass('ui-disabled');
		else
			$('#filter').removeClass('ui-disabled');
	
		if (selectSearchType != 0)
			$('#sName').addClass('ui-disabled');
		else
			$('#sName').removeClass('ui-disabled');
		
		if (typeof page != "number")
			page = 1;
		$('#lPage').text('Page '+page+' / '+pages);
	}
	
	function loadConfig(){
		$.post('/search/genres')
			.done(function(data){
				$('#with_genre').empty();
				$('#without_genre').empty();
				var html = '<div class="ui-bar ui-bar-a">With Genre</div>';
				$(html).appendTo($('#with_genre'));
				html = '<div class="ui-bar ui-bar-a">Without Genre</div>';
				$(html).appendTo($('#without_genre'));
		
				for(var key in data) {
					var gen = data[key];
					var html = '<label for="wg'+gen.ID+'">'+gen.Name+'</label>';
					html += '<input data-mini="true" type="checkbox" id="wg'+gen.ID+'">';
					$(html).appendTo($('#with_genre'));
					html = '<label for="wog'+gen.ID+'">'+gen.Name+'</label>';
					html += '<input data-mini="true" type="checkbox" id="wog'+gen.ID+'">';
					$(html).appendTo($('#without_genre'));
				}
				$('#with_genre').enhanceWithin();
				$('#without_genre').enhanceWithin();
			}).fail(function( data ) {
				alert(data.responseJSON.message);
			});
	}
		
	var selectSearchType = 0;
		
	$('#sType').on('change', function () {
    	updateUI();
	});
		
	function searchTorrents() {
		$('#search').prop("disabled", true);
		$('#PagePrev').prop("disabled", true);
		$('#PageNext').prop("disabled", true);
		
		var hide = $('#hideWOTorrents').prop('checked');
		if (typeof page != "number")
			page = 1;
		
		var SRequest = {"Name":$('#sName').val(), "Type":selectSearchType, "Page":page, "HideWTorrent":hide};
		if (selectSearchType == 2){
			SRequest.Filter = getFilter(); 
		}
		
		$.post('/search/request', JSON.stringify(SRequest))
			.done(function( data ) {
				var torrents = $("#torrents");
				torrents.empty();
				$('<hr>').appendTo(torrents);
				pages = data.Pages;
				for(var key in data.Movies) {
					var tor = data.Movies[key];
					var html = '';
					html+='<div onclick="toggleInfo(\'#torr'+key+'\')" class="movie" style="background-image: url('+tor.BackdropUrl+');">';
					html+=' <img class="poster" src="'+tor.PosterUrl+'"/>';
					html+=' <div class="description">';
					html+='  <h4 style="padding-top: 10px;">'+tor.Title+' / '+tor.OrigTitle+'</h4>';
					html+='  <p style="float:left">'+tor.Date+'</p>';
					var movtype = "Фильм"; 
					if (tor.IsTv)
						movtype = "Сериал";
					html+='  <p style="float:right">'+movtype+'</p>';
					html+=' </div>';
					html+='</div>';
					html+='<div style="clear:both"></div>';
					html+= getTorrList(key, tor.Torrents, tor.Overview);
					html+='<hr>'; 
					$(html).appendTo(torrents);
				}
				torrents.enhanceWithin();
				$('#search').prop("disabled", false);
				$('#PagePrev').prop("disabled", false);
				$('#PageNext').prop("disabled", false);
				updateUI();
			})
			.fail(function( data ) {
				$('#search').prop("disabled", false);
				$('#PagePrev').prop("disabled", false);
				$('#PageNext').prop("disabled", false);
				updateUI();
				alert(data.responseJSON.message);
			});
	}
	
	function getFilter(){
		var asc = $("#asc").prop("checked");
		var withg = [];
		var withoutg = [];
		$('input', '#with_genre').each(function() {
			switch($(this).prop("type")) { 
					case "checkbox":
						if ($(this).prop("checked"))
							withg.push(+$(this).prop("id").replace('wg', ''));
						break;  
				}
		});
		$('input', '#without_genre').each(function() {
			switch($(this).prop("type")) { 
					case "checkbox":   
						if ($(this).prop("checked"))
							withoutg.push(+$(this).prop("id").replace('wog', ''));   
						break;  
				}
		});
		return {"SortAsc":asc,"SortBy":"popularity","DateLte":"","DateGte":"","WithGenres":withg,"WithoutGenres":withoutg};
	}
	
	function toggleInfo(key){
		$(key).toggle(50);
	}
	
	function getTorrList(key, torrList, torrOverview){
		var html = ''; 
		html+= '<div class="hidden" id="torr'+key+'">';
		html+= '<p>'+torrOverview+'</p>';
		for(var key in torrList) {
			torr = torrList[key];
			var dl = '';
			if (torr.PeersDl>=0){
				dl = '| ▼ '+torr.PeersDl;
				dl+= '| ▲ '+torr.PeersUl;
			}
			html+='<button data-icon="plus" onclick="doTorrent(\''+torr.Magnet+'\', this)">'+torr.Name+" "+torr.Size+dl+'</button>'
		}
		html+= '</div>'
		return html;
	}
		
	function doTorrent(magnet, elem){
		$(elem).prop("disabled", true);
		var magJS = JSON.stringify({ Link: magnet });
		$.post('/torrent/add',magJS)
		.done(function( data ) {
			$(elem).prop("disabled", false);
		})
		.fail(function( data ) {
			$(elem).prop("disabled", false);
			alert(data.responseJSON.message);
		});
	}
</script>
</body>
</html>
`

func (t *Template) parseSearchPage() {
	parsePage(t, "searchPage", searchPage)
}
