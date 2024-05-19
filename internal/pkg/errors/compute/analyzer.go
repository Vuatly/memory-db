package compute

import "errors"

var (
	ZeroTokensError             = errors.New("zero tokens")
	InvalidCommandError         = errors.New("invalid command")
	InvalidArgumentsNumberError = errors.New("invalid number of arguments")
)
