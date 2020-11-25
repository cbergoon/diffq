<h1 align="center">DIFFQ</h1>
<p align="center">
<a href="https://travis-ci.org/cbergoon/diffq"><img src="https://travis-ci.org/cbergoon/diffq.svg?branch=master" alt="Build"></a>
<a href="https://goreportcard.com/report/github.com/cbergoon/diffq"><img src="https://goreportcard.com/badge/github.com/cbergoon/diffq?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/diffq"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Identify and query for changes and differences between two objects. 

Diffq provides a query language to identify changes using a declarative style language. It utilizes the r3labs/diff library to identify the changes and applies a query to those changes allowing changes to be identified dynamically. 

### Install

```
go get -u github.com/cbergoon/diffq
```

### Example

```go
package main 

import (
	"fmt"
	"log"
	"time"

	"github.com/cbergoon/diffq"
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

	var d *diffq.Diff
    d, _ = diffq.Differential(OT1, OT2)
    
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
			EVAL(SS.1 => "SS2UX"), 
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
}
```

### Language Survey

The diffq language consists of two major constructs: boolean expresion statements and evaluation statements which are nested in conditional statements. 

The language generally takes the form of an outer boolean operator with evaluation statements and additional boolean operators nested within. 

Below is an informal definition of the language structure: 

```
Statement   := [AND|OR]([Statement|Evaluator]+)
Evaluator   := EVAL(Identifier Previous Operator Literal)
Identifier  := [A-Z|a-z|.|-|_|*|$]+
Previous    := Literal
Operator    := [=>|=!>|=LT>|=GT>|=LTE>|=GTE>]
```

#### Example

Given a struct `S` and two arbitrary instances of the struct `a` and `b`: 

```go
type S struct {
    Step    string
    Status  string
    Value   int
    Dur     time.Duration
    Aliases []string
    PtrInt  int
}
```

Below are some valid queries demonstrating the structure of the language:

```
// This statement evaluates to true when the `Value` field does not change to 100.
AND(
    EVAL(Value =!> 100)
)
```

```
// This statement evaluates to true when the `Value` field goes greater than 100 and either or both of 1) the `Status` 
// field changes from "New" to "Scheduled" or 2) the `Step` field changes to "PROC-2".
AND(
    EVAL(Value =GT> 100),
    OR(
        EVAL(Status ["New"] => "Scheduled"), 
        EVAL(Step => "PROC-2")
    )
)
```

```
// This statement evaluates to true when the an element of the `Aliases` field is created. 
AND(
    EVAL(Aliases.* => $created)
)
```

```
// This statement evaluates to true when the `Step` field changes to any value.
AND(
    EVAL(Step => *)
)
```

#### Identifiers

Identifiers indicate the field by using a path-like syntax. Nested types are accessed by concatenating the field names with a '.' (period). Array and map indicies are accessed in the same manner but by using the appropriate key or modifier no square brackets or quotes required. 

Special modifiers are available for arrays to access the first and last elements of an array. These are `$first` and `$last` and can be interspurced in the directly in the identifier. 

```
AND(
    EVAL(Aliases.$first => "Test")
) 
```

Additionally, asterisks can be used as wildcards to match arbitrary fields or indicies in a path. 

```
AND(
    EVAL(Aliases.* => "Test")
) 
```

#### Operators

Operators indicate how the evaluator should compare the changed values. Operators can be though of as "goes to" operators (i.e. a value "goes to" 3) with modifiers. THe full list of operators are listed below: 

```
GOES TO:                        => 
DOES NOT GO TO:                 =!>
GOES GREATER THAN:              =GT>
GOES GREATER THAN OR EQUAL:     =GTE>
GOES LESS THAN:                 =LT>
GOES LESS THAN OR EQUAL:        =LTE>
```

Operators always directly follow the identifier or the previous value if present and semantically are relative to the change or new value. 

#### Literal Values

Literal values are the represent the types that can be compared to the changed values. Literal values in the diffq language are int, float, string, boolean, time and, duration. These type are represented as shown below: 

```
INT:        100, -123
FLOAT:      1.5, 100.5, -3.1415
STRING:     "diffq string"
BOOLEAN:    TRUE, FALSE, true, false
TIME:       t"2020-01-01T12:00:00-04:00"
DURATION:   d"24h"
```

There are also 4 additional types of special literal valules: asterisk, nil, $created and, $deleted. These sepcial literal values are used as show below: 

```
ASTERISK:   EVAL(Step => *) // Step changes to any value
NIL:        EVAL(PtrInt => nil) // PtrInt goes to nil
CREATED:    EVAL(Aliases.* => $created) // An element of Aliases is created
DELETED:    EVAL(Aliases.* => $deleted) // An element of Aliases is deleted
```

#### Previous 

Previous values are optional and allow the evaluator to more selectively control a match. The previous value signifies that the match must have changed from the specified value in order to be considered a match. If no previous value is provided then it is not considered when matching a rule and the previous value can be any. 

```
EVAL(Step ["PROC-1"] => "PROC-2") // Step must change from "PROC-1" to "PROC-2"
```

#### Comments

Comments use the "/* */" format and can be used with a statement.

### License

MIT - See [LICENSE](https://github.com/cbergoon/diffq/blob/master/LICENSE) file.