package goforth

import "strconv"

type LexState struct {
	Line     int
	Column   int
	Position int
	Token    string
	Tokens   []Token
}

func (state *LexState) ConsumeToken() {
	typ := toTokenType(state.Token)
	tok := Token{
		Value:  state.Token,
		Length: len(state.Token),
		Line:   state.Line,
		Column: state.Column,
		Type:   typ,
	}
	if typ == TokenInt {
		i, _ := strconv.Atoi(tok.Value)
		tok.ValueInt = i
	}
	state.Tokens = append(state.Tokens, tok)
	state.Token = ""
}

type TokenType int

const (
	TokenWord = iota
	TokenQuote
	TokenColon
	TokenSemiColon
	TokenInt
	TokenIf
	TokenElse
	TokenThen
	TokenDo
	TokenLoop
	TokenVariable
)

func toTokenType(s string) TokenType {
	if _, err := strconv.Atoi(s); err == nil {
		return TokenInt
	}

	switch s {
	case "\"":
		return TokenQuote
	case ":":
		return TokenColon
	case ";":
		return TokenSemiColon
	case "if":
		return TokenIf
	case "else":
		return TokenElse
	case "then":
		return TokenThen
	case "do":
		return TokenDo
	case "loop":
		return TokenLoop
	case "variable":
		return TokenVariable
	default:
		return TokenWord
	}
}

type Token struct {
	Value    string
	ValueInt int
	Length   int
	Line     int
	Column   int
	Type     TokenType
}

func isNewLine(c uint8) bool {
	return c == '\n'
}

func isWhitespace(c uint8) bool {
	return c == ' ' || c == '\t' || c == '\r' || isNewLine(c)
}

func Tokenize(code string) []Token {
	state := &LexState{}
	codeLen := len(code)
	for state.Position < codeLen {
		char := code[state.Position]
		switch {
		case isWhitespace(char):
			if state.Token != "" {
				state.ConsumeToken()
			}

			if isNewLine(char) {
				state.Column = 0
				state.Line++
			} else {
				state.Column++
			}

			state.Position++
		default:
			state.Token += string(char)
			state.Column++
			state.Position++
		}

	}

	// consume last token
	if state.Token != "" {
		state.ConsumeToken()
	}

	return state.Tokens
}
