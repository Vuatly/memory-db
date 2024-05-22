package compute

type Query struct {
	command   Command
	arguments []string
}

func NewQuery(command Command, arguments []string) Query {
	return Query{
		command:   command,
		arguments: arguments,
	}
}

func (q *Query) GetCommand() Command {
	return q.command
}

func (q *Query) GetArguments() []string {
	return q.arguments
}
