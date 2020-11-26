package diffq

import "testing"

func TestLookupIdent(t *testing.T) {
	tt := lookupIdent("$created")
	if tt != cCREATED {
		t.Errorf("incorrect identifier found from lookup, got: %s, want: %s", tt, cCREATED)
	}
	tt = lookupIdent("test.identifier")
	if tt != cIDENT {
		t.Errorf("incorrect identifier found from lookup, got: %s, want: %s", tt, cIDENT)
	}
}

func TestTokenString(t *testing.T) {
	tok := token{
		ttype:    cAND,
		tliteral: "AND",
	}
	tokstr := tok.String()
	if tokstr != `{Type: "AND", Literal: "AND"}` {
		t.Errorf("incorrect string for type token, got: %s, want: %s", tokstr, "Type: AND, Literal: AND")
	}
}
