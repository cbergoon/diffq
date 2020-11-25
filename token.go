package diffq

import "fmt"

type tokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	COMMENT = "COMMENT"

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
	GOESGT    = "=GT>"
	GOESLT    = "=LT>"
	GOESGTE   = "=GTE>"
	GOESLTE   = "=LTE>"
	NIL       = "NIL"
	CREATED   = "$created"
	DELETED   = "$deleted"
)

type token struct {
	ttype    tokenType
	tliteral string
}

func (t *token) String() string {
	return fmt.Sprintf("Type: %s, Literal: %s", t.ttype, t.tliteral)
}

var keywords = map[string]tokenType{
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
	"=gt>":  GOESGT,
	"=lt>":  GOESLT,
	"=gte>": GOESGTE,
	"=lte>": GOESLTE,
	"=GT>":  GOESGT,
	"=LT>":  GOESLT,
	"=GTE>": GOESGTE,
	"=LTE>": GOESLTE,

	"nil": NIL,
	"NIL": NIL,

	"$created": CREATED,
	"$deleted": DELETED,
	"$CREATED": CREATED,
	"$DELETED": DELETED,
}

func lookupIdent(ident string) tokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
