package diffq

import "fmt"

// tokenType represents the type of a token using predefines constants.
type tokenType string

const (
	// cILLEGAL represents an unsupported token
	cILLEGAL = "ILLEGAL"
	//cEOF represents end of file
	cEOF = "EOF"

	cCOMMENT = "COMMENT"

	// Identifiers and Literals

	cIDENT    = "IDENT"    // field, field.val, array.0.val
	cINT      = "INT"      // 1343456, -123456
	cSTRING   = "STRING"   // "foobar"
	cFLOAT    = "FLOAT"    // 123.456, -123.456
	cDURATION = "DURATION" // d"12h30m"
	cTIME     = "TIME"     // t"2006-01-02T15:04:05+07:00" t"2006-01-02T15:04:05Z" (time.RFC3339)

	cASTERISK = "*"

	// Delimiters

	cCOMMA = ","

	cLPAREN = "("
	cRPAREN = ")"

	cLBRACKET = "["
	cRBRACKET = "]"

	// Keywords

	cTRUE      = "TRUE"
	cFALSE     = "FALSE"
	cAND       = "AND"
	cOR        = "OR"
	cEVAL      = "EVAL"
	cGOESTO    = "=>"
	cNOTGOESTO = "=!>"
	cGOESGT    = "=GT>"
	cGOESLT    = "=LT>"
	cGOESGTE   = "=GTE>"
	cGOESLTE   = "=LTE>"
	cNIL       = "NIL"
	cCREATED   = "$created"
	cDELETED   = "$deleted"
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
	"true":  cTRUE,
	"false": cFALSE,
	"TRUE":  cTRUE,
	"FALSE": cFALSE,
	"or":    cOR,
	"and":   cAND,
	"eval":  cEVAL,
	"OR":    cOR,
	"AND":   cAND,
	"EVAL":  cEVAL,

	"=>":    cGOESTO,
	"=!>":   cNOTGOESTO,
	"=gt>":  cGOESGT,
	"=lt>":  cGOESLT,
	"=gte>": cGOESGTE,
	"=lte>": cGOESLTE,
	"=GT>":  cGOESGT,
	"=LT>":  cGOESLT,
	"=GTE>": cGOESGTE,
	"=LTE>": cGOESLTE,

	"nil": cNIL,
	"NIL": cNIL,

	"$created": cCREATED,
	"$deleted": cDELETED,
	"$CREATED": cCREATED,
	"$DELETED": cDELETED,
}

// lookupIdent first checks for and returns a matching keyword otherwise returns
// generic IDENT type.
func lookupIdent(ident string) tokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return cIDENT
}
