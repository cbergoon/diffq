package diffq

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/r3labs/diff"
)

// Changes and Change are abstracted to eliminate tight-coupling with underlying diff library.

// Changes represents a list of changes identified by the diff.
type Changes []Change

// Change represents a single change identified by the diff.
type Change struct {
	Type string
	Path []string
	From interface{}
	To   interface{}
}

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

func (d *Diff) EvaluateStatement(statement string) (bool, error) {
	err := Validate(statement)
	if err != nil {
		return false, err
	}
	result := Evaluate(statement, d)
	return result, nil
}

func (d *Diff) GetStructSliceFieldLenByName(selector string, v interface{}) (int, error) {
	components := strings.Split(selector, ".")
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Invalid {
		return 0, errors.New("invalid type encountered for initial reflection")
	}
	for _, c := range components {
		if r.Kind() == reflect.Ptr {
			r = reflect.Indirect(r).FieldByName(c)
		} else if r.Kind() == reflect.Slice || r.Kind() == reflect.Array {
			i, _ := strconv.ParseInt(c, 10, 64)
			r = r.Index(int(i))
		} else if r.Kind() == reflect.Map {
			r = r.MapIndex(reflect.ValueOf(c))
		} else {
			tmp := r.FieldByName(c)
			r = tmp
		}
	}
	if r.Kind() == reflect.Slice || r.Kind() == reflect.Array {
		return r.Len(), nil
	}
	return 0, errors.New("terminal type from selector is not indexible")
}

func (d *Diff) GetStructFieldByName(selector string, v interface{}) (interface{}, error) {
	components := strings.Split(selector, ".")
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Invalid {
		return 0, errors.New("invalid type encountered for initial reflection")
	}
	for _, c := range components {
		if r.Kind() == reflect.Ptr {
			r = reflect.Indirect(r).FieldByName(c)
		} else if r.Kind() == reflect.Slice || r.Kind() == reflect.Array {
			i, _ := strconv.ParseInt(c, 10, 64)
			r = r.Index(int(i))
		} else if r.Kind() == reflect.Map {
			r = r.MapIndex(reflect.ValueOf(c))
		} else {
			tmp := r.FieldByName(c)
			r = tmp
		}
	}
	return r.Interface(), nil
}
