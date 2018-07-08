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
				<div class="col-sm-4">
					<div class="input-group">
						<div class="input-group-prepend">
							<div class="input-group-text">Type</div>
						</div>
						<select id="sType">
							<option value="byName">Search by name</option>
							<option value="byWatching">Now watching</option>
							<option value="byFilter">Search by filter</option>
							<option value="byTorrents">Search torrents</option>
						</select>
					</div>
				</div>
				<div class="col-sm-4">
					<div class="form-check">
						<input id="SearchMovie" checked="checked" class="form-check-input" type="checkbox" autocomplete="off">
						<label for="SearchMovie" class="active">Search movie</label>
					</div>
				</div>
				<div class="col-sm-4">
					<div class="form-check">
						<input id="SearchTV" checked="checked" class="form-check-input" type="checkbox" autocomplete="off">
						<label for="SearchTV" class="active">Search TV</label>
					</div>
				</div>
			</div>
		</div>
        <br>

        <div id="sbName">
            <div class="input-group">
                <div class="input-group-prepend">
                    <div class="input-group-text">Name</div>
                </div>
                <input type="text" name="search_movie" id="sName" value="" class="form-control">
            </div>
        </div>

        <div id="sbFilter">
            <button class="btn btn-primary w-100" type="button" data-toggle="collapse" data-target="#filter">
                Filter
            </button>
            <div class="collapse" id="filter">
                <div class="btn-group btn-group-toggle" data-toggle="buttons">
                    <label class="btn btn-secondary active">
                        <input type="radio" name="options" id="asc" autocomplete="off" checked>Ascend &lt;
                    </label>
                    <label class="btn btn-secondary">
                        <input type="radio" name="options" id="desc" autocomplete="off">Descend &gt;
                    </label>
                </div>
                <div class="container" id="genres"></div>
            </div>
        </div>

        <br>

        <button id="search" class="btn btn-primary w-100" type="button">Search</button>
        <br>
		<br>
        <div id="movies" class="thumbnail-mousey"></div>
		<div id="torrents" class="content"></div>
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
					<div id="seasonsButtons" class="btn-group flex-wrap"></div>
					<div id="infoTorrents"></div>
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-danger" data-dismiss="modal">Close</button>
				</div>
			</div>
		</div>
	</div>
			
    <script>
        $(document).ready(function() {
            updateGenres();
            updateUI();
        });

        $("#sName").keyup(function(event) {
            if (event.keyCode === 13)
                $("#search").click();
        });

        $("#search").click(function() {
            page = 1;
            search();
        });

        var selectSearchType = $('#sType')[0].selectedIndex;
        $('#sType').on('change', function() {
            updateUI();
        });

        var page = 1;
        var pages = 1;

        function goPage(val) {
            if (page == val)
                return;
            page = val;
            searchMovies();
        }

        function updateUI() {
            selectSearchType = $('#sType')[0].selectedIndex;
            if (selectSearchType == 0 || selectSearchType == 3) {
                $('#sbName').show(200);
                $('#sbFilter').hide(200);
            } else if (selectSearchType == 1) {
                $('#sbName').hide(200);
                $('#sbFilter').hide(200);
            } else if (selectSearchType == 2) {
                $('#sbName').hide(200);
                $('#sbFilter').show(200);
            }
            updatePages();
        }

        function updatePages() {
            if (pages == 1) {
                $('#pagesBlock').hide(0);
                return;
            } else
                $('#pagesBlock').show(0);
            $('#pages').empty();
            var html = "";
            for (i = 1; i <= pages; i++) {
                if (i == page)
                    html += '<li class="page-item active"><button class="page-link">' + i + '</button></li>';
                else
                    html += '<li class="page-item"><button class="page-link" onclick="goPage(' + i + ')">' + i + '</button></li>';
            }
            $(html).appendTo("#pages");
        }

        function updateGenres() {
            $.post('/search/genres')
                .done(function(data) {
                    $('#genres').empty();
                    var html = '<div class="row"><div class="col">With genre</div><div class="col">Without genre</div></div>';
                    for (var key in data) {
                        var gen = data[key];
                        html += '<div class="row"><div class="col">';
                        html += '<input id="wg' + gen.ID + '" class="form-check-input with_genre" type="checkbox" autocomplete="off">';
                        html += '<label for="wg' + gen.ID + '">' + gen.Name + '</label>';
                        html += '</div><div class="col">';
                        html += '<input id="wog' + gen.ID + '" class="form-check-input without_genre" type="checkbox" autocomplete="off">';
                        html += '<label for="wog' + gen.ID + '">' + gen.Name + '</label>';
                        html += '</div></div>';
                    }
                    $(html).appendTo('#genres');
                }).fail(function(data) {
                    alert(data.responseJSON.message);
                });
        }
			
		function search(){
			if (selectSearchType<3)
				searchMovies();
			else
				searchTorrents();
		}

        function searchMovies() {
            $('#search').prop("disabled", true);
            $('#pagesBlock').prop("disabled", true);

            if (typeof page != "number")
                page = 1;

            var SRequest = {
                "Name": $('#sName').val(),
                "Type": selectSearchType,
                "Page": page,
				"SearchMovie": $("#SearchMovie").prop('checked'),
				"SearchTV": $("#SearchTV").prop('checked')
            };
            if (selectSearchType == 2) {
                SRequest.Filter = getFilter();
            }

            $.post('/search/request', JSON.stringify(SRequest))
                .done(function(data) {
                    var movies = $("#movies");
					$("#torrents").empty();
                    movies.empty();
                    pages = data.Pages;
                    for (var key in data.Movies) {
                        var tor = data.Movies[key];
						var movtype = "Фильм";
                        if (tor.IsTv)
                            movtype = "Сериал";
						var name = tor.Title;
						if (!name)
							name = tor.OrigTitle;
						var year = tor.Date.substring(0,4);
						var overview = tor.Overview.replace(/"/g, '&quot;');
						overview = overview.replace(/'/g, '&apos;');
						overview = overview.replace(/\n/g, "<br>");
						var genres = "";
						if (tor.Genres)
							genres = tor.Genres.join(", ");
                        var html = '';
						html+= '<div id="m'+tor.Id+'" onclick="showModal(\''+name+'\',\''+overview+'\',\''+year+'\','+tor.Seasons+',\'\', \''+tor.BackdropUrl+'\')">';
						html+= '	<div class="thumbnail shadow">';
						html+= '		<h3>';
						html+= 				name + ' ('+ year +')<br>';
					    html+= '			<small>'+ movtype +'<br>'+ genres +'</small>';
						html+= '		</h3>';
						html+= '		<img class="img-responsive" src="'+ tor.PosterUrl +'">';
						html+= '	</div>';
						html+= '</div>';
                        $(html).appendTo(movies);
                    }
                    $('#search').prop("disabled", false);
                    $('#pagesBlock').prop("disabled", false);
                    updateUI();
                })
                .fail(function(data) {
                    $('#search').prop("disabled", false);
                    $('#pagesBlock').prop("disabled", false);
                    updateUI();
                    alert(data.responseJSON.message);
                });
        }
			
		function searchTorrents(){
			$('#search').prop("disabled", true);
            $('#pagesBlock').prop("disabled", true);
			$.post('/search/torrents', $('#sName').val())
			.done(function(torrList) {
				$("#movies").empty();
				var html = '<div class="btn-group-vertical d-flex" role="group">';
				for (var key in torrList) {
					torr = torrList[key];
					var dl = '';
					if (torr.PeersDl >= 0) {
						dl = '| ▼ ' + torr.PeersDl;
						dl += '| ▲ ' + torr.PeersUl;
					}
					html += '<button type="button" class="btn btn-secondary wrap w-100" onclick="doTorrent(\'' + torr.Magnet + '\', this)"><i class="fas fa-plus"></i> ' + torr.Name + " " + torr.Size + dl +'</button>';
				}
				html += '</div>';
				$('#torrents').html(html);
				$('#search').prop("disabled", false);
            	$('#pagesBlock').prop("disabled", false);
			})
			.fail(function(data) {
				$('#torrents').text(data.responseJSON.message);
				$('#search').prop("disabled", false);
            	$('#pagesBlock').prop("disabled", false);
			});
		}
			
		function showModal(Name, Overview, Year, SeasonsCount, Season, Backdrop){
			$('#infoModal').modal('show');
			$('#infoName').text(Name+ ' ' +Year);
			var img = '<img src="'+Backdrop+'" class="rounded leftimg">';
			$('#infoOverview').html(img + Overview);
			var fndStr = Name;
			if (Year && !Season && !SeasonsCount)
				fndStr += ' '+Year;
			if (Season)
				fndStr += ' S'+Season;
			if (SeasonsCount>0){
				var html = '<button type="button" class="btn btn-primary" onclick="showModal(\''+Name+'\',\''+Overview+'\',\''+Year+'\','+SeasonsCount+',\'\', \''+Backdrop+'\')">All</button>';
				for (var i = 1; i < SeasonsCount; i++){
					var ses = ('0' + i).slice(-2)
					html += '<button type="button" class="btn btn-primary" onclick="showModal(\''+Name+'\',\''+Overview+'\',\''+Year+'\','+SeasonsCount+',\''+ses+'\', \''+Backdrop+'\')">S'+ses+'</button>';
				}
				$('#seasonsButtons').html(html);
			}else{
				$('#seasonsButtons').empty();
			}
			$.post('/search/torrents', fndStr)
                .done(function(torrList) {
					var html = '<div class="btn-group-vertical d-flex" role="group">';
					for (var key in torrList) {
						torr = torrList[key];
						var dl = '';
						if (torr.PeersDl >= 0) {
							dl = '| ▼ ' + torr.PeersDl;
							dl += '| ▲ ' + torr.PeersUl;
						}
						html += '<button type="button" class="btn btn-secondary wrap w-100" onclick="doTorrent(\'' + torr.Magnet + '\', this)"><i class="fas fa-plus"></i> ' + torr.Name + " " + torr.Size + dl +'</button>';
					}
					html += '</div>';
					$('#infoTorrents').html(html);
				})
				.fail(function(data) {
					$('#infoTorrents').text(data.responseJSON.message);
				});
		}

        function getFilter() {
            var asc = $("#asc").prop("checked");
            var withg = [];
            var withoutg = [];
            $('.with_genre').each(function() {
                switch ($(this).prop("type")) {
                    case "checkbox":
                        if ($(this).prop("checked"))
                            withg.push(+$(this).prop("id").replace('wg', ''));
                        break;
                }
            });
            $('.without_genre').each(function() {
                switch ($(this).prop("type")) {
                    case "checkbox":
                        if ($(this).prop("checked"))
                            withoutg.push(+$(this).prop("id").replace('wog', ''));
                        break;
                }
            });
            return {
                "SortAsc": asc,
                "SortBy": "popularity",
                "DateLte": "",
                "DateGte": "",
                "WithGenres": withg,
                "WithoutGenres": withoutg
            };
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
