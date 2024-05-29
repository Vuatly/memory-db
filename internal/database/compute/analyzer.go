package compute

import (
	interrors "memory-db/internal/errors/compute"

	"go.uber.org/zap"
)

var commandArgumentNumber = map[Command]int{
	GetCommand:    1,
	SetCommand:    2,
	DeleteCommand: 1,
}

type Analyzer struct {
	logger *zap.Logger
}

func NewAnalyzer(logger *zap.Logger) *Analyzer {
	return &Analyzer{
		logger: logger,
	}
}

func (a *Analyzer) AnalyzeTokens(tokens []string) (Query, error) {
	if len(tokens) == 0 {
		a.logger.Debug("no tokens provided")
		return Query{}, interrors.ZeroTokensError
	}

	command := CommandFromString(tokens[0])
	if command == InvalidCommand {
		a.logger.Debug("invalid command", zap.String("command", tokens[0]))
		return Query{}, interrors.InvalidCommandError
	}

	arguments := tokens[1:]
	if len(arguments) != commandArgumentNumber[command] {
		return Query{}, interrors.InvalidArgumentsNumberError
	}

	query := NewQuery(command, arguments)
	return query, nil
}
