package diffq

import (
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/r3labs/diff"
)

// Changes and Change are abstracted to eliminate tight-coupling with underlying diff library.

// Change represents a single change identified by the diff.
type Change struct {
	Type string
	Path []string
	From interface{}
	To   interface{}
}

// Changes represents a list of changes identified by the diff.
type Changes []Change

type Diff struct {
	Changed      bool
	Changes      Changes
	ChangeLogMap map[string]Change
	A            interface{}
	B            interface{}
}

func Differential(a, b interface{}) (*Diff, error) {
	changes, err := diff.Diff(a, b)
	if err != nil {
		return nil, err
	}
	result := &Diff{
		ChangeLogMap: make(map[string]Change),
		A:            a,
		B:            b,
	}

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

func HumanDifferential(a, b interface{}) string {
	diff := cmp.Diff(a, b)
	return diff
}

func (d *Diff) EvaluateStatement(statement string) (bool, error) {
	err := Validate(statement)
	if err != nil {
		return false, err
	}
	result, err := Evaluate(statement, d)
	if err != nil {
		return false, errors.Wrap(err, "error: failed to evaluate")
	}
	return result, nil
}
