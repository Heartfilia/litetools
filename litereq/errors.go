package litereq

import "strconv"

// ErrorKind indicates where an error was returned in the process of building, validating, and handling a request.
// Errors returned by Builder can be tested for their ErrorKind using errors.Is or errors.As.
type ErrorKind int8

//go:generate stringer -type=ErrorKind

// Enum values for type ErrorKind
const (
	ErrURL       ErrorKind = iota // error building URL
	ErrRequest                    // error building the request
	ErrTransport                  // error connecting
	ErrValidator                  // validator error
	ErrHandler                    // handler error
)

func (ek ErrorKind) Error() string {
	return ek.String()
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrURL-0]
	_ = x[ErrRequest-1]
	_ = x[ErrTransport-2]
	_ = x[ErrValidator-3]
	_ = x[ErrHandler-4]
}

const errorKindName = "ErrURLErrRequestErrTransportErrValidatorErrHandler"

var errorKindIndex = [...]uint8{0, 6, 16, 28, 40, 50}

func (ek ErrorKind) String() string {
	if ek < 0 || ek >= ErrorKind(len(errorKindIndex)-1) {
		return "ErrorKind(" + strconv.FormatInt(int64(ek), 10) + ")"
	}
	return errorKindName[errorKindIndex[ek]:errorKindIndex[ek+1]]
}
