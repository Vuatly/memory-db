package compute

import (
	"memory-db/internal/pkg/interror"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCompute_HandleQuery(t *testing.T) {
	t.Parallel()

	t.Run("should successfully handle query", func(t *testing.T) {
		t.Parallel()
		fx := tearUp(t)
		expQuery := Query{
			command:   GetCommand,
			arguments: []string{"key"},
		}

		fx.parser.EXPECT().ParseQuery("GET key").Return([]string{"GET", "key"}, nil)
		fx.analyzer.EXPECT().AnalyzeTokens([]string{"GET", "key"}).Return(expQuery, nil)

		query, err := fx.compute.HandleQuery("GET key")
		require.NoError(t, err)
		require.Equal(t, expQuery, query)
	})

	t.Run("should return err if parser fails", func(t *testing.T) {
		t.Parallel()
		fx := tearUp(t)

		fx.parser.EXPECT().ParseQuery("get key&").Return(nil, interror.InvalidSymbolError)

		_, err := fx.compute.HandleQuery("get key&")
		require.Error(t, err, interror.InvalidSymbolError)
	})

	t.Run("should return err if analyzer fails with zero tokens error", func(t *testing.T) {
		t.Parallel()
		fx := tearUp(t)

		fx.parser.EXPECT().ParseQuery("").Return([]string{}, nil)
		fx.analyzer.EXPECT().AnalyzeTokens([]string{}).Return(Query{}, interror.ZeroTokensError)

		_, err := fx.compute.HandleQuery("")
		require.Error(t, err, interror.ZeroTokensError)
	})

	t.Run("should return err if analyzer fails with command error", func(t *testing.T) {
		t.Parallel()
		fx := tearUp(t)

		fx.parser.EXPECT().ParseQuery("get key").Return([]string{"get", "key"}, nil)
		fx.analyzer.EXPECT().AnalyzeTokens([]string{"get", "key"}).Return(Query{}, interror.InvalidCommandError)

		_, err := fx.compute.HandleQuery("get key")
		require.Error(t, err, interror.InvalidCommandError)
	})

	t.Run("should return err if analyzer fails with invalid arguments", func(t *testing.T) {
		t.Parallel()
		fx := tearUp(t)

		fx.parser.EXPECT().ParseQuery("get key value").Return([]string{"get", "key", "value"}, nil)
		fx.analyzer.EXPECT().AnalyzeTokens([]string{"get", "key", "value"}).Return(Query{}, interror.InvalidArgumentsNumberError)

		_, err := fx.compute.HandleQuery("get key value")
		require.Error(t, err, interror.InvalidArgumentsNumberError)
	})
}

type fixture struct {
	compute *Compute

	analyzer *AnalyzerMock
	parser   *ParserMock
}

func tearUp(t *testing.T) *fixture {
	analyzerMock := NewAnalyzerMock(t)
	parserMock := NewParserMock(t)

	compute, err := NewCompute(analyzerMock, parserMock, zap.NewNop())
	require.NoError(t, err)

	return &fixture{
		compute:  compute,
		analyzer: analyzerMock,
		parser:   parserMock,
	}
}
