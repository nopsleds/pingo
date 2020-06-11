package impl

import (
	"log"
	"time"

	"github.com/nopsleds/pingo/impl/timeserie"
)

type TargetStatus int

const (
	TargetStatusUnknown TargetStatus = 0
	TargetStatusOk      TargetStatus = 1
	TargetStatusError   TargetStatus = 2
)

type TargetState struct {
	Config     ConfigTarget
	LastCheck  time.Time
	LastChange time.Time
	Status     TargetStatus
	Reason     string
	Timeseries map[string]*timeserie.Timeserie
}

type TargetEntry struct {
	Probe           Probe
	State           TargetState
	PollingInterval time.Duration
}
type PingoInstance struct {
	Targets map[string]*TargetEntry
}

func NewPingoInstance(config Config) (*PingoInstance, error) {
	targets := make(map[string]*TargetEntry)
	for targetName, targetConfig := range config.Targets {
		targetState, err := processTarget(targetConfig)
		if err != nil {
			return nil, err
		}
		targets[targetName] = targetState
	}

	return &PingoInstance{
		Targets: targets,
	}, nil
}

func TargetTimeserie(length int, ticker timeserie.Ticker) *timeserie.Timeserie {
	ts, err := timeserie.New(timeserie.Options{
		DefaultValue: int(TargetStatusUnknown),
		Ticker:       ticker,
		Length:       length,
		ValueReducer: func(current, newVal int) int {
			if newVal > current {
				return newVal
			}
			return current
		},
	})
	if err != nil {
		panic(err)
	}
	return ts
}

func processTarget(targetConfig ConfigTarget) (*TargetEntry, error) {
	// polling interval
	pollingInterval, err := time.ParseDuration(targetConfig.PollingInterval)
	if err != nil {
		return nil, err
	}
	probe, err := MakeProbe(targetConfig)
	if err != nil {
		return nil, err
	}
	return &TargetEntry{
		State: TargetState{
			Status: TargetStatusUnknown,
			Config: targetConfig,
			Timeseries: map[string]*timeserie.Timeserie{
				"last 60m": TargetTimeserie(12, timeserie.TickByMinutes(5)),
				"last 24h": TargetTimeserie(24, timeserie.TickByHours(1)),
				"last 7d":  TargetTimeserie(7, timeserie.TickByDay),
			},
		},
		Probe:           probe,
		PollingInterval: pollingInterval,
	}, nil
}

type targetResult struct {
	TargetName string
	Result     *ProbeError
}

func (this *PingoInstance) Run() error {
	results := make(chan targetResult)

	for targetName, t := range this.Targets {
		log.Printf("launching target '%s' > polling every %s", targetName, t.PollingInterval)
		go pollTarget(targetName, t.Probe, t.PollingInterval, results)
	}

	// handling results...
	go func() {
		for result := range results {
			this.processResult(result)
		}
	}()

	this.LogState()
	return nil
}

func pollTarget(targetName string, probe Probe, pollingInterval time.Duration, results chan targetResult) {
	tick := time.Tick(pollingInterval)
	for range tick {
		result := probe.Test()
		results <- targetResult{
			TargetName: targetName,
			Result:     result,
		}
	}
}

func (this *PingoInstance) processResult(result targetResult) {
	if targetEntry, ok := this.Targets[result.TargetName]; ok {
		previousStatus := targetEntry.State.Status

		now := time.Now()
		targetEntry.State.LastCheck = now
		newStatus := TargetStatusUnknown
		newReason := ""
		if result.Result != nil {
			newStatus = TargetStatusError
			newReason = result.Result.Error()
		} else {
			newStatus = TargetStatusOk
		}

		// save
		targetEntry.State.Status = newStatus
		targetEntry.State.Reason = newReason
		for _, ts := range targetEntry.State.Timeseries {
			ts.Insert(int(newStatus))
		}

		if previousStatus != newStatus {
			log.Printf("target %s changed status from %v to %v", result.TargetName, previousStatus, newStatus)
			targetEntry.State.LastChange = now
		}

	} else {
		log.Printf("result for unknown target %s\n", result.TargetName)
	}
}

func (this *PingoInstance) State() map[string]TargetState {
	res := make(map[string]TargetState)
	for targetName, targetEntry := range this.Targets {
		res[targetName] = targetEntry.State
	}
	return res
}

func (this *PingoInstance) LogState() {
	log.Printf("======")
	for targetName, t := range this.Targets {
		log.Printf("target %s: %+v", targetName, t)
	}
	log.Printf("======")
}
