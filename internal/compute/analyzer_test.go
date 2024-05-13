package compute

import (
	"memory-db/internal/pkg/interror"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAnalyzer_AnalyzeTokens(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tokens   []string
		expQuery Query
		err      error
	}{
		"should successfully analyze set command": {
			tokens: []string{"SET", "key", "test"},
			expQuery: Query{
				command:   SetCommand,
				arguments: []string{"key", "test"},
			},
		},
		"should successfully analyze get command": {
			tokens: []string{"GET", "key"},
			expQuery: Query{
				command:   GetCommand,
				arguments: []string{"key"},
			},
		},
		"should successfully analyze del command": {
			tokens: []string{"DEL", "key"},
			expQuery: Query{
				command:   DeleteCommand,
				arguments: []string{"key"},
			},
		},

		"should return err if invalid number of arguments for set": {
			tokens: []string{"SET", "key"},
			err:    interror.InvalidArgumentsNumberError,
		},
		"should return err if invalid number of arguments for get": {
			tokens: []string{"GET", "key", "val"},
			err:    interror.InvalidArgumentsNumberError,
		},
		"should return err if invalid number of arguments for del": {
			tokens: []string{"DEL", "key", "val"},
			err:    interror.InvalidArgumentsNumberError,
		},

		"should return err if zero tokens provided": {
			tokens: []string{},
			err:    interror.ZeroTokensError,
		},
		"should return err if invalid command provided": {
			tokens: []string{"set", "key"},
			err:    interror.InvalidCommandError,
		},
	}

	analyzerObj := NewAnalyzer(zap.NewNop())

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			query, err := analyzerObj.AnalyzeTokens(test.tokens)
			require.Equal(t, test.err, err)
			require.Equal(t, test.expQuery, query)
		})
	}
}
