package compute

const (
	GetCommandString    = "GET"
	SetCommandString    = "SET"
	DeleteCommandString = "DEL"
)

type Command uint8

const (
	InvalidCommand Command = iota
	SetCommand
	GetCommand
	DeleteCommand
)

func CommandFromString(s string) Command {
	switch s {
	case SetCommandString:
		return SetCommand
	case GetCommandString:
		return GetCommand
	case DeleteCommandString:
		return DeleteCommand
	}

	return InvalidCommand
}
