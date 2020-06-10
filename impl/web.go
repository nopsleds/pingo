package impl

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"
)

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
			<div class="target-title">last check</div>
			<div class="target-title">status</div>
		</div>
		{{ range .Targets }}
			<div class="target target-status-{{ .Class }}">
				<div class="target-name">{{ .Name }}</div>
				<div class="target-type">{{ .Type }}</div>
				<div class="target-last-check">{{ .LastCheck }}</div>
				<div class="target-status">{{ .Status }} - {{ .Reason }}</div>
			</div>
		{{ end }}
	</main>
	<footer>current time: {{ .Now }}</footer>
    <script>//setInterval(() => location.reload(), 1000)</script>
</body>
</html>
`

const TPL_STYLE_CSS = `
html {
	margin: 0;
    color: #CCC;
    background-color: black;
    font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen-Sans,Ubuntu,Cantarell,"Helvetica Neue",sans-serif;
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
    padding: 8px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    line-height: 24px;
    grid-template-columns: repeat(6, 1fr);
}
.target-title {
    line-height: 12px;
    font-size: 12px;
    color: #666;
}
.target.target-status-ok {
    background-color: rgb(42, 131, 42);
}
.target.target-status-error {
    background-color: rgb(197, 10, 10);
}
.target-name {
    font-weight: bold;
}
`

var (
	tplIndex    = template.Must(template.New("index.html").Parse(TPL_INDEX))
	tplStyleCss = template.Must(template.New("style.css").Parse(TPL_STYLE_CSS))
)

const dateFormat = "Jan 02, 2006 15:04:05 UTC"

type IndexDataTarget struct {
	Name      string
	Type      string
	LastCheck string
	Status    string
	Class     string
	Reason    string
}
type IndexData struct {
	Now     string
	Targets []IndexDataTarget
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByName []IndexDataTarget

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func statusToLabel(status TargetStatus) string {
	switch status {
	case TargetStatusOk:
		return "OK"
	case TargetStatusError:
		return "Error"
	default:
		return "Unknown"
	}
}
func statusToClass(status TargetStatus) string {
	switch status {
	case TargetStatusOk:
		return "ok"
	case TargetStatusError:
		return "error"
	default:
		return "unknown"
	}
}

func handleGetIndex(instance *PingoInstance) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dataTargets []IndexDataTarget
		for targetName, targetState := range instance.State() {
			dataTargets = append(dataTargets, IndexDataTarget{
				Name:      targetName,
				Type:      targetState.Config.Type,
				LastCheck: targetState.LastCheck.Format(dateFormat),
				Status:    statusToLabel(targetState.Status),
				Class:     statusToClass(targetState.Status),
				Reason:    targetState.Reason,
			})
		}
		sort.Sort(ByName(dataTargets))
		data := IndexData{
			Now:     time.Now().Format(dateFormat),
			Targets: dataTargets,
		}
		err := tplIndex.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func handleGetStyleCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/css")
	tplStyleCss.Execute(w, nil)
}

func RunWeb(config ConfigWeb, instance *PingoInstance) error {
	log.Printf("starting web UI on http://localhost:%d\n", config.Port)
	addr := fmt.Sprintf(":%d", config.Port)
	http.HandleFunc("/style.css", handleGetStyleCss)
	http.HandleFunc("/", handleGetIndex(instance))
	return http.ListenAndServe(addr, nil)
}
