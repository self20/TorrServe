package templates

import "server/version"

var settingsPage = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link href="` + faviconB64 + `" rel="icon" type="image/x-icon">
	<link rel="stylesheet" href="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css">
	<script src="http://code.jquery.com/jquery-1.11.3.min.js"></script>
	<script src="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js"></script>
	<title>Torrent settings</title>
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

		<form id="settings">
			<div class="ui-widget">
				<label for="CacheSize">Cache Size</label>
				<input id="CacheSize" type="number" autocomplete="off">
			</div>
	
			<div class="ui-widget">
				<label for="PreloadBufferSize">Preload Buffer Size</label>
				<input id="PreloadBufferSize" type="number" autocomplete="off">
			</div>
			<div class="ui-widget">
				<label for="DisableTCP">Disable TCP</label>
				<input id="DisableTCP" type="checkbox" autocomplete="off">
			</div>
			<div class="ui-widget">
				<label for="DisableUTP">Disable UTP</label>
				<input id="DisableUTP" type="checkbox" autocomplete="off">
			</div>
			<div class="ui-widget">
				<label for="DisableUPNP">Disable UPNP</label>
				<input id="DisableUPNP" type="checkbox" autocomplete="off">
			</div>
			<div class="ui-widget">
				<label for="DisableDHT">Disable DHT</label>
				<input id="DisableDHT" type="checkbox" autocomplete="off">
			</div>
			<div class="ui-widget">	
				<label for="DisableUpload">Disable Upload</label>
				<input id="DisableUpload" type="checkbox" autocomplete="off">
			</div>
			<div class="ui-widget">
				<label for="Encryption">Encryption</label>
				<input id="Encryption" type="number" autocomplete="off">
			</div>

			<div class="ui-widget">
				<label for="ConnectionsLimit">Connections Limit</label>
				<input id="ConnectionsLimit" type="number" autocomplete="off">
			</div>
			<h4>Download/Upload Rate Limit setup in kilobytes, 0 for infinite</h4>
			<div class="ui-widget">
				<label for="DownloadRateLimit">Download Rate Limit</label>
				<input id="DownloadRateLimit" type="number" autocomplete="off">
			</div>
			<div class="ui-widget">
				<label for="UploadRateLimit">Upload Rate Limit</label>
				<input id="UploadRateLimit" type="number" autocomplete="off">
			</div>
			<h4>0 - Do nothing; 1 - Add retrackers; 2 - Remove retrackers</h4>
			<div class="ui-widget">
				<label for="RetrackersMode">Retrackers Mode</label>
				<input id="RetrackersMode" type="number" autocomplete="off">
			</div>
		</form>

		<br>
		<div class="ui-grid-a">
			<div class="ui-block-a"><button id="buttonSave" data-icon="check" onclick="saveSettings()">Save</button></div>
			<div class="ui-block-b"><button id="buttonRefresh" data-icon="refresh" onclick="refreshSettings()">Refresh</button></div>
		</div>

	</div>

	<div data-role="footer">
	<center><p><a rel="external" style="text-decoration: none;" href="/about">About</a></p></center>
	</div>
</div> 

<script>
	function saveSettings(){
		var data = {};
		$('input', '#settings').each(function() {
			switch($(this).prop("type")) { 
					case "checkbox":   
						data[$(this).prop("id")] = $(this).prop("checked");   
						break;  
					default:
						data[$(this).prop("id")] = parseInt($(this).val(),10);
				}
		});
		
		$.post("/settings/write", JSON.stringify(data))
			.done(function( data ) {
				alert(data);
			})
			.fail(function( data ){
				alert(data.responseJSON.message);
			});
	}

	function refreshSettings(){
		$.post("/settings/read")
		.done(function(data){
			var frm = '#settings';
			$.each(data, function(key, value) {  
				var ctrl = $('[id='+key+']', frm);
				switch(ctrl.prop("type")) { 
					case "checkbox":   
						ctrl.prop("checked", value).checkboxradio('refresh');
						break;  
					default:
						ctrl.val(value); 
				}  
			});
		});
	};

	$( document ).ready(function() {
		refreshSettings();
	});
</script>
</body>
</html>
`

func (t *Template) parseSettingsPage() {
	parsePage(t, "settingsPage", settingsPage)
}
