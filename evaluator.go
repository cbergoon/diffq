package diffq

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// tokenStack reprents a stack of tokens. tokenStack is used by the evaluation
// operations as the evaluation stack.
type tokenStack struct {
	Stack []*token
}

// push appends a token to the bottom of the stack.
func (s *tokenStack) push(token *token) {
	s.Stack = append(s.Stack, token)
}

// pop returns and removes a token from the top of the stack.
func (s *tokenStack) pop() *token {
	if s.isEmpty() {
		return nil
	}
	last := s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]
	return last
}

// peek returns the token from the top of the stack but does not remove it.
func (s *tokenStack) peek() *token {
	if s.isEmpty() {
		return nil
	}
	last := s.Stack[len(s.Stack)-1]
	return last
}

// isEmpty return true if the stack is empty; false otherwise.
func (s *tokenStack) isEmpty() bool {
	if len(s.Stack) <= 0 {
		return true
	}
	return false
}

// size returns the length of the stack.
func (s *tokenStack) size() int {
	return len(s.Stack)
}

// wildcardPathMatch matches the composed identifier "filter" with wildcards to
// the components of the path. Example: A.B.* matches A.B.C, A.B.D, etc.; A.*.C
// matches A.B.C, A.D.C, etc.
func wildcardPathMatch(filter, path []string) bool {
	for i, f := range filter {
		if len(path) < i+1 {
			return false
		}
		if f != path[i] && f != "*" {
			return false
		}
	}

	return true
}

// getStructSliceFieldLenByName reflects on the provided value 'v' to determine
// the length of the field provided identified by 'selector'. It is assumed that
// the selector identifies an array. The function returns the length of the
// identified value.
func (d *Diff) getStructSliceFieldLenByName(selector string, v interface{}) (int, error) {
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

// getStructFieldByName returns the value of the field represented by 'selector'
// in the value 'v' provided.
func (d *Diff) getStructFieldByName(selector string, v interface{}) (interface{}, error) {
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

// validateTransformStack ensures that the transform stack is valid for
// operation. The transform stack represents the actual operations/comparisons
// to be executed (the portions of the statement that is contained in EVAL
// statements). This validation step enforces the structure of EVAL statements
// contents ensuring that the components of the expression are of the correct
// type and semantically correct.
func validateTransformStack(stack *tokenStack) error {
	if stack.size() == 3 { // Standard expression - assuming previous value not present expecting 3 components
		if stack.Stack[2].ttype != cIDENT {
			return errors.Errorf("validation error: expected identifier got %s", stack.Stack[2].tliteral)
		}
		if stack.Stack[1].ttype != cGOESTO && stack.Stack[1].ttype != cNOTGOESTO && stack.Stack[1].ttype != cGOESGT && stack.Stack[1].ttype != cGOESGTE && stack.Stack[1].ttype != cGOESLT && stack.Stack[1].ttype != cGOESLTE {
			return errors.Errorf("validation error: expected operator got %s", stack.Stack[1].tliteral)
		}
		if stack.Stack[0].ttype != cSTRING && stack.Stack[0].ttype != cINT && stack.Stack[0].ttype != cFLOAT && stack.Stack[0].ttype != cASTERISK && stack.Stack[0].ttype != cDURATION && stack.Stack[0].ttype != cTIME && stack.Stack[0].ttype != cTRUE && stack.Stack[0].ttype != cFALSE && stack.Stack[0].ttype != cNIL && stack.Stack[0].ttype != cCREATED && stack.Stack[0].ttype != cDELETED {
			return errors.Errorf("validation error: expected literal got %s", stack.Stack[0].tliteral)
		}
		// If operator is comparison literal cannot be 'nil' or '*'
		if stack.Stack[1].ttype == cGOESGT || stack.Stack[1].ttype == cGOESGTE || stack.Stack[1].ttype == cGOESLT || stack.Stack[1].ttype == cGOESLTE {
			if stack.Stack[0].ttype == cASTERISK || stack.Stack[0].ttype == cNIL || stack.Stack[0].ttype == cCREATED || stack.Stack[0].ttype == cDELETED {
				return errors.New("validation error: cannot use literal values '*' or 'nil' with comparison operators")
			}
		}
	} else if stack.size() == 4 { // Previous value - assuming previous value present expecting 4 components
		if stack.Stack[3].ttype != cIDENT {
			return errors.Errorf("validation error: expected identifier got %s", stack.Stack[3].tliteral)
		}
		if stack.Stack[2].ttype != cSTRING && stack.Stack[2].ttype != cINT && stack.Stack[2].ttype != cFLOAT && stack.Stack[2].ttype != cASTERISK && stack.Stack[2].ttype != cDURATION && stack.Stack[2].ttype != cTIME && stack.Stack[2].ttype != cTRUE && stack.Stack[2].ttype != cFALSE && stack.Stack[2].ttype != cNIL {
			return errors.Errorf("validation error: expected literal got %s", stack.Stack[2].tliteral)
		}
		if stack.Stack[1].ttype != cGOESTO && stack.Stack[1].ttype != cNOTGOESTO && stack.Stack[1].ttype != cGOESGT && stack.Stack[1].ttype != cGOESGTE && stack.Stack[1].ttype != cGOESLT && stack.Stack[1].ttype != cGOESLTE {
			return errors.Errorf("validation error: expected operator got %s", stack.Stack[1].tliteral)
		}
		if stack.Stack[0].ttype != cSTRING && stack.Stack[0].ttype != cINT && stack.Stack[0].ttype != cFLOAT && stack.Stack[0].ttype != cASTERISK && stack.Stack[0].ttype != cDURATION && stack.Stack[0].ttype != cTIME && stack.Stack[0].ttype != cTRUE && stack.Stack[0].ttype != cFALSE && stack.Stack[0].ttype != cNIL {
			// If length of stack is 4 then assume using previous value; cannot
			// use deleted or created with previous value
			if stack.Stack[0].ttype == cCREATED || stack.Stack[0].ttype == cDELETED {
				return errors.New("validation error: cannot specify action literal of $created or $deleted when using previous value")
			}
			return errors.Errorf("validation error: expected literal got %s", stack.Stack[0].tliteral)
		}
		// If operator is comparison literal cannot be 'nil' or '*'
		if stack.Stack[1].ttype == cGOESGT || stack.Stack[1].ttype == cGOESGTE || stack.Stack[1].ttype == cGOESLT || stack.Stack[1].ttype == cGOESLTE {
			if stack.Stack[0].ttype == cASTERISK || stack.Stack[0].ttype == cNIL {
				return errors.New("validation error: cannot use literal values '*' or 'nil' with comparison operators")
			}
		}
	} else {
		return errors.New("validation error: invalid number of arguments in eval")
	}

	return nil
}

// evaluateTransformStack evaluates the transform stack which represents the
// actual comparison operations inside EVAL expressions. This function returns
// the validity of the expression provided in the transform stack as either true
// or false.
func evaluateTransformStack(stack *tokenStack, d *Diff) bool {
	var identifier, previous, operator, literal *token

	if stack.size() == 3 {
		identifier = stack.pop()
		operator = stack.pop()
		literal = stack.pop()
	} else if stack.size() == 4 {
		identifier = stack.pop()
		previous = stack.pop()
		operator = stack.pop()
		literal = stack.pop()
	} else {
		return false // TODO (cbergoon): error here?
	}

	// rewrite/expand expression
	exprChangeIdentifier := identifier

	identifierParts := strings.Split(exprChangeIdentifier.tliteral, ".")
	for i := 1; i <= len(identifierParts); i++ {
		cumulativeParts := strings.Join(identifierParts[:i], ".")
		field, _ := d.getStructFieldByName(cumulativeParts, d.New)
		fieldKind := reflect.ValueOf(field).Kind()
		if fieldKind == reflect.Array || fieldKind == reflect.Slice {
			length, _ := d.getStructSliceFieldLenByName(cumulativeParts, d.New)
			if length > 0 {
				if len(identifierParts) > i {
					if identifierParts[i] == "$first" {
						identifierParts[i] = "0"
						i++
					} else if identifierParts[i] == "$last" {
						identifierParts[i] = fmt.Sprint(length - 1)
						i++
					}
				}
			}
		}
	}

	expandedPath := identifierParts

	var matchedChanges Changes
	for _, c := range d.Changes {
		if wildcardPathMatch(expandedPath, c.Path) {
			matchedChanges = append(matchedChanges, c)
		}
	}

	// TODO (cbergoon): handle errors below?
	foundValidChange := false
	if len(matchedChanges) == 0 && operator.ttype == cNOTGOESTO {
		foundValidChange = true
	} else {
		for _, mc := range matchedChanges {

			previousConditionValid := true
			if previous != nil {
				previousConditionValid = false
				if previous.ttype == cINT {
					i, err := strconv.ParseInt(previous.tliteral, 10, 64)
					if err != nil {

					}
					if i == cast.ToInt64(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == cFLOAT {
					i, err := strconv.ParseFloat(previous.tliteral, 64)
					if err != nil {

					}
					if i == cast.ToFloat64(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == cSTRING {
					s := previous.tliteral
					if s == cast.ToString(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == cDURATION {
					d, err := time.ParseDuration(previous.tliteral)
					if err != nil {

					}
					if d == cast.ToDuration(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == cTIME {
					t, err := time.Parse(time.RFC3339, previous.tliteral)
					if err != nil {
					}
					if t.Equal(cast.ToTime(mc.From)) {
						previousConditionValid = true
					}
				} else if previous.ttype == cTRUE {
					bv := true
					if bv == mc.From {
						previousConditionValid = true
					}
				} else if previous.ttype == cFALSE {
					bv := false
					if bv == mc.From {
						previousConditionValid = true
					}
				} else if previous.ttype == cNIL {
					if mc.From == nil {
						previousConditionValid = true
					}
				} else if previous.ttype == cASTERISK {
					previousConditionValid = true
				}
			}

			if previousConditionValid {
				if operator.ttype == cGOESTO {
					if literal.ttype == cINT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if i == cast.ToInt64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cFLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if i == cast.ToFloat64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cSTRING {
						s := literal.tliteral
						if s == cast.ToString(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cDURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if d == cast.ToDuration(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cTIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {

						}
						if t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					} else if literal.ttype == cTRUE {
						bv := true
						if bv == mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == cFALSE {
						bv := false
						if bv == mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == cNIL {
						if mc.To == nil {
							foundValidChange = true
						}
					} else if literal.ttype == cASTERISK {
						// Change matches; so change went to some value
						// therefore true
						foundValidChange = true
					} else if literal.ttype == cCREATED {
						if mc.Type == "create" {
							foundValidChange = true
						}
					} else if literal.ttype == cDELETED {
						if mc.Type == "delete" {
							foundValidChange = true
						}
					}
				} else if operator.ttype == cGOESGT {
					if literal.ttype == cINT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) > i {
							foundValidChange = true
						}
					} else if literal.ttype == cFLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) > i {
							foundValidChange = true
						}
					} else if literal.ttype == cSTRING {
						s := literal.tliteral
						if cast.ToString(mc.To) > s {
							foundValidChange = true
						}
					} else if literal.ttype == cDURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) > d {
							foundValidChange = true
						}
					} else if literal.ttype == cTIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {

						}
						if cast.ToTime(mc.To).After(t) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == cGOESGTE {
					if literal.ttype == cINT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) >= i {
							foundValidChange = true
						}
					} else if literal.ttype == cFLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) >= i {
							foundValidChange = true
						}
					} else if literal.ttype == cSTRING {
						s := literal.tliteral
						if cast.ToString(mc.To) >= s {
							foundValidChange = true
						}
					} else if literal.ttype == cDURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) >= d {
							foundValidChange = true
						}
					} else if literal.ttype == cTIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {

						}
						if cast.ToTime(mc.To).After(t) || t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == cGOESLT {
					if literal.ttype == cINT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) < i {
							foundValidChange = true
						}
					} else if literal.ttype == cFLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) < i {
							foundValidChange = true
						}
					} else if literal.ttype == cSTRING {
						s := literal.tliteral
						if cast.ToString(mc.To) < s {
							foundValidChange = true
						}
					} else if literal.ttype == cDURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) < d {
							foundValidChange = true
						}
					} else if literal.ttype == cTIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {

						}
						if cast.ToTime(mc.To).Before(t) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == cGOESLTE {
					if literal.ttype == cINT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) <= i {
							foundValidChange = true
						}
					} else if literal.ttype == cFLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) <= i {
							foundValidChange = true
						}
					} else if literal.ttype == cSTRING {
						s := literal.tliteral
						if cast.ToString(mc.To) <= s {
							foundValidChange = true
						}
					} else if literal.ttype == cDURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) <= d {
							foundValidChange = true
						}
					} else if literal.ttype == cTIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {

						}
						if cast.ToTime(mc.To).Before(t) || t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == cNOTGOESTO {
					if literal.ttype == cINT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if i != cast.ToInt64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cFLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if i != cast.ToFloat64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cSTRING {
						s := literal.tliteral
						if s != cast.ToString(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cDURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if d != cast.ToDuration(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == cTIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {

						}
						if !t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					} else if literal.ttype == cTRUE {
						bv := true
						if bv != mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == cFALSE {
						bv := false
						if bv != mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == cNIL {
						if mc.To != nil {
							foundValidChange = true
						}
					} else if literal.ttype == cASTERISK {
						notFound := true
						for _, ch := range matchedChanges {
							if wildcardPathMatch(expandedPath, ch.Path) {
								notFound = false
							}
						}
						if notFound {
							foundValidChange = notFound
						}
					} else if literal.ttype == cCREATED {
						notFound := true
						for _, ch := range matchedChanges {
							if ch.Type == "create" {
								if wildcardPathMatch(expandedPath, ch.Path) {
									notFound = false
								}
							}
						}
						if notFound {
							foundValidChange = notFound
						}
					} else if literal.ttype == cDELETED {
						notFound := true
						for _, ch := range matchedChanges {
							if ch.Type == "delete" {
								if wildcardPathMatch(expandedPath, ch.Path) {
									notFound = false
								}
							}
						}
						if notFound {
							foundValidChange = notFound
						}
					}
				}
			} else {
				foundValidChange = false
			}

			if foundValidChange {
				return foundValidChange
			}

		}
	}

	return foundValidChange
}

// validate ensures the structure and semantics of the complete statement.
// Validate returns an error if the statement is not valid or is otherwise
// incorrect.
func validate(statement string) error {
	if statement == "" {
		return errors.Errorf("validation error: empty statement")
	}

	ts := &tokenStack{}

	lexer := newLexer(statement)

	isBalanced := true
	isOpComplete := true
	tokenCount := 0

	token := lexer.nextToken()
	for isBalanced && token.ttype != cEOF {
		tokenCount++
		switch token.ttype {
		case cILLEGAL:
			return errors.Errorf("validation error: illegal token %s", token.tliteral)
		case cAND:
			fallthrough
		case cOR:
			fallthrough
		case cEVAL:
			fallthrough
		case cLPAREN:
			ts.push(token)
		case cRPAREN:
			if ts.isEmpty() {
				isBalanced = false
			} else {
				expectedOpenDelimeter := ts.pop()
				if expectedOpenDelimeter.ttype != cLPAREN {
					isBalanced = false
				}
				expectedOpDelimeter := ts.pop()
				if expectedOpDelimeter.ttype != cAND && expectedOpDelimeter.ttype != cOR && expectedOpDelimeter.ttype != cEVAL {
					isOpComplete = false
				}
			}
		}
		token = lexer.nextToken()
	}

	if !ts.isEmpty() {
		isBalanced = false
	}

	if !isOpComplete {
		return errors.New("validation error: missing operation")
	}
	if !isBalanced {
		return errors.New("validation error: mismatched parentheses")
	}
	if tokenCount == 0 {
		return errors.New("validation error: token length zero for provided statement")
	}
	return nil
}

// evaluate executes the the 'statement' against the Diff 'd' provided. Returns
// the boolean result of the statement and an error if encountered. Evaluate
// manages the entire execution and handles validation as well as the execution
// of the individual transform stacks within the statement.
func evaluate(statement string, d *Diff) (bool, error) {
	ts := &tokenStack{}

	lexer := newLexer(statement)

	tok := lexer.nextToken()
	for tok.ttype != cEOF {

		if tok.ttype != cRPAREN { // continue populating stack until hit right paren
			ts.push(tok)
		} else { // if right paren encountered then begin execution of the component until the most recent (previous) left paren.
			// populate EVAL transfor stack
			curexpts := &tokenStack{}
			for !ts.isEmpty() {
				ct := ts.pop()
				if ct.ttype != cCOMMA && ct.ttype != cLPAREN && ct.ttype != cRPAREN && ct.ttype != cLBRACKET && ct.ttype != cRBRACKET {
					curexpts.push(ct)
				}
				if ct.ttype == cLPAREN {
					break
				}
			}

			op := ts.pop()

			if op.ttype == cEVAL {
				// if operator is EVAL then validate and execute pushing result
				// onto the stack
				err := validateTransformStack(curexpts)
				if err != nil {
					return false, errors.Wrap(err, "error: invalid transform stack; failed validation")
				}
				expres := evaluateTransformStack(curexpts, d)
				if expres == true {
					ts.push(&token{ttype: cTRUE, tliteral: cTRUE})
				} else {
					ts.push(&token{ttype: cFALSE, tliteral: cFALSE})
				}
			} else {
				// if operator is not EVAL (AND or OR) then evaluate the boolean
				// expression
				seenTrue := false
				seenFalse := false
				for !curexpts.isEmpty() {
					t := curexpts.pop()
					if t.ttype == cFALSE {
						seenFalse = true
					} else {
						seenTrue = true
					}
				}
				if op.ttype == cOR {
					if seenTrue {
						ts.push(&token{ttype: cTRUE, tliteral: cTRUE})
					} else {
						ts.push(&token{ttype: cFALSE, tliteral: cFALSE})
					}
				} else if op.ttype == cAND {
					if !seenFalse {
						ts.push(&token{ttype: cTRUE, tliteral: cTRUE})
					} else {
						ts.push(&token{ttype: cFALSE, tliteral: cFALSE})
					}
				}
			}

		}
		tok = lexer.nextToken()
	}

	// last result represents the result of the statement
	result := ts.pop().tliteral
	if result == cTRUE {
		return true, nil
	} else if result == cFALSE {
		return false, nil
	} else {
		return false, nil
	}
}
