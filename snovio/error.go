package snovio

type errorCheckEmailStatus struct {
	msg string
}

func newError(msg string) *errorCheckEmailStatus {
	return &errorCheckEmailStatus{msg}
}

func (e *errorCheckEmailStatus) Error() string {
	return e.msg
}
