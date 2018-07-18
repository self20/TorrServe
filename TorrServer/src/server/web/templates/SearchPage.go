package templates

import "server/version"

var searchPage = `
<!DOCTYPE html>
<html lang="ru">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="` + faviconB64 + `" rel="icon" type="image/x-icon">
    <script src="/js/api.js"></script>
    <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.1.0/css/all.css" integrity="sha384-lKuwvrZot6UHsBSfcMvOkWwlCMgc0TaWr+30HWe3a4ltaBwTZhyTEggF5tJv8tbt" crossorigin="anonymous">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css" integrity="sha384-WskhaSGFgHYWDcbwN70/dfYBj47jz9qbsMId/iRN3ewGhXQFZCSftd1LZCfmhktB" crossorigin="anonymous">
    <script src="http://code.jquery.com/jquery-1.11.3.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js" integrity="sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/js/bootstrap.min.js" integrity="sha384-smHYKdLADwkXOn1EmN1qk/HfnUcbVRZyYmZ4qpPea6sjB/pTJ0euyQp0Mk8ck+5T" crossorigin="anonymous"></script>
    <title>TorrServer ` + version.Version + `</title>
</head>

<body>
    <style>
		#movies {
			display: grid; 
			grid-template-columns: repeat(auto-fit, minmax(186px, 1fr));
    		justify-items: center;
		}
    
		.thumbnail {
			width: 185px;
			margin-bottom: 3px;
			line-height: 1.42857143;
			background-color: #282828;
			border: 1px solid #4a4a4a;
			border-radius: 0;
			transition: border .2s ease-in-out
		}
		
		.thumbnail-mousey .thumbnail {
			position: relative;
			overflow: hidden
		}
		
		.thumbnail-mousey .thumbnail h3 {
			position: absolute;
			bottom: 0;
			font-family: noto sans, sans-serif;
			font-weight: 400;
			font-size: 16px;
			text-shadow: 2px 2px 4px #000;
			width: 100%;
			margin: 0;
			padding: 4px;
			background-color: rgba(0, 0, 0, .6)
		}
		
		.thumbnail-mousey .thumbnail h3 {
			text-shadow: -1px 0 #333, 0 1px #333, 1px 0 #333, 0 -1px #333, #000 0 0 5px;
			color: #fff;
		}
		
		.thumbnail-mousey .thumbnail h3 small {
			text-shadow: -1px 0 #333, 0 1px #333, 1px 0 #333, 0 -1px #333, #000 0 0 5px;
			color: #ddd;
		}
		
		.thumbnail-mousey .thumbnail>img {
			width: 185px;
    		height: 278px;
		}
    
        .wrap {
			white-space: normal;
			word-wrap: break-word;
			word-break: break-all;
		}
    	.content {
    		padding: 20px;
    	}
    	.modal-lg {
			max-width: 90% !important;
    		margin: 20px auto;
		}
    	.leftimg {
    		float:left;
    		margin: 7px 7px 7px 0;
    		max-width: 300px;
    		max-height: 170px;
   		}
    
		.sk-cube-grid {
		  	width: 40px;
		  	height: 40px;
		  	margin: 10px auto;
		}
		
		.sk-cube-grid .sk-cube {
		  width: 33%;
		  height: 33%;
		  background-color: #333;
		  float: left;
		  -webkit-animation: sk-cubeGridScaleDelay 1.3s infinite ease-in-out;
				  animation: sk-cubeGridScaleDelay 1.3s infinite ease-in-out; 
		}
		.sk-cube-grid .sk-cube1 {
		  -webkit-animation-delay: 0.2s;
				  animation-delay: 0.2s; }
		.sk-cube-grid .sk-cube2 {
		  -webkit-animation-delay: 0.3s;
				  animation-delay: 0.3s; }
		.sk-cube-grid .sk-cube3 {
		  -webkit-animation-delay: 0.4s;
				  animation-delay: 0.4s; }
		.sk-cube-grid .sk-cube4 {
		  -webkit-animation-delay: 0.1s;
				  animation-delay: 0.1s; }
		.sk-cube-grid .sk-cube5 {
		  -webkit-animation-delay: 0.2s;
				  animation-delay: 0.2s; }
		.sk-cube-grid .sk-cube6 {
		  -webkit-animation-delay: 0.3s;
				  animation-delay: 0.3s; }
		.sk-cube-grid .sk-cube7 {
		  -webkit-animation-delay: 0s;
				  animation-delay: 0s; }
		.sk-cube-grid .sk-cube8 {
		  -webkit-animation-delay: 0.1s;
				  animation-delay: 0.1s; }
		.sk-cube-grid .sk-cube9 {
		  -webkit-animation-delay: 0.2s;
				  animation-delay: 0.2s; }
		
		@-webkit-keyframes sk-cubeGridScaleDelay {
		  0%, 70%, 100% {
			-webkit-transform: scale3D(1, 1, 1);
					transform: scale3D(1, 1, 1);
		  } 35% {
			-webkit-transform: scale3D(0, 0, 1);
					transform: scale3D(0, 0, 1); 
		  }
		}
		
		@keyframes sk-cubeGridScaleDelay {
		  0%, 70%, 100% {
			-webkit-transform: scale3D(1, 1, 1);
					transform: scale3D(1, 1, 1);
		  } 35% {
			-webkit-transform: scale3D(0, 0, 1);
					transform: scale3D(0, 0, 1);
		  } 
		}
    </style>

    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    	<a class="btn navbar-btn pull-left" href="/"><i class="fas fa-arrow-left"></i></a>
        <span class="navbar-brand mx-auto">
			TorrServer ` + version.Version + `
			</span>
    </nav>
    <div class="content">
		<div class="container-fluid">
			<div class="row">
			{{if not .IsTorrent}}
				<div class="col-auto">
					<div class="btn-group btn-group-toggle" data-toggle="buttons">
						<label id="stFBN" class="btn btn-secondary" onclick="update_search(0)">
							<input type="radio" name="stype">Find by name
						</label>
						<label id="stDiscover" class="btn btn-secondary" onclick="update_search(1)">
							<input type="radio" name="stype">Discover
						</label>
                	</div>
				</div>
			{{end}}
				<div class="col-auto">
					<div class="btn-group">
						<a id="stMovies" class="btn btn-secondary" href="?vt=movie">Movies</a>
						<a id="stShows" class="btn btn-secondary" href="?vt=show">Shows</a>
						<a id="stTorrents" class="btn btn-secondary" href="?vt=torrent">Torrents</a>
                	</div>
				</div>
			</div>
		</div>
        <br>
		{{if .IsTorrent}}
		<div>
			<div class="input-group">
                <div class="input-group-prepend">
                    <div class="input-group-text">Filter</div>
                </div>
                <input type="text" name="search_filter" id="sFilter" placeholder="2017;S01|01x;LostFilm|Кубик в Кубе;720|1080|BDRemux" value="" class="form-control">
            </div>
		</div>
		<br>
		{{end}}
        <div id="sbName">
            <div class="input-group">
                <div class="input-group-prepend">
                    <div class="input-group-text">Name</div>
                </div>
                <input type="text" name="search_movie" id="sName" value="" class="form-control">
            </div>
        </div>
		<br>
		{{if not .IsTorrent}}
        <div id="sbFilter">
			<div class="container-fluid">
				<div class="row">
					<div class="col-auto">
						<div class="input-group">
							<div class="input-group-prepend">
								<div class="input-group-text">Year</div>
							</div>
							<select id="fYear">
								<option></option>
								{{range .Years}}<option>{{.}}</option>{{end}}
							</select>
						</div>
					</div>
					<div class="col-auto">
						<div class="input-group">
							<div class="input-group-prepend">
								<div class="input-group-text">Sort</div>
							</div>
							<select id="fSort">
								<option></option>
								{{range .Sorts}}<option>{{.}}</option>{{end}}
							</select>
						</div>
					</div>
					<div class="col w-100">
						<button class="btn btn-primary w-100" type="button" data-toggle="collapse" data-target="#fGenres">
							Genres
						</button>
						<div class="collapse" id="fGenres">
							{{range .Genres}}
								<label><input class="gcheckbox" type="checkbox" id="g{{.ID}}">{{.Name}}</label>
							{{end}}
						</div>
					</div>
				</div>
			</div>
        </div>
        <br>
		{{end}}

        <button id="search" class="btn btn-primary w-100" type="button">Search</button>
        <br>
		<br>
		{{if .IsTorrent}}
		<div id="loader" style="display:none" class="sk-cube-grid">
		  <div class="sk-cube sk-cube1"></div>
		  <div class="sk-cube sk-cube2"></div>
		  <div class="sk-cube sk-cube3"></div>
		  <div class="sk-cube sk-cube4"></div>
		  <div class="sk-cube sk-cube5"></div>
		  <div class="sk-cube sk-cube6"></div>
		  <div class="sk-cube sk-cube7"></div>
		  <div class="sk-cube sk-cube8"></div>
		  <div class="sk-cube sk-cube9"></div>
		</div>
		{{end}}
			
		{{if not .IsTorrent}}
        <div id="movies" class="thumbnail-mousey">
		{{range .Items}}
			<div id="m{{.ID}}" onclick="showModal('{{.OriginalName}}','{{.Name}}','{{.Overview}}','{{.Year}}','{{.Seasons}}','', '{{.Backdrop}}')">
				<div class="thumbnail shadow">
					<h3>
						{{.Name}} ({{.Year}})<br>
						<small>{{range $index, $gen := .Genres}}{{if $index}},{{end}} {{$gen.Name}}{{end}}</small>
					</h3>
					<img class="img-responsive" src="{{.Poster}}">
				</div>
			</div>
    	{{end}}	
		</div>
		{{end}}
		{{if .IsTorrent}}
		<div id="torrents" class="content">
			{{range .Items}}
			<div class="btn-group d-flex" role="group">
				<button type="button" class="btn btn-secondary wrap w-100" onclick="doTorrent('{{.OriginalName}}', this)"><i class="fas fa-plus"></i>{{.Name}} {{.Year}}{{if gt .Seasons -1}} | ▲ {{.Seasons}} | ▼ {{.Episodes}}{{end}}</button>
				<a type="button" class="btn btn-secondary" href="/torrent/play?link={{.OriginalName}}&m3u=true">...</a>
			</div>
			{{end}}
		</div>
		{{end}}
        <br>
        <div id="pagesBlock">
            <ul id="pages" class="pagination justify-content-center flex-wrap">
            </ul>
        </div>
    </div>
    <footer class="page-footer navbar-dark bg-dark">
        <span class="navbar-brand d-flex justify-content-center">
			<a rel="external" style="text-decoration: none;" href="/about">About</a>
			</span>
    </footer>
	{{if not .IsTorrent}}
	<div class="modal fade" id="infoModal" role="dialog">
		<div class="modal-dialog modal-lg">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="infoName"></h5>
					<button type="button" class="close" data-dismiss="modal" aria-label="Close">
						<span aria-hidden="true">&times;</span>
					</button>
				</div>
				<div class="modal-body">
					<small id="infoOverview"></small>
					<div style="clear:both"></div>
					<div id="loader" class="sk-cube-grid">
					  <div class="sk-cube sk-cube1"></div>
					  <div class="sk-cube sk-cube2"></div>
					  <div class="sk-cube sk-cube3"></div>
					  <div class="sk-cube sk-cube4"></div>
					  <div class="sk-cube sk-cube5"></div>
					  <div class="sk-cube sk-cube6"></div>
					  <div class="sk-cube sk-cube7"></div>
					  <div class="sk-cube sk-cube8"></div>
					  <div class="sk-cube sk-cube9"></div>
					</div>
					<br>
					<div id="seasonsButtons" class="btn-group flex-wrap"></div>
					<div id="infoTorrents"></div>
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-danger" data-dismiss="modal">Close</button>
				</div>
			</div>
		</div>
	</div>
	{{end}}
    <script>
		var currentPage = 1;
        $(document).ready(function() {
            $('#sbName').show(0);
			$('#sbFilter').hide(0);
			
			var params = new URLSearchParams(document.location.search.substring(1));
			
			if (params.get('type')=='discover')
				update_search(1);
			
			if (params.get('sort_by'))
				$('#fSort').val(params.get('sort_by'));
			
			if (params.get('primary_release_year'))
				$('#fYear').val(params.get('primary_release_year'));
			
			if (params.get('with_genres')){
				var genres = params.get('with_genres').split(',');
				genres.forEach(function(genre){
					$('#g'+genre).prop("checked", true);
				});
			}
			
			if (params.get('query'))
				$('#sName').val(params.get('query'));
			
			{{if not .IsTorrent}}
			if (params.get('page'))
				currentPage = params.get('page');
			updatePages();
			{{end}}
			{{if .IsTorrent}}
			if (params.get('ft')){
				var fts = params.getAll('ft');
				$('#sFilter').val(fts.join(";"));
			}
			{{end}}
			
			if (params.get('vt')=="movie" || !params.get('vt'))
				$('#stMovies').addClass("active");
			if (params.get('vt')=="show")
				$('#stShows').addClass("active");
			if (params.get('vt')=="torrent")
				$('#stTorrents').addClass("active");
			
			if (params.get('type')=="discover" || !params.get('type'))
				$('#stDiscover').click();
			if (params.get('type')=="search")
				$('#stFBN').click();
        });
			
		{{if not .IsTorrent}}
		var searchType = 0;
			
		function update_search(stype){
			searchType = stype;
			if (stype==1){
				$('#sbName').hide(200);
				$('#sbFilter').show(200);
			}else{
				$('#sbName').show(200);
				$('#sbFilter').hide(200);
			}
		}
		{{end}}
		
        $("#sName").keyup(function(event) {
            if (event.keyCode === 13)
                $("#search").click();
        });

        $("#search").click(function() {
            search();
        });
		
		function search(){
			var params = new URLSearchParams(document.location.search.substring(1));
			var qparam = [];
			var vt = params.get("vt");
			if (vt != null)
				qparam.push('vt='+vt);
			{{if not .IsTorrent}}
			var lang = params.get("language");
			if (lang != null) 
				qparam.push('language='+lang);
			
			if (searchType==1){
				qparam.push('type=discover');
				var year = $("#fYear option:selected").text();
				var sort = $("#fSort option:selected").text();
			
				var genres = [];
				$('.gcheckbox').each(function(i,obj) {
					if ($(obj).prop("checked")){
       					var gid = $(obj).attr('id').substring(1);
						genres.push(gid);
					}
  				});
			
				if (year)
					qparam.push('primary_release_year='+year);
				if (sort)
					qparam.push('sort_by='+sort);
				if (genres.length>0)
					qparam.push('with_genres='+genres.join(","));
				window.location.href = '/search?'+qparam.join('&');
			}else{
				qparam.push('type=search');
				var query = $('#sName').val();
				if (query){
					qparam.push('query='+query);
					window.location.href = '/search?'+qparam.join('&');
				}
			}
			{{end}}
			{{if .IsTorrent}}
				var query = $('#sName').val();
				if (query){
					$('#loader').show(0);
					var filter = $('#sFilter').val().split(";");
					if (filter.length){
						var ft = 'ft='+filter.join("&ft=");
						qparam.push(ft);
					}
					qparam.push('query='+query);
					window.location.href = '/search?'+qparam.join('&');
				}
			{{end}}
		}
		
		{{if not .IsTorrent}}
		function goPage(page){
			var params = new URLSearchParams(document.location.search.substring(1));
			if (params.get('page')!=page){
				params.set('page', page);
				window.location.href = '/search?'+params.toString();
			}
		}
			
        function updatePages() {
            if (pages == 1) {
                $('#pagesBlock').hide(0);
                return;
            } else
                $('#pagesBlock').show(0);
            $('#pages').empty();
            var html = "";
            for (i = 1; i <= {{.Pages}}; i++) {
                if (i == currentPage)
                    html += '<li class="page-item active"><button class="page-link">' + i + '</button></li>';
                else
                    html += '<li class="page-item"><button class="page-link" onclick="goPage(' + i + ')">' + i + '</button></li>';
            }
            $(html).appendTo("#pages");
        }
		{{end}}
			
		function showModal(OrigName, Name, Overview, Year, SeasonsCount, Season, Backdrop){
			$('#infoTorrents').empty();
			$('#infoName').text(Name+ ' ' +Year);
			var img = '<img src="'+Backdrop+'" class="rounded leftimg">';
			$('#infoOverview').html(img + Overview);
			$('#infoModal').modal('show');
			$('#loader').fadeIn(700);
			
			var filter = [];
			if (Year && !Season && !SeasonsCount)
				filter.push(Year);
			if (Season){
				var ses = Season.padStart(2,"0");
				filter.push('S'+ses+'|'+ses+'x');
			}
			if (SeasonsCount>0){
				var html = '<button type="button" class="btn btn-primary" onclick="showModal(\''+OrigName+'\',\''+Name+'\',\''+Overview+'\',\''+Year+'\','+SeasonsCount+',\'\', \''+Backdrop+'\')">All</button>';
				for (var i = 1; i < +SeasonsCount+1; i++){
					var ses = (""+i).padStart(2,"0");
					html += '<button type="button" class="btn btn-primary" onclick="showModal(\''+OrigName+'\',\''+Name+'\',\''+Overview+'\',\''+Year+'\','+SeasonsCount+',\''+ses+'\', \''+Backdrop+'\')">S'+ses+'</button>';
				}
				$('#seasonsButtons').html(html);
			}else{
				$('#seasonsButtons').empty();
			}
			
			var fn = function(torrList) {
				var html = '';
				for (var key in torrList) {
					torr = torrList[key];
					var dl = '';
					if (torr.PeersDl >= 0) {
						dl += '| ▲ ' + torr.PeersUl;
						dl += '| ▼ ' + torr.PeersDl;
					}
					html += '<div class="btn-group d-flex" role="group">'
					html += '<button type="button" class="btn btn-secondary wrap w-100" onclick="doTorrent(\'' + torr.Magnet + '\', this)"><i class="fas fa-plus"></i> ' + torr.Name + " | " + torr.Size + dl +'</button>';
					html += '<a type="button" class="btn btn-secondary" href="/torrent/play?link='+encodeURIComponent(torr.Magnet)+'&m3u=true">...</a>'
					html += '</div>';
				}
				$('#loader').fadeOut(700);
				$('#infoTorrents').html(html);
			};
			
			searchTorrent(OrigName,filter,function(torrList){
				if (!torrList)
					searchTorrent(Name,filter,fn);
				else
					fn(torrList);
			});
		}
			
		function searchTorrent(query, filter, done, fail){
			var ftstr = 'ft='+filter.join("&ft=");
			$.get('/search/torrent?query='+query+'&'+ftstr)
			.done(function(torrList){
				done(torrList);
			})
			.fail(function(data){
				if (fail)
					fail();
			})
		}

        function doTorrent(magnet, elem) {
            $(elem).prop("disabled", true);
            var magJS = JSON.stringify({
                Link: magnet
            });
            $.post('/torrent/add', magJS)
                .fail(function(data) {
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
