package diffq

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() *Token {
	tok := &Token{}

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '>' {
			ch1 := l.ch
			l.readChar()
			tok = &Token{Type: GOESTO, Literal: string(ch1) + string(l.ch)}
		} else if l.peekChar() == '!' {
			ch1 := l.ch
			l.readChar()
			if l.peekChar() == '>' {
				ch2 := l.ch
				l.readChar()
				tok = &Token{Type: NOTGOESTO, Literal: string(ch1) + string(ch2) + string(l.ch)}
			}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			if l.ch == 'd' && l.peekChar() == '"' {
				l.readChar()
				tok.Literal = l.readString()
				tok.Type = DURATION
				l.readChar()
				return tok
			} else if l.ch == 't' && l.peekChar() == '"' {
				l.readChar()
				tok.Literal = l.readString()
				tok.Type = TIME
				l.readChar()
				return tok
			}
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) || isNegativeSign(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || isConcatenator(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, TokenType) {
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

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isConcatenator(ch byte) bool {
	return ch == '_' || ch == '.'
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

func newToken(tokenType TokenType, ch byte) *Token {
	return &Token{Type: tokenType, Literal: string(ch)}
}
