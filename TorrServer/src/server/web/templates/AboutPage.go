package templates

import "server/version"

var aboutPage = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link href="` + faviconB64 + `" rel="icon" type="image/x-icon">
	<link rel="stylesheet" href="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css">
	<script src="http://code.jquery.com/jquery-1.11.3.min.js"></script>
	<script src="http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js"></script>
	<title>About</title>
</head>
<body>
<style type="text/css">
.inline{
	display:inline;
	padding-left: 2%;
}
.center {
    display: block;
    margin-left: auto;
    margin-right: auto;
}
</style>
<div data-role="page">
	<div data-role="header"><h3>About</h3></div>
	<div data-role="content">
		<img class="center" src='` + faviconB64 + `'/>
		<h3 align="middle">TorrServer</h3>
		<h4 align="middle">` + version.Version + `</h4>
		
		<h4>Поддержка проекта:</h4>
		<a class="inline" target="_blank" href="https://www.paypal.me/yourok">PayPal</a>
		<br>
		<a class="inline" target="_blank" href="https://money.yandex.ru/to/410013733697114/100">Yandex.Деньги</a>
		<br>
		<hr align="left" width="25%">
		<br>
		
		<h4>Инструкция по использованию:</h4>
			<a class="inline" target="_blank" href="https://4pda.ru/forum/index.php?showtopic=896840&st=0#entry72570782">4pda.ru</a>
			<p class="inline">Спасибо <b>MadAndron</b></p> 
		<br>
		<hr align="left" width="25%">
		<br>
		
		<h4>Автор:</h4>
			<b class="inline">YouROK</b>
			<br>
			<i class="inline">Email:</i>
			<a target="_blank" class="inline" href="mailto:8yourok8@gmail.com">8YouROK8@gmail.com</a>
			<br>
			<i class="inline">Site: </i>
			<a target="_blank" class="inline" href="https://github.com/YouROK">GitHub.com/YouROK</a>
		<br>
		<hr align="left" width="25%">
		<br>
		
		<h4>Спасибо всем, кто тестировал и помогал:</h4>
			<b class="inline">kuzzman</b>
			<br>
			<i class="inline">Site: </i>
			<a target="_blank" class="inline" href="https://4pda.ru/forum/index.php?showuser=1259550">4pda.ru</a>
			<a target="_blank" class="inline" href="http://tv-box.pp.ua">tv-box.pp.ua</a>
		<br>
		<br>
			<b class="inline">MadAndron</b>
			<br>
			<i class="inline">Site:</i>
			<a target="_blank" class="inline" href="https://4pda.ru/forum/index.php?showuser=1543999">4pda.ru</a>
		<br>
		<br>
			<b class="inline">SpAwN_LMG</b>
			<br>
			<i class="inline">Site:</i>
			<a target="_blank" class="inline" href="https://4pda.ru/forum/index.php?showuser=700929">4pda.ru</a>
		<br>
		<br>
			<b class="inline">Zivio</b>
			<br>
			<i class="inline">Site:</i>
			<a target="_blank" class="inline" href="https://4pda.ru/forum/index.php?showuser=1195633">4pda.ru</a>
			<a target="_blank" class="inline" href="http://forum.hdtv.ru/index.php?showtopic=19020">forum.hdtv.ru</a>
		
		<br>
		<br>
		<hr align="left" width="25%">
		<br>
	</div>
	<div data-role="footer">
		<center><h4>TorrServer ` + version.Version + `</h4></center>
	</div>
</div> 
</body>
</html>
`

func (t *Template) parseAboutPage() {
	parsePage(t, "aboutPage", aboutPage)
}
