package diffq

import (
	"testing"
)

func TestLexerTokenization(t *testing.T) {
	testInputs := []string{
		`/* test comment */
		AND( /* test comment */
			EVAL(S ["StringS"] => "StringSU"),
			EVAL(T => t"2020-01-01T12:00:00-04:00"),
			EVAL(D => d"2h"),
			OR(
				/* test comment */
				EVAL(F64 => 100.5),
				EVAL(NTS.0.NSS.0 => "ANSS1u"), /* test comment */
			),
			EVAL(B => true),
			EVAL(NTP => nil),
			EVAL(I32 =!> *),
			EVAL(F64 =LTE> -110.0),
			EVAL(F64 =LT> 110.0),
			EVAL(F64 =GTE> 110.5),
			EVAL(F64 =GT> 110.0),
			EVAL(M.one => 2),
			EVAL(SS.* => *),
			EVAL(SS.*.K => *),
			EVAL(SS.$first => "SS1U"),
			EVAL(SS.$last => "SS4U"), /* test comment */
			EVAL(SS.* => $created), 
			EVAL(SS.* => $deleted), 
			EVAL(SS.1 => "SS2UX")
		)
		/* test comment */`,
	}

	testTokens := [][]*token{
		[]*token{
			{ttype: "COMMENT", tliteral: " test comment "},
			{ttype: "AND", tliteral: "AND"},
			{ttype: "(", tliteral: "("},
			{ttype: "COMMENT", tliteral: " test comment "},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "S"},
			{ttype: "[", tliteral: "["},
			{ttype: "STRING", tliteral: "StringS"},
			{ttype: "]", tliteral: "]"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "STRING", tliteral: "StringSU"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "T"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "TIME", tliteral: "2020-01-01T12:00:00-04:00"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "D"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "DURATION", tliteral: "2h"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "OR", tliteral: "OR"},
			{ttype: "(", tliteral: "("},
			{ttype: "COMMENT", tliteral: " test comment "},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "F64"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "FLOAT", tliteral: "100.5"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "NTS.0.NSS.0"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "STRING", tliteral: "ANSS1u"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "COMMENT", tliteral: " test comment "},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "B"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "TRUE", tliteral: "true"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "NTP"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "NIL", tliteral: "nil"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "I32"},
			{ttype: "=!>", tliteral: "=!>"},
			{ttype: "*", tliteral: "*"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "F64"},
			{ttype: "=LTE>", tliteral: "=LTE>"},
			{ttype: "FLOAT", tliteral: "-110.0"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "F64"},
			{ttype: "=LT>", tliteral: "=LT>"},
			{ttype: "FLOAT", tliteral: "110.0"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "F64"},
			{ttype: "=GTE>", tliteral: "=GTE>"},
			{ttype: "FLOAT", tliteral: "110.5"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "F64"},
			{ttype: "=GT>", tliteral: "=GT>"},
			{ttype: "FLOAT", tliteral: "110.0"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "M.one"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "INT", tliteral: "2"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "SS.*"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "*", tliteral: "*"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "SS.*.K"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "*", tliteral: "*"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "SS.$first"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "STRING", tliteral: "SS1U"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "SS.$last"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "STRING", tliteral: "SS4U"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "COMMENT", tliteral: " test comment "},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "SS.*"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "$created", tliteral: "$created"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "SS.*"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "$deleted", tliteral: "$deleted"},
			{ttype: ")", tliteral: ")"},
			{ttype: ",", tliteral: ","},
			{ttype: "EVAL", tliteral: "EVAL"},
			{ttype: "(", tliteral: "("},
			{ttype: "IDENT", tliteral: "SS.1"},
			{ttype: "=>", tliteral: "=>"},
			{ttype: "STRING", tliteral: "SS2UX"},
			{ttype: ")", tliteral: ")"},
			{ttype: ")", tliteral: ")"},
			{ttype: "COMMENT", tliteral: " test comment "},
		},
	}

	for inputSetIndex, inputSetInput := range testInputs {
		tokens := []*token{}
		lex := newLexer(inputSetInput)
		tok := lex.nextToken()
		for tok.ttype != EOF {
			tokens = append(tokens, tok)
			tok = lex.nextToken()
		}
		if len(tokens) != len(testTokens[inputSetIndex]) {
			t.Errorf("incorrect length of resulting token list, got: %d, want: %d", len(tokens), len(testTokens))
		}
		for tokIndex, gotTok := range tokens {
			if gotTok.ttype != testTokens[inputSetIndex][tokIndex].ttype {
				t.Errorf("incorrect token in input set %d at token %d, got: %s, want: %s", inputSetIndex, tokIndex, gotTok, testTokens[inputSetIndex][tokIndex])
			}
			if gotTok.tliteral != testTokens[inputSetIndex][tokIndex].tliteral {
				t.Errorf("incorrect token in input set %d at token %d, got: %s, want: %s", inputSetIndex, tokIndex, gotTok, testTokens[inputSetIndex][tokIndex])
			}
		}
	}
}
