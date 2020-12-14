package goforth

import (
	"fmt"
	"log"
)

const (
	ExprExec = iota
	ExprPush
	ExprDefine
	ExprIfThen
	ExprLoop
	ExprVariable
)

type ExprType int

type Expr struct {
	Type     ExprType
	Value    string
	ValueInt int
	Body     AST
	AltBody  AST
}

type AST []Expr

func (expr Expr) ToString() string {
	switch expr.Type {
	case ExprPush:
		return fmt.Sprintf("( push %+v )", expr.ValueInt)
	case ExprExec:
		return fmt.Sprintf("( %+v )", expr.Value)
	case ExprDefine:
		return fmt.Sprintf("( : %+v %+v )", expr.Value, expr.Body.ToString())
	case ExprIfThen:
		return fmt.Sprintf("( if %+v else %+v then )", expr.Body.ToString(), expr.AltBody.ToString())
	case ExprLoop:
		return fmt.Sprintf("( loop %+v )", expr.Body.ToString())
	case ExprVariable:
		return fmt.Sprintf("( variable %+v )", expr.Value)
	default:
		return "( UNKNOWN )"
	}
}

func (ast AST) ToString() string {
	var s string
	for _, expr := range ast {
		s += expr.ToString()
	}
	return s
}

type ParseState struct {
	AST      AST
	Tokens   []Token
	Position int
}

func (s *ParseState) parsePush() Expr {
	return Expr{
		Type:     ExprPush,
		ValueInt: s.consumeToken().ValueInt,
	}
}

func (s *ParseState) parseExec() Expr {
	return Expr{
		Type:  ExprExec,
		Value: s.consumeToken().Value,
	}
}

func (s *ParseState) parseIfThen() Expr {
	expr := Expr{
		Type: ExprIfThen,
	}
	// consume the `if`
	s.consumeToken()
	body := &ParseState{}
	tok := s.consumeToken()
	for tok.Type != TokenElse && tok.Type != TokenThen {
		body.Tokens = append(body.Tokens, tok)
		tok = s.consumeToken()
	}

	expr.Body = body.parse()

	if tok.Type == TokenThen {
		return expr
	}

	// alternative body
	altBody := &ParseState{}
	tok = s.consumeToken()
	for tok.Type != TokenThen {
		altBody.Tokens = append(altBody.Tokens, tok)
		tok = s.consumeToken()
	}
	expr.AltBody = altBody.parse()

	return expr
}

func (s *ParseState) parseDefinition() Expr {
	s.consumeToken()
	name := s.consumeToken()
	tok := s.consumeToken()
	body := &ParseState{}
	for tok.Type != TokenSemiColon {
		body.Tokens = append(body.Tokens, tok)
		tok = s.consumeToken()
	}

	return Expr{
		Type:  ExprDefine,
		Value: name.Value,
		Body:  body.parse(),
	}
}

func (s *ParseState) parseLoop() Expr {
	s.consumeToken()
	tok := s.consumeToken()
	body := &ParseState{}
	for tok.Type != TokenLoop {
		body.Tokens = append(body.Tokens, tok)
		tok = s.consumeToken()
	}

	return Expr{
		Type: ExprLoop,
		Body: body.parse(),
	}
}

func (s *ParseState) parseVariable() Expr {
	s.consumeToken()
	tok := s.consumeToken()

	return Expr{
		Type:  ExprVariable,
		Value: tok.Value,
	}
}

func (s *ParseState) consumeToken() Token {
	token := s.nextToken()
	s.Position++
	return token
}

func (s *ParseState) nextToken() Token {
	return s.Tokens[s.Position]
}

func (s *ParseState) parse() AST {
	var expr AST
	for s.Position < len(s.Tokens) {
		t := s.nextToken()
		switch t.Type {
		case TokenColon:
			expr = append(expr, s.parseDefinition())
		case TokenInt:
			expr = append(expr, s.parsePush())
		case TokenWord:
			expr = append(expr, s.parseExec())
		case TokenIf:
			expr = append(expr, s.parseIfThen())
		case TokenDo:
			expr = append(expr, s.parseLoop())
		case TokenVariable:
			expr = append(expr, s.parseVariable())
		default:
			log.Fatalf("Don't know about expression %+v", expr)
		}
	}
	return expr
}

func Parse(tokens []Token) AST {
	parser := &ParseState{
		Tokens: tokens,
	}
	return parser.parse()
}
