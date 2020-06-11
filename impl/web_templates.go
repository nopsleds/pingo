package impl

import "html/template"

const TPL_INDEX = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf8"/>
	<title>Pingo</title>
	<link rel="stylesheet" href="style.css"/>
</head>
<body>
	<header>
        <div class="appname">Pingo</div>
    </header>
	<main>
		<div class="target">
			<div class="target-title">target</div>
			<div class="target-title">type</div>
			<div class="target-title">status</div>
			{{ if (len .Targets) gt 0 }}
				{{ range $key, $value := (index .Targets 0).Timeseries }}
					<div class="target-title">{{$key}}</div>
				{{ end }}
			{{ end }}
		</div>
		{{ range .Targets }}
			<div class="target target-status-{{ .Class }}">
				<div class="target-name">{{ .Name }}</div>
				<div class="target-type">{{ .Type }}</div>
				<div class="target-status">
					<div class="target-status-value">
						{{ .Status }} - {{ .Reason }}
					</div>
					<div class="target-status-detail">since {{ .LastChange }}</div>
					<div class="target-status-detail">last check {{ .LastCheck }}</div>
				</div>
				{{ range .Timeseries }}
					<div class="target-history">
						<div class="timeserie">
						{{- range . }}
							<div class="timeserie-item timeserie-item-{{ .Value }}">
								<div class="timeserie-item-info">
									value = {{.Value }}<br>
									at {{ .Time }}
								</div>
							</div>
						{{- end -}}
						</div>
					</div>
				{{ end }}
			</div>
		{{ end }}
	</main>
	<footer>current time: {{ .Now }}</footer>
    <script>setInterval(() => location.reload(), 1000)</script>
</body>
</html>
`

const TPL_STYLE_CSS = `

html {
	margin: 0;
    color: #CCC;
    background-color: black;
    font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen-Sans,Ubuntu,Cantarell,"Helvetica Neue",sans-serif;

    --color-ok: rgb(42, 131, 42);
    --color-error: rgb(197, 10, 10);
    --color-unknown: rgb(202, 202, 202, 0.3);
}
body {
    padding: 48px;
}

header {
    margin-bottom: 24px;
}
footer {
    margin-top: 24px;
    font-size: 12px;
    color: #666;
}

.target {
    display: grid;
    grid-gap: 2px;
    padding: 8px;
    line-height: 24px;
    grid-template-columns: repeat(6, 1fr);
    margin-bottom: 2px;
    background-color: rgba(250, 250, 250, 0.1);
}
.target-title {
    line-height: 12px;
    font-size: 12px;
    color: #666;
}
.target-status {
    padding: 4px;
    border-radius: 3px;
    text-align: center;
}
.target.target-status-ok .target-status {
    background-color: var(--color-ok);
}
.target.target-status-error .target-status {
    background-color: var(--color-error);
}

.target-status-value {
    font-weight: bold;
}
.target-status-detail {
	font-size: 12px;
}
.target-name {
    font-weight: bold;
}

.timeserie {
    display: flex;
    width: 240px;
    height: 24px;
}
.timeserie-item {
    border-radius: 3px;
    border: 1px solid black;
    flex: 1;
	position: relative;
}
.timeserie-item:hover {
    border-color: white;
}
.timeserie-item .timeserie-item-info {
	display: none;
}
.timeserie-item:hover .timeserie-item-info {
	display: block;
	position: absolute;
	z-index: 1;
	width: 200px;
	font-size: 12px;
	background-color: black;
	padding: 8px;
	top: 12px;
	left: 12px;
}
.timeserie-item-0 {
    background-color: var(--color-unknown);
}
.timeserie-item-1 {
    background-color: var(--color-ok);
}
.timeserie-item-2 {
    background-color: var(--color-error);
}
`

var (
	tplIndex    = template.Must(template.New("index.html").Parse(TPL_INDEX))
	tplStyleCss = template.Must(template.New("style.css").Parse(TPL_STYLE_CSS))
)
