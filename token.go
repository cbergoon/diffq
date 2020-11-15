package diffq

import "fmt"

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT    = "IDENT"    // field, field.val, array.0.val
	INT      = "INT"      // 1343456, -123456
	STRING   = "STRING"   // "foobar"
	FLOAT    = "FLOAT"    // 123.456, -123.456
	DURATION = "DURATION" // d"12h30m"
	TIME     = "TIME"     // t"2006-01-02T15:04:05+07:00" t"2006-01-02T15:04:05Z" (time.RFC3339)

	ASTERISK = "*"

	// Delimiters
	COMMA = ","

	LPAREN = "("
	RPAREN = ")"

	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	AND       = "AND"
	OR        = "OR"
	EVAL      = "EVAL"
	GOESTO    = "=>"
	NOTGOESTO = "=!>"
	// GOESGT    = "=GT>"
	// GOESLT    = "=LT>"
	// GOESGTE   = "=GTE>"
	// GOESLTE   = "=LTE>"
	NIL     = "NIL"
	CREATED = "$created"
	DELETED = "$deleted"
)

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) String() string {
	return fmt.Sprintf("Type: %s, Literal: %s", t.Type, t.Literal)
}

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"TRUE":  TRUE,
	"FALSE": FALSE,
	"or":    OR,
	"and":   AND,
	"eval":  EVAL,
	"OR":    OR,
	"AND":   AND,
	"EVAL":  EVAL,
	"=>":    GOESTO,
	"=!>":   NOTGOESTO,

	// "=GT>":  GOESGT,
	// "=LT>":  GOESLT,
	// "=GTE>": GOESGTE,
	// "=LTE>": GOESLTE,

	"nil": NIL,

	"$created": CREATED,
	"$deleted": DELETED,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
