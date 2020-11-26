package diffq

import "testing"

func TestLookupIdent(t *testing.T) {
	tt := lookupIdent("$created")
	if tt != CREATED {
		t.Errorf("incorrect identifier found from lookup, got: %s, want: %s", tt, CREATED)
	}
	tt = lookupIdent("test.identifier")
	if tt != IDENT {
		t.Errorf("incorrect identifier found from lookup, got: %s, want: %s", tt, IDENT)
	}
}

func TestTokenString(t *testing.T) {
	tok := token{
		ttype:    AND,
		tliteral: "AND",
	}
	tokstr := tok.String()
	if tokstr != `{Type: "AND", Literal: "AND"}` {
		t.Errorf("incorrect string for type token, got: %s, want: %s", tokstr, "Type: AND, Literal: AND")
	}
}
