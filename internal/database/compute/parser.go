package compute

import (
	"unicode"

	"go.uber.org/zap"
)

var validPunctuationSymbols = []rune{'*', '_', '/'}

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

func (p *Parser) ParseQuery(query string) ([]string, error) {
	machine := newFSM()
	return machine.parse(query)
}

func isWhitespace(symbol rune) bool {
	return unicode.IsSpace(symbol)
}

func isValidSymbol(symbol rune) bool {
	return unicode.IsLetter(symbol) || unicode.IsNumber(symbol) || isValidPunctuation(symbol)
}

func isValidPunctuation(symbol rune) bool {
	for _, s := range validPunctuationSymbols {
		if symbol == s {
			return true
		}
	}

	return false
}
