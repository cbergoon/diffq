package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cbergoon/diffq"
)

// DONE (cbergoon): Add Negative Numbers
// DONE (cbergoon): Add Floats
// DONE (cbergoon): Add Duration
// DONE (cbergoon): Add Time
// DONE (cbergoon): Add Booleans

// TODO (cbergoon): Array Elements?
// TODO (cbergoon): Add Ability to Capture Creates and Deletes in Arrays
// TODO (cbergoon): Add First/Last Index Shortcuts
// TODO (cbergoon): Add `*` for Any Value Literal
// TODO (cbergoon): Add `*` for Any Array Index in Identifier

// TODO (cbergoon): Goes to Nil?
// TODO (cbergoon): Add `*` for Any Value Literal

// TODO (cbergoon): Add Operators for GOESGREATERTHAN, GOESLESSTHAN, GOESTOGREATERTHANEQUAL, GOESTOLESSTHANEQUAL

// TODO (cbergoon): Add NOT Operator

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
}

func main() {

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
	}

	newTime, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00-04:00")

	OT2 := &OuterType{
		S:   "StringSU",
		I:   12,
		I64: 64,
		I32: 32,
		B:   true,
		F64: -3.14159,
		F32: 2.718,
		T:   newTime,
		D:   time.Duration(time.Hour * 2),
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
				NSS: []string{"ANSS1u", "ANSS2"},
			},
			&NestedType{
				NS:  "BStringNS",
				NI:  123,
				NSS: []string{"BNSS1", "BNSS2"},
			},
		},
	}

	d, err := diffq.Differential(OT1, OT2)
	for _, di := range d.ChangeLog {
		fmt.Printf("%s %s %v %v\n", di.Type, di.Path, di.From, di.To)
	}
	if err != nil {
		fmt.Println(err)
	}

	examples := []string{
		// `AND(EVAL(S => "StringSU"), EVAL(I => 12), EVAL(NTS.0.NSS.0 => "ANSS1u"))`,
		`AND(
			EVAL(S =!> "StringSUX"),
			EVAL(T => t"2020-01-01T12:00:00-04:00"),
			EVAL(D => d"2h"),
			OR(
				EVAL(F64 => -3.14159),
				EVAL(NTS.0.NSS.0 => "ANSS1u")
			),
			EVAL(B => true)
		)`,
	}

	for _, ex := range examples {
		fmt.Println("Running: ", ex)
		result, err := d.EvaluateStatement(ex)
		if err != nil {
			log.Fatalf("error: failed to evaluate statement: %v", err)
		}
		fmt.Println(result)
	}
}
