package errs

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrInvalidTargetUrl = Error("revers proxy: invalid target url")
)

