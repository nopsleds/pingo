package impl

import (
	"net/http"
)

type HttpProbe struct {
	URL string
}

func (probe *HttpProbe) Test() *ProbeError {
	_, err := http.Get(probe.URL)
	if err != nil {
		return &ProbeError{err.Error()}
	}
	return nil
}
