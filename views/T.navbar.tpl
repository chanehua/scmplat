{{define "navbar"}}
<a class="navbar-brand" href="/"></a>
<div>
	<ul class="nav navbar-nav">
		<li {{if .IsHome}}class="active"{{end}}><a href="/">首页</a></li>
		<li {{if .IsProc}}class="active"{{end}}><a href="/proc?p=1">进程管理</a></li>
		<li {{if .IsPub}}class="active"{{end}}><a href="/pub?p=1">发布管理</a></li>
		<li {{if .IsExec}}class="active"{{end}}><a href="/exec?p=1">启停管理</a></li>
		<li {{if .IsSec}}class="active"{{end}} style="display:{{.SecMg}};"><a href="/sec?p=1">权限管理</a></li>
		<!-- <li {{if .IsdbMg}}class="active"{{end}}><a href="/db">db管理</a></li> -->
	</ul>
</div>

<div class="pull-right">
	<ul class="nav navbar-nav">
		{{if .IsLogin}}
		<li><a href="/login?exit=true">退出登录</a></li>
		{{else}}
		<li><a href="/login">登录</a></li>
		{{end}}
	</ul>
</div>
{{end}}