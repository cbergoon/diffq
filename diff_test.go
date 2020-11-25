package diffq

import (
	"fmt"
	"log"
	"time"
)

type NestedType struct {
	NS  string
	NI  int
	NSS []string
}

type OuterType struct {
	S   string
	I   int
	I64 int64
	I32 int32
	B   bool
	F64 float64
	F32 float32
	SS  []string
	IS  []int

	T time.Time
	D time.Duration

	NT  NestedType
	NTP *NestedType
	NTS []*NestedType

	M map[string]int
}

func ExampleDifferential() {

	OT1 := &OuterType{
		S:   "StringS",
		I:   1,
		I64: 64,
		I32: 32,
		B:   false,
		F64: 3.1415,
		F32: 2.718,
		T:   time.Now(),
		D:   time.Duration(time.Hour),
		SS:  []string{"SS1", "SS2", "SS3"},
		IS:  []int{1, 2, 3},
		NT: NestedType{
			NS:  "StringNS",
			NI:  123,
			NSS: []string{"NSS1", "NSS2"},
		},
		NTP: &NestedType{
			NS:  "StringNS",
			NI:  123,
			NSS: []string{"NSS1", "NSS2"},
		},
		NTS: []*NestedType{
			&NestedType{
				NS:  "AStringNS",
				NI:  123,
				NSS: []string{"ANSS1", "ANSS2"},
			},
			&NestedType{
				NS:  "BStringNS",
				NI:  123,
				NSS: []string{"BNSS1", "BNSS2"},
			},
		},
		M: make(map[string]int),
	}

	newTime, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00-04:00")

	OT2 := &OuterType{
		S:   "StringSU",
		I:   12,
		I64: 64,
		I32: 32,
		B:   true,
		F64: 100.5,
		F32: 2.718,
		T:   newTime,
		D:   time.Duration(time.Hour * 2),
		SS:  []string{"SS1U", "SS2UX", "SS3U", "SS4U"},
		IS:  []int{1, 2, 3},
		NT: NestedType{
			NS:  "StringNS",
			NI:  123,
			NSS: []string{"NSS1", "NSS2"},
		},
		NTP: nil,
		NTS: []*NestedType{
			&NestedType{
				NS:  "AStringNS",
				NI:  123,
				NSS: []string{"ANSS1u", "ans", "ANSS2"},
			},
			&NestedType{
				NS:  "BStringNS",
				NI:  123,
				NSS: []string{"BNSS1", "BNSS2"},
			},
		},
		M: make(map[string]int),
	}

	OT1.M["one"] = 1
	OT1.M["two"] = 2

	OT2.M["one"] = 2
	OT2.M["two"] = 3

	var d *Diff
	d, _ = Differential(OT1, OT2)

	examples := []string{
		`AND(
			EVAL(S ["StringS"] => "StringSU"),
			EVAL(T => t"2020-01-01T12:00:00-04:00"),
			EVAL(D => d"2h"),
			OR(
				EVAL(F64 => 100.5),
				EVAL(NTS.0.NSS.0 => "ANSS1u")
			),
			EVAL(B => true),
			EVAL(NTP => nil),
			EVAL(I32 =!> *),
			EVAL(SS.$first => "SS1U"),
			EVAL(SS.* => *),
			EVAL(M.one => 2),
			EVAL(SS.$last => "SS4U"),
			EVAL(SS.* => $created),
			EVAL(SS.1 => "SS2UX"), /* when the second element is SS2UX */
			EVAL(F64 =LTE> 110.0)
		)`,
		`AND(
			EVAL(SS.* ["SS2"] => "SS2UX"),
		)`,
	}

	for i := 0; i < 1; i++ {
		for _, ex := range examples {

			result, err := d.EvaluateStatement(ex)
			if err != nil {
				log.Fatalf("error: failed to evaluate statement: %v", err)
			}

			fmt.Println(result)
		}
	}

	// Output:
	// true
	// true

}
