package probe

import (
	"fmt"

	"../config"
)

type ProbeError struct {
	message string
}

func (err *ProbeError) Error() string {
	return err.message
}

type Probe interface {
	Test() *ProbeError
}

var (
	ErrUnreachable    = &ProbeError{message: "unreachble host"}
	ErrUnknown        = &ProbeError{message: "unknown"}
	ErrNotImplemented = &ProbeError{message: "not implemented"}
)

func MakeProbe(target config.Target) (res Probe, err error) {
	switch target.Type {
	case config.TypeHttp:
		return &HttpProbe{URL: target.HttpUrl}, nil
	default:
		return nil, fmt.Errorf("unsupported probe type '%s'", target.Type)
	}
}
