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

type tokenStack struct {
	Stack []*token
}

func (s *tokenStack) push(token *token) {
	s.Stack = append(s.Stack, token)
}

func (s *tokenStack) pop() *token {
	if s.isEmpty() {
		return nil
	}
	last := s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]
	return last
}

func (s *tokenStack) peek() *token {
	if s.isEmpty() {
		return nil
	}
	last := s.Stack[len(s.Stack)-1]
	return last
}

func (s *tokenStack) isEmpty() bool {
	if len(s.Stack) <= 0 {
		return true
	}
	return false
}

func (s *tokenStack) size() int {
	return len(s.Stack)
}

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

func validateTransformStack(stack *tokenStack) error {
	if stack.size() == 3 {
		if stack.Stack[2].ttype != IDENT {
			return errors.Errorf("validation error: expected identifier got %s", stack.Stack[2].tliteral)
		}
		if stack.Stack[1].ttype != GOESTO && stack.Stack[1].ttype != NOTGOESTO && stack.Stack[1].ttype != GOESGT && stack.Stack[1].ttype != GOESGTE && stack.Stack[1].ttype != GOESLT && stack.Stack[1].ttype != GOESLTE {
			return errors.Errorf("validation error: expected operator got %s", stack.Stack[1].tliteral)
		}
		if stack.Stack[0].ttype != STRING && stack.Stack[0].ttype != INT && stack.Stack[0].ttype != FLOAT && stack.Stack[0].ttype != ASTERISK && stack.Stack[0].ttype != DURATION && stack.Stack[0].ttype != TIME && stack.Stack[0].ttype != TRUE && stack.Stack[0].ttype != FALSE && stack.Stack[0].ttype != NIL && stack.Stack[0].ttype != CREATED && stack.Stack[0].ttype != DELETED {
			return errors.Errorf("validation error: expected literal got %s", stack.Stack[0].tliteral)
		}
		// If operator is comparison literal cannot be 'nil' or '*'
		if stack.Stack[1].ttype == GOESGT || stack.Stack[1].ttype == GOESGTE || stack.Stack[1].ttype == GOESLT || stack.Stack[1].ttype == GOESLTE {
			if stack.Stack[0].ttype == ASTERISK || stack.Stack[0].ttype == NIL || stack.Stack[0].ttype == CREATED || stack.Stack[0].ttype == DELETED {
				return errors.New("validation error: cannot use literal values '*' or 'nil' with comparison operators")
			}
		}
	} else if stack.size() == 4 {
		if stack.Stack[3].ttype != IDENT {
			return errors.Errorf("validation error: expected identifier got %s", stack.Stack[3].tliteral)
		}
		if stack.Stack[2].ttype != STRING && stack.Stack[2].ttype != INT && stack.Stack[2].ttype != FLOAT && stack.Stack[2].ttype != ASTERISK && stack.Stack[2].ttype != DURATION && stack.Stack[2].ttype != TIME && stack.Stack[2].ttype != TRUE && stack.Stack[2].ttype != FALSE && stack.Stack[2].ttype != NIL {
			return errors.Errorf("validation error: expected literal got %s", stack.Stack[2].tliteral)
		}
		if stack.Stack[1].ttype != GOESTO && stack.Stack[1].ttype != NOTGOESTO && stack.Stack[1].ttype != GOESGT && stack.Stack[1].ttype != GOESGTE && stack.Stack[1].ttype != GOESLT && stack.Stack[1].ttype != GOESLTE {
			return errors.Errorf("validation error: expected operator got %s", stack.Stack[1].tliteral)
		}
		if stack.Stack[0].ttype != STRING && stack.Stack[0].ttype != INT && stack.Stack[0].ttype != FLOAT && stack.Stack[0].ttype != ASTERISK && stack.Stack[0].ttype != DURATION && stack.Stack[0].ttype != TIME && stack.Stack[0].ttype != TRUE && stack.Stack[0].ttype != FALSE && stack.Stack[0].ttype != NIL {
			// If length of stack is 4 then assume using previous value; cannot use deleted or created with previous value
			if stack.Stack[0].ttype == CREATED || stack.Stack[0].ttype == DELETED {
				return errors.New("validation error: cannot specify action literal of $created or $deleted when using previous value")
			}
			return errors.Errorf("validation error: expected literal got %s", stack.Stack[0].tliteral)
		}
		// If operator is comparison literal cannot be 'nil' or '*'
		if stack.Stack[1].ttype == GOESGT || stack.Stack[1].ttype == GOESGTE || stack.Stack[1].ttype == GOESLT || stack.Stack[1].ttype == GOESLTE {
			if stack.Stack[0].ttype == ASTERISK || stack.Stack[0].ttype == NIL {
				return errors.New("validation error: cannot use literal values '*' or 'nil' with comparison operators")
			}
		}
	} else {
		return errors.New("validation error: invalid number of arguments in eval")
	}

	return nil
}

func evaluateTransformStack(stack *tokenStack, d *Diff) bool {
	var identifier, previous, operator, literal *token
	// TODO (cbergoon): what if its not equal to 3
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

	//Rewrite
	exprChangeIdentifier := identifier

	identifierParts := strings.Split(exprChangeIdentifier.tliteral, ".")
	for i := 1; i <= len(identifierParts); i++ {
		cumulativeParts := strings.Join(identifierParts[:i], ".")
		field, _ := d.getStructFieldByName(cumulativeParts, d.B)
		fieldKind := reflect.ValueOf(field).Kind()
		if fieldKind == reflect.Array || fieldKind == reflect.Slice {
			length, _ := d.getStructSliceFieldLenByName(cumulativeParts, d.B)
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

	foundValidChange := false
	if len(matchedChanges) == 0 && operator.ttype == NOTGOESTO {
		foundValidChange = true
	} else {
		for _, mc := range matchedChanges {

			previousConditionValid := true
			if previous != nil {
				previousConditionValid = false
				if previous.ttype == INT {
					i, err := strconv.ParseInt(previous.tliteral, 10, 64)
					if err != nil {

					}
					if i == cast.ToInt64(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == FLOAT {
					i, err := strconv.ParseFloat(previous.tliteral, 64)
					if err != nil {

					}
					if i == cast.ToFloat64(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == STRING {
					s := previous.tliteral
					if s == cast.ToString(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == DURATION {
					d, err := time.ParseDuration(previous.tliteral)
					if err != nil {

					}
					if d == cast.ToDuration(mc.From) {
						previousConditionValid = true
					}
				} else if previous.ttype == TIME {
					t, err := time.Parse(time.RFC3339, previous.tliteral)
					if err != nil {
						fmt.Println(err)
					}
					if t.Equal(cast.ToTime(mc.From)) {
						previousConditionValid = true
					}
				} else if previous.ttype == TRUE {
					bv := true
					if bv == mc.From {
						previousConditionValid = true
					}
				} else if previous.ttype == FALSE {
					bv := false
					if bv == mc.From {
						previousConditionValid = true
					}
				} else if previous.ttype == NIL {
					if mc.From == nil {
						previousConditionValid = true
					}
				} else if previous.ttype == ASTERISK {
					previousConditionValid = true
				}
			}

			if previousConditionValid {
				if operator.ttype == GOESTO {
					if literal.ttype == INT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if i == cast.ToInt64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == FLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if i == cast.ToFloat64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == STRING {
						s := literal.tliteral
						if s == cast.ToString(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == DURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if d == cast.ToDuration(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == TIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {
							fmt.Println(err)
						}
						if t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					} else if literal.ttype == TRUE {
						bv := true
						if bv == mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == FALSE {
						bv := false
						if bv == mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == NIL {
						if mc.To == nil {
							foundValidChange = true
						}
					} else if literal.ttype == ASTERISK {
						// Change matches; so change went to some value therefore true
						foundValidChange = true
					} else if literal.ttype == CREATED {
						if mc.Type == "create" {
							foundValidChange = true
						}
					} else if literal.ttype == DELETED {
						if mc.Type == "delete" {
							foundValidChange = true
						}
					}
				} else if operator.ttype == GOESGT {
					if literal.ttype == INT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) > i {
							foundValidChange = true
						}
					} else if literal.ttype == FLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) > i {
							foundValidChange = true
						}
					} else if literal.ttype == STRING {
						s := literal.tliteral
						if cast.ToString(mc.To) > s {
							foundValidChange = true
						}
					} else if literal.ttype == DURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) > d {
							foundValidChange = true
						}
					} else if literal.ttype == TIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {
							fmt.Println(err)
						}
						if cast.ToTime(mc.To).After(t) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == GOESGTE {
					if literal.ttype == INT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) >= i {
							foundValidChange = true
						}
					} else if literal.ttype == FLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) >= i {
							foundValidChange = true
						}
					} else if literal.ttype == STRING {
						s := literal.tliteral
						if cast.ToString(mc.To) >= s {
							foundValidChange = true
						}
					} else if literal.ttype == DURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) >= d {
							foundValidChange = true
						}
					} else if literal.ttype == TIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {
							fmt.Println(err)
						}
						if cast.ToTime(mc.To).After(t) || t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == GOESLT {
					if literal.ttype == INT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) < i {
							foundValidChange = true
						}
					} else if literal.ttype == FLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) < i {
							foundValidChange = true
						}
					} else if literal.ttype == STRING {
						s := literal.tliteral
						if cast.ToString(mc.To) < s {
							foundValidChange = true
						}
					} else if literal.ttype == DURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) < d {
							foundValidChange = true
						}
					} else if literal.ttype == TIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {
							fmt.Println(err)
						}
						if cast.ToTime(mc.To).Before(t) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == GOESLTE {
					if literal.ttype == INT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if cast.ToInt64(mc.To) <= i {
							foundValidChange = true
						}
					} else if literal.ttype == FLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if cast.ToFloat64(mc.To) <= i {
							foundValidChange = true
						}
					} else if literal.ttype == STRING {
						s := literal.tliteral
						if cast.ToString(mc.To) <= s {
							foundValidChange = true
						}
					} else if literal.ttype == DURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if cast.ToDuration(mc.To) <= d {
							foundValidChange = true
						}
					} else if literal.ttype == TIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {
							fmt.Println(err)
						}
						if cast.ToTime(mc.To).Before(t) || t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					}
				} else if operator.ttype == NOTGOESTO {
					if literal.ttype == INT {
						i, err := strconv.ParseInt(literal.tliteral, 10, 64)
						if err != nil {

						}
						if i != cast.ToInt64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == FLOAT {
						i, err := strconv.ParseFloat(literal.tliteral, 64)
						if err != nil {

						}
						if i != cast.ToFloat64(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == STRING {
						s := literal.tliteral
						if s != cast.ToString(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == DURATION {
						d, err := time.ParseDuration(literal.tliteral)
						if err != nil {

						}
						if d != cast.ToDuration(mc.To) {
							foundValidChange = true
						}
					} else if literal.ttype == TIME {
						t, err := time.Parse(time.RFC3339, literal.tliteral)
						if err != nil {

						}
						if !t.Equal(cast.ToTime(mc.To)) {
							foundValidChange = true
						}
					} else if literal.ttype == TRUE {
						bv := true
						if bv != mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == FALSE {
						bv := false
						if bv != mc.To {
							foundValidChange = true
						}
					} else if literal.ttype == NIL {
						if mc.To != nil {
							foundValidChange = true
						}
					} else if literal.ttype == ASTERISK {
						notFound := true
						for _, ch := range matchedChanges {
							if wildcardPathMatch(expandedPath, ch.Path) {
								notFound = false
							}
						}
						if notFound {
							foundValidChange = notFound
						}
					} else if literal.ttype == CREATED {
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
					} else if literal.ttype == DELETED {
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

func Validate(statement string) error {
	if statement == "" {
		return errors.Errorf("validation error: empty statement")
	}

	ts := &tokenStack{}

	lexer := newLexer(statement)

	isBalanced := true
	isOpComplete := true
	tokenCount := 0

	token := lexer.nextToken()
	for isBalanced && token.ttype != EOF {
		tokenCount++
		switch token.ttype {
		case ILLEGAL:
			return errors.Errorf("validation error: illegal token %s", token.tliteral)
		case AND:
			fallthrough
		case OR:
			fallthrough
		case EVAL:
			fallthrough
		case LPAREN:
			ts.push(token)
		case RPAREN:
			if ts.isEmpty() {
				isBalanced = false
			} else {
				expectedOpenDelimeter := ts.pop()
				if expectedOpenDelimeter.ttype != LPAREN {
					isBalanced = false
				}
				expectedOpDelimeter := ts.pop()
				if expectedOpDelimeter.ttype != AND && expectedOpDelimeter.ttype != OR && expectedOpDelimeter.ttype != EVAL {
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

func Evaluate(statement string, d *Diff) (bool, error) {
	ts := &tokenStack{}

	lexer := newLexer(statement)

	tok := lexer.nextToken()
	for tok.ttype != EOF {
		if tok.ttype != RPAREN {
			ts.push(tok)
		} else {
			curexpts := &tokenStack{}
			for !ts.isEmpty() {
				ct := ts.pop()
				if ct.ttype != COMMA && ct.ttype != LPAREN && ct.ttype != RPAREN && ct.ttype != LBRACKET && ct.ttype != RBRACKET {
					curexpts.push(ct)
				}
				if ct.ttype == LPAREN {
					break
				}
			}

			op := ts.pop()

			if op.ttype == EVAL {
				err := validateTransformStack(curexpts)
				if err != nil {
					return false, errors.Wrap(err, "error: invalid transform stack; failed validation")
				}
				expres := evaluateTransformStack(curexpts, d)
				if expres == true {
					ts.push(&token{ttype: TRUE, tliteral: TRUE})
				} else {
					ts.push(&token{ttype: FALSE, tliteral: FALSE})
				}
			} else {
				seenTrue := false
				seenFalse := false
				for !curexpts.isEmpty() {
					t := curexpts.pop()
					if t.ttype == FALSE {
						seenFalse = true
					} else {
						seenTrue = true
					}
				}
				if op.ttype == OR {
					if seenTrue {
						ts.push(&token{ttype: TRUE, tliteral: TRUE})
					} else {
						ts.push(&token{ttype: FALSE, tliteral: FALSE})
					}
				} else if op.ttype == AND {
					if !seenFalse {
						ts.push(&token{ttype: TRUE, tliteral: TRUE})
					} else {
						ts.push(&token{ttype: FALSE, tliteral: FALSE})
					}
				}
			}

		}
		tok = lexer.nextToken()
	}

	result := ts.pop().tliteral
	if result == TRUE {
		return true, nil
	} else if result == FALSE {
		return false, nil
	} else {
		return false, nil
	}
}
