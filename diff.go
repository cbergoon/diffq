package diffq

import (
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/r3labs/diff"
)

// Changes and Change are abstracted to eliminate tight-coupling with underlying
// diff library.

// Change represents a single change identified by the differential. Change is
// intentionally abstracted from the r3labs/diff library to avoid tight coupling
// to the dependancy.
type Change struct {
	// Type indicates the change type: create, delete or update.
	Type string
	// Path is an array of field names representing the path to the field in
	// question from the outer struct.
	Path []string
	// From contains the original value of the field.
	From interface{}
	// To contains the new value of the field.
	To interface{}
}

// Changes represents a list of changes identified by the differential process.
type Changes []Change

// Diff represents the differential of two arbitrary structs and is used to
// evaluate statements against it.
type Diff struct {
	// Changed indicates the presence of any changes between two objects.
	Changed bool
	// Changes holds the complete list of changes identified by the diff
	// process.
	Changes Changes
	// ChangeLogMap maps the field identifiers of the changed values to the
	// change.
	ChangeLogMap map[string]Change
	// Original holds the original struct value
	Original interface{}
	// New holds the new struct value
	New interface{}
}

// Differential calculates the differential of a and b returning an initialized
// Diff and an error if encountered.
func Differential(a, b interface{}) (*Diff, error) {
	// calculate diff using r3labs/diff
	changes, err := diff.Diff(a, b)
	if err != nil {
		return nil, err
	}

	result := &Diff{
		ChangeLogMap: make(map[string]Change),
		Original:     a,
		New:          b,
	}

	// map changes to internal change type and build lookup map
	for _, c := range changes {
		nc := Change{
			Type: c.Type,
			Path: c.Path,
			To:   c.To,
			From: c.From,
		}
		result.Changes = append(result.Changes, nc)
		ident := strings.Join(c.Path, ".")
		result.ChangeLogMap[ident] = nc
	}

	if len(changes) > 0 {
		result.Changed = true
	}

	return result, nil
}

// HumanDifferential calculates the differential between Original and New fields
// on Diff, d, and returns a human readable string representing the diff.
func (d *Diff) HumanDifferential() string {
	diff := cmp.Diff(d.Original, d.New)
	return diff
}

// EvaluateStatement executes statement provided against Diff, d, and returns
// the validity of the statement relative to the calculated diff and any errors
// encountered.
func (d *Diff) EvaluateStatement(statement string) (bool, error) {
	err := validate(statement)
	if err != nil {
		return false, err
	}
	result, err := evaluate(statement, d)
	if err != nil {
		return false, errors.Wrap(err, "error: failed to evaluate")
	}
	return result, nil
}
