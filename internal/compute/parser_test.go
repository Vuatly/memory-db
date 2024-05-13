package compute

import (
	"memory-db/internal/pkg/interror"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParser_ParseQuery(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		query  string
		tokens []string
		err    error
	}{
		"should parse empty query": {
			query: "",
		},
		"should parse empty query with whitespaces": {
			query: "    ",
		},
		"should parse valid query": {
			query:  "SET key value_*/",
			tokens: []string{"SET", "key", "value_*/"},
		},
		"should parse valid query with extra whitespaces": {
			query:  "    SET   key    value ",
			tokens: []string{"SET", "key", "value"},
		},

		"should fail in initial state if contains invalid symbol": {
			query: "&",
			err:   interror.InvalidSymbolError,
		},
		"should fail from whitespace state if contains invalid symbol": {
			query: "SET %",
			err:   interror.InvalidSymbolError,
		},
		"should fail from symbol found state if contains invalid symbol": {
			query: "SET va%",
			err:   interror.InvalidSymbolError,
		},
	}

	parserObj := NewParser(zap.NewNop())

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tokens, err := parserObj.ParseQuery(test.query)
			require.Equal(t, test.err, err)
			require.Equal(t, tokens, test.tokens)
		})
	}
}
