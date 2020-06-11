package impl

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/nopsleds/pingo/impl/timeserie"
)

const dateFormat = "Jan 02, 2006 15:04:05 UTC"

type IndexDataTarget struct {
	Name       string
	Type       string
	LastCheck  string
	LastChange string
	Status     string
	Class      string
	Reason     string
	Timeseries map[string][]timeserie.Entry
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
			timeseries := make(map[string][]timeserie.Entry)
			for k, ts := range targetState.Timeseries {
				timeseries[k] = ts.Get()
			}
			dataTargets = append(dataTargets, IndexDataTarget{
				Name:       targetName,
				Type:       targetState.Config.Type,
				LastCheck:  targetState.LastCheck.Format(dateFormat),
				LastChange: targetState.LastChange.Format(dateFormat),
				Status:     statusToLabel(targetState.Status),
				Class:      statusToClass(targetState.Status),
				Reason:     targetState.Reason,
				Timeseries: timeseries,
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
