//go:generate mockery --case underscore --name parser --structname ParserMock --with-expecter --inpackage --inpackage-suffix
//go:generate mockery --case underscore --name analyzer --structname AnalyzerMock --with-expecter --inpackage --inpackage-suffix

package compute

import (
	"errors"

	"go.uber.org/zap"
)

type analyzer interface {
	AnalyzeTokens(tokens []string) (Query, error)
}

type parser interface {
	ParseQuery(query string) ([]string, error)
}

type Compute struct {
	analyzer analyzer
	parser   parser
	logger   *zap.Logger
}

func NewCompute(analyzer analyzer, parser parser, logger *zap.Logger) (*Compute, error) {
	if analyzer == nil {
		return nil, errors.New("analyzer is required")
	}

	if parser == nil {
		return nil, errors.New("parser is required")
	}

	return &Compute{
		analyzer: analyzer,
		parser:   parser,
		logger:   logger,
	}, nil
}

func (c *Compute) HandleQuery(query string) (Query, error) {
	tokens, err := c.parser.ParseQuery(query)
	c.logger.Debug("tokens parsed", zap.Strings("tokens", tokens))
	if err != nil {
		return Query{}, err
	}

	return c.analyzer.AnalyzeTokens(tokens)
}
