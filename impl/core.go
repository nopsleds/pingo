package impl

import (
	"log"
	"time"
)

type TargetStatus string

const (
	TargetStatusUnknown TargetStatus = "Unknown"
	TargetStatusOk      TargetStatus = "OK"
	TargetStatusError   TargetStatus = "Error"
)

type TargetState struct {
	LastCheck time.Time
	Status    TargetStatus
	Reason    string
	Config    ConfigTarget
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
			if targetEntry, ok := this.Targets[result.TargetName]; ok {
				previousStatus := targetEntry.State.Status

				targetEntry.State.LastCheck = time.Now()
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

				if previousStatus != newStatus {
					log.Printf("target %s changed status from %v to %v", result.TargetName, previousStatus, newStatus)
				}

			} else {
				log.Printf("result for unknown target %s\n", result.TargetName)
			}
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
