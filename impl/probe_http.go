package impl

import (
	"fmt"
	"net/http"
)

type HttpProbe struct {
	URL            string
	ExpectedStatus int
}

func (probe *HttpProbe) Test() *ProbeError {
	res, err := http.Get(probe.URL)
	if err != nil {
		return &ProbeError{err.Error()}
	}
	if probe.ExpectedStatus > 0 {
		if probe.ExpectedStatus != res.StatusCode {
			msg := fmt.Sprintf("HTTP status %d (expected: %d)", res.StatusCode, probe.ExpectedStatus)
			return &ProbeError{
				message: msg,
			}
		}
	}
	return nil
}
