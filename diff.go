package diffq

import (
	"strings"

	"github.com/r3labs/diff"
)

type Diff struct {
	Changed      bool
	ChangeLog    diff.Changelog
	ChangeLogMap map[string]diff.Change
}

func Differential(a, b interface{}) (*Diff, error) {
	changes, err := diff.Diff(a, b)
	if err != nil {
		return nil, err
	}
	result := &Diff{
		ChangeLog:    changes,
		ChangeLogMap: make(map[string]diff.Change),
	}
	if len(changes) > 0 {
		result.Changed = true
	}
	for _, c := range changes {
		ident := strings.Join(c.Path, ".")
		result.ChangeLogMap[ident] = c
	}
	return result, nil
}

func (d *Diff) EvaluateStatement(statement string) (bool, error) {
	err := Validate(statement)
	if err != nil {
		return false, err
	}
	result := Evaluate(statement, d)
	return result, nil
}
