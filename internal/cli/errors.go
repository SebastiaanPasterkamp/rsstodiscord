package cli

import "fmt"

var (
	// ErrParsingFailed is the error returned when the command line arguments
	// could not be parsed
	ErrParsingFailed = fmt.Errorf("CLI parsing failed")
)
