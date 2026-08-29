package tools

import (
	"errors"
	"strconv"
)

func evalExpr(s string) (float64, error) {
	p := &parser{toks: tokenize(s)}
	v, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	if p.pos < len(p.toks) {
		return 0, errors.New("unexpected token: " + p.toks[p.pos].lex)
	}
	return v, nil
}

type tokKind int

const (
	tokNum tokKind = iota
	tokOp
	tokLParen
	tokRParen
)

type token struct {
	kind tokKind
	lex  string
	val  float64
}

func tokenize(s string) []token {
	var toks []token
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case c == ' ' || c == '\t':
			i++
		case c >= '0' && c <= '9' || c == '.':
			j := i
			for j < len(s) && (s[j] >= '0' && s[j] <= '9' || s[j] == '.') {
				j++
			}
			f, err := strconv.ParseFloat(s[i:j], 64)
			if err != nil {
				return nil
			}
			toks = append(toks, token{kind: tokNum, val: f, lex: s[i:j]})
			i = j
		case c == '+' || c == '-' || c == '*' || c == '/':
			toks = append(toks, token{kind: tokOp, lex: string(c)})
			i++
		case c == '(':
			toks = append(toks, token{kind: tokLParen, lex: "("})
			i++
		case c == ')':
			toks = append(toks, token{kind: tokRParen, lex: ")"})
			i++
		default:
			return nil
		}
	}
	return toks
}

// parser is a tiny recursive-descent parser for + - * / ( ).
type parser struct {
	toks []token
	pos  int
}

func (p *parser) peek() (token, bool) {
	if p.pos >= len(p.toks) {
		return token{}, false
	}
	return p.toks[p.pos], true
}

func (p *parser) next() token {
	t := p.toks[p.pos]
	p.pos++
	return t
}

func (p *parser) parseExpr() (float64, error) {
	v, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for {
		t, ok := p.peek()
		if !ok || t.kind != tokOp || (t.lex != "+" && t.lex != "-") {
			return v, nil
		}
		p.next()
		rhs, err := p.parseTerm()
		if err != nil {
			return 0, err
		}
		if t.lex == "+" {
			v += rhs
		} else {
			v -= rhs
		}
	}
}

func (p *parser) parseTerm() (float64, error) {
	v, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for {
		t, ok := p.peek()
		if !ok || t.kind != tokOp || (t.lex != "*" && t.lex != "/") {
			return v, nil
		}
		p.next()
		rhs, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		if t.lex == "*" {
			v *= rhs
		} else {
			if rhs == 0 {
				return 0, errors.New("division by zero")
			}
			v /= rhs
		}
	}
}

func (p *parser) parseFactor() (float64, error) {
	t, ok := p.peek()
	if !ok {
		return 0, errors.New("unexpected end of expression")
	}
	if t.kind == tokNum {
		p.next()
		return t.val, nil
	}
	if t.kind == tokLParen {
		p.next()
		v, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		cl, ok := p.peek()
		if !ok || cl.kind != tokRParen {
			return 0, errors.New("missing closing parenthesis")
		}
		p.next()
		return v, nil
	}
	return 0, errors.New("unexpected token: " + t.lex)
}
