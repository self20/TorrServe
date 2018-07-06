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
        .movie {
            height: 150px;
            background: #ccc;
        }
    
    	.torrList {
    		text-shadow: -1px 0 #333, 0 1px #333, 1px 0 #333, 0 -1px #333, #000 0 0 5px;
            color: #56ffaa;
       		background: #ccc;
            background-repeat: no-repeat;
            background-size: cover;
            background-position: center center;
    	}
    	
    	.button-text-shadow {
	    	text-shadow: -1px 0 #333, 0 1px #333, 1px 0 #333, 0 -1px #333, #000 0 0 5px;
            color: #56ffaa;
    	}
        
        .poster {
            height: 100%;
            float: left;
            padding-right: 10px;
        }
        
        .description {
            text-shadow: -1px 0 #333, 0 1px #333, 1px 0 #333, 0 -1px #333, #000 0 0 5px;
            color: #56ffaa;
            margin-right: 10px;
            height: 100%;
        }
        
        .ui-btn {
            word-wrap: break-word !important;
            white-space: normal !important;
        }
        
        .hidden {
            display: none;
        }
        
        .content {
            margin: 1%;
        }
    </style>

    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
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
        <div class="form-check">
            <input id="hideWOTorrents" class="form-check-input" type="checkbox" autocomplete="off">
            <label for="hideWOTorrents">Hide without torrents</label>
        </div>

        <button id="search" class="btn btn-primary w-100" type="button">Search</button>
        <br>
        <div id="torrents"></div>
        <br>
        <div id="pagesBlock">
            <ul id="pages" class="pagination justify-content-center flex-wrap">
                <li class="page-item">
                    <button class="page-link" href="#">Previous</button>
                </li>
                <li class="page-item active">
                    <button class="page-link" href="#">1</button>
                </li>
                <li class="page-item">
                    <button class="page-link" href="#">2</button>
                </li>
                <li class="page-item">
                    <button class="page-link" href="#">3</button>
                </li>
                <li class="page-item">
                    <button class="page-link" href="#">Next</button>
                </li>
            </ul>
        </div>
    </div>
    <footer class="page-footer navbar-dark bg-dark">
        <span class="navbar-brand d-flex justify-content-center">
			<a rel="external" style="text-decoration: none;" href="/about">About</a>
			</span>
    </footer>
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
            searchTorrents();
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
            searchTorrents();
        }

        function updateUI() {
            selectSearchType = $('#sType')[0].selectedIndex;
            if (selectSearchType == 0) {
                $('#sbName').show(200);
                $('#sbFilter').hide(200);
            } else if (selectSearchType == 1) {
                $('#sbName').hide(200);
                $('#sbFilter').hide(200);
            } else {
                $('#sbName').hide(200);
                $('#sbFilter').show(200);
            }
            updatePages();
        }

        function updatePages() {
            if (pages == 1) {
                $('#pagesBlock').addClass('hidden');
                return;
            } else
                $('#pagesBlock').removeClass('hidden');
            $('#pages').empty();
            var html = "";
            if (page > 1)
                html += '<li class="page-item"><button class="page-link" onclick="goPage(' + 1 + ')">Previous</button></li>';
            for (i = 1; i <= pages; i++) {
                if (i == page)
                    html += '<li class="page-item active"><button class="page-link">' + i + '</button></li>';
                else
                    html += '<li class="page-item"><button class="page-link" onclick="goPage(' + i + ')">' + i + '</button></li>';
            }
            if (page < pages)
                html += '<li class="page-item"><button class="page-link" onclick="goPage(' + pages + ')">Next</button></li>';
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

        function searchTorrents() {
            $('#search').prop("disabled", true);
            $('#pagesBlock').prop("disabled", true);

            var hide = $('#hideWOTorrents').prop('checked');
            if (typeof page != "number")
                page = 1;

            var SRequest = {
                "Name": $('#sName').val(),
                "Type": selectSearchType,
                "Page": page,
                "HideWTorrent": hide,
				"SearchMovie": $("#SearchMovie").prop('checked'),
				"SearchTV": $("#SearchTV").prop('checked')
            };
            if (selectSearchType == 2) {
                SRequest.Filter = getFilter();
            }

            $.post('/search/request', JSON.stringify(SRequest))
                .done(function(data) {
                    var torrents = $("#torrents");
                    torrents.empty();
                    $('<hr>').appendTo(torrents);
                    pages = data.Pages;
                    for (var key in data.Movies) {
                        var tor = data.Movies[key];
                        var html = '';
                        html += '<div onclick="toggleInfo(\'#torr' + key + '\')" class="movie">';
                        html += ' <img class="poster rounded float-left" src="' + tor.PosterUrl + '"/>';
                        html += ' <div class="description">';
                        html += '  <h4 style="padding-top: 10px;">' + tor.Title + ' / ' + tor.OrigTitle + '</h4>';
                        html += '  <p style="float:left">' + tor.Date + '</p>';
                        var movtype = "Фильм";
                        if (tor.IsTv)
                            movtype = "Сериал";
                        html += '  <p style="float:right">' + movtype + '</p>';
                        html += '  <br><p style="float:right">' + tor.Genres + '</p>';
                        html += ' </div>';
                        html += '</div>';
                        html += '<div style="clear:both"></div>';
                        html += getTorrList(key, tor.Torrents, tor.Overview, tor.BackdropUrl);
                        html += '<hr>';
                        $(html).appendTo(torrents);
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

        function toggleInfo(key) {
            $(key).toggle(50);
        }

        function getTorrList(key, torrList, torrOverview, BackdropUrl) {
            var html = '';
            html += '<div style="background-image: url(' + BackdropUrl + ');" class="hidden torrList" id="torr' + key + '">';
            html += '<p>' + torrOverview + '</p>';
            html += '<div class="btn-group-vertical d-flex" role="group">';
            for (var key in torrList) {
                torr = torrList[key];
                var dl = '';
                if (torr.PeersDl >= 0) {
                    dl = '| ▼ ' + torr.PeersDl;
                    dl += '| ▲ ' + torr.PeersUl;
                }
                html += '<button class="btn-outline-primary button-text-shadow w-100" onclick="doTorrent(\'' + torr.Magnet + '\', this)"><i class="fas fa-plus"></i> ' + torr.Name + " " + torr.Size + dl + '</button>';
            }
            html += '</div>';
            html += '</div>';
            return html;
        }

        function doTorrent(magnet, elem) {
            $(elem).prop("disabled", true);
            var magJS = JSON.stringify({
                Link: magnet
            });
            $.post('/torrent/add', magJS)
                .done(function(data) {
                    $(elem).prop("disabled", false);
                })
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
