package diffq

import "fmt"

// tokenType represents the type of a token using predefines constants.
type tokenType string

const (
	// ILLEGAL represents an unsupported token
	ILLEGAL = "ILLEGAL"
	//EOF represents end of file
	EOF = "EOF"

	COMMENT = "COMMENT"

	// Identifiers and Literals

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

// token represents the output of the lexer representing each component of the
// statement as a type and literal pair.
type token struct {
	// ttype represents the type of the token
	ttype tokenType
	// tliteral represents the actual value parsed by the lexer
	tliteral string
}

// String returns a human readable string format of token.
func (t *token) String() string {
	return fmt.Sprintf(`{Type: "%s", Literal: "%s"}`, t.ttype, t.tliteral)
}

// keywords is a lookup map for the token type based on literal value of the
// keyword.
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

// lookupIdent first checks for and returns a matching keyword otherwise returns
// generic IDENT type.
func lookupIdent(ident string) tokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
