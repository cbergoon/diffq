package diffq

// lexer represents the lexical analyzer that parses/tokenizes the diffq
// language.
type lexer struct {
	// input is the statement to parse.
	input string
	// position is the current position in input (points to the current char)
	position int
	// readPosition is the current position in the input (after the current
	// char)
	readPosition int
	// ch is the current character to parse
	ch byte
}

// newLexer initializes a new lexer with the provided input.
func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

// nextToken advances the parser based on the current character. The current
// character is inspected and the lexer is advanced according to the context and
// the inspection creating tokens based on the character sequence. The next
// token derived from the input is returned with the identified type and the
// literal value.
func (l *lexer) nextToken() *token {
	tok := &token{}

	// whitespace is insignificant in the diffq language other than the
	// separation of tokens.
	l.skipWhitespace()

	switch l.ch {
	// comments
	case '/':
		if l.peekChar() == '*' {
			l.readChar()
			tok.tliteral = l.readComment()
			tok.ttype = cCOMMENT
			l.readChar()
			return tok
		}
		tok = newToken(cILLEGAL, l.ch)
	// asterisk
	case '*':
		tok = newToken(cASTERISK, l.ch)
	// comma
	case ',':
		tok = newToken(cCOMMA, l.ch)
	// left parenthesis
	case '(':
		tok = newToken(cLPAREN, l.ch)
	// right parenthesis
	case ')':
		tok = newToken(cRPAREN, l.ch)
	// left bracket
	case '[':
		tok = newToken(cLBRACKET, l.ch)
	// right bracket
	case ']':
		tok = newToken(cRBRACKET, l.ch)
	// string / quote
	case '"':
		tok.ttype = cSTRING
		tok.tliteral = l.readString()
	// special keywords; $created, $deleted
	case '$':
		tok.tliteral = l.readIdentifier()
		tok.ttype = lookupIdent(tok.tliteral)
		return tok
	// operators
	case '=':
		tok.tliteral = l.readOperator()
		tok.ttype = lookupIdent(tok.tliteral)
		return tok
	// end of file
	case 0:
		tok.tliteral = ""
		tok.ttype = cEOF
	// literals; string, time, duration, integer, float, true, false, etc.
	default:
		if isLetter(l.ch) {
			if l.ch == 'd' && l.peekChar() == '"' {
				l.readChar()
				tok.tliteral = l.readString()
				tok.ttype = cDURATION
				l.readChar()
				return tok
			} else if l.ch == 't' && l.peekChar() == '"' {
				l.readChar()
				tok.tliteral = l.readString()
				tok.ttype = cTIME
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
			tok = newToken(cILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace consumes whitespace blindly. Called on each iteration of
// nextToken().
func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readChar advances the current position of the lexer in the input.
func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar looks ahead to the next character in the input.
func (l *lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// readOperator advances the input until the end of the operator. The start of
// the operator is indicated by the '=' character.
func (l *lexer) readOperator() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || isConcatenator(l.ch) || isSpecial(l.ch) || isEqualSign(l.ch) || isAngleBracket(l.ch) || isExclamationPoint(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readIdentifier advances the input until the end of the identifier.
func (l *lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || isConcatenator(l.ch) || isSpecial(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber advances the input until the end of the number.
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
		return l.input[position:l.position], cFLOAT
	}
	return l.input[position:l.position], cINT
}

// readString advances the input until the end of the string. The start and end
// is indicated by a '"'.
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

// readComment advances the input until the end of the comment. The start and
// end of the comment are indicated by "/*" and "*/" respectively.
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

// newToken creates a new token with the type and literal value specified.
func newToken(tokenType tokenType, ch byte) *token {
	return &token{ttype: tokenType, tliteral: string(ch)}
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
