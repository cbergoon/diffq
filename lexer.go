package diffq

type lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func (l *lexer) nextToken() *token {
	tok := &token{}

	l.skipWhitespace()

	switch l.ch {
	case '/':
		if l.peekChar() == '*' {
			l.readChar()
			tok.tliteral = l.readComment()
			tok.ttype = COMMENT
			l.readChar()
			return tok
		}
		tok = newToken(ILLEGAL, l.ch)
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '[':
		tok = newToken(LBRACKET, l.ch)
	case ']':
		tok = newToken(RBRACKET, l.ch)
	case '"':
		tok.ttype = STRING
		tok.tliteral = l.readString()
	case '$':
		tok.tliteral = l.readIdentifier()
		tok.ttype = lookupIdent(tok.tliteral)
		return tok
	case '=':
		tok.tliteral = l.readOperator()
		tok.ttype = lookupIdent(tok.tliteral)
		return tok
	case 0:
		tok.tliteral = ""
		tok.ttype = EOF
	default:
		if isLetter(l.ch) {
			if l.ch == 'd' && l.peekChar() == '"' {
				l.readChar()
				tok.tliteral = l.readString()
				tok.ttype = DURATION
				l.readChar()
				return tok
			} else if l.ch == 't' && l.peekChar() == '"' {
				l.readChar()
				tok.tliteral = l.readString()
				tok.ttype = TIME
				l.readChar()
				return tok
			}
			tok.tliteral = l.readIdentifier()
			tok.ttype = lookupIdent(tok.tliteral)
			return tok
		} else if isDigit(l.ch) || isNegativeSign(l.ch) {
			tok.tliteral, tok.ttype = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *lexer) readOperator() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || isConcatenator(l.ch) || isSpecial(l.ch) || isEqualSign(l.ch) || isAngleBracket(l.ch) || isExclamationPoint(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || isConcatenator(l.ch) || isSpecial(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *lexer) readNumber() (string, tokenType) {
	position := l.position
	if l.ch == '-' {
		l.readChar()
	}
	seenDecimal := false
	for isDigit(l.ch) || (isDecimal(l.ch) && !seenDecimal) {
		if isDecimal(l.ch) {
			seenDecimal = true
		}
		l.readChar()
	}
	if seenDecimal {
		return l.input[position:l.position], FLOAT
	}
	return l.input[position:l.position], INT
}

func (l *lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *lexer) readComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			break
		}
	}
	return l.input[position : l.position-1]
}

func isConcatenator(ch byte) bool {
	return ch == '_' || ch == '.'
}

func isSpecial(ch byte) bool {
	return ch == '$' || ch == '-' || ch == '*'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isNegativeSign(ch byte) bool {
	return ch == '-'
}

func isDecimal(ch byte) bool {
	return ch == '.'
}

func isEqualSign(ch byte) bool {
	return ch == '='
}

func isAngleBracket(ch byte) bool {
	return ch == '>' || ch == '<'
}

func isExclamationPoint(ch byte) bool {
	return ch == '!'
}

func newToken(tokenType tokenType, ch byte) *token {
	return &token{ttype: tokenType, tliteral: string(ch)}
}
