package probe

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
