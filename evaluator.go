package diffq

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type TokenStack struct {
	Stack []*Token
}

func (s *TokenStack) Push(token *Token) {
	s.Stack = append(s.Stack, token)
}

func (s *TokenStack) Pop() *Token {
	if s.IsEmpty() {
		return nil
	}
	last := s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]
	return last
}

func (s *TokenStack) Peek() *Token {
	if s.IsEmpty() {
		return nil
	}
	last := s.Stack[len(s.Stack)-1]
	return last
}

func (s *TokenStack) IsEmpty() bool {
	if len(s.Stack) <= 0 {
		return true
	}
	return false
}

func (s *TokenStack) Size() int {
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

func ValidateTransformStack(stack *TokenStack) error {
	if stack.Size() != 3 {
		return errors.New("error: invalid number of arguments in eval")
	}
	if stack.Stack[2].Type != IDENT {
		return fmt.Errorf("error: expected identifier got %s", stack.Stack[2].Literal)
	}
	if stack.Stack[1].Type != GOESTO && stack.Stack[1].Type != NOTGOESTO {
		return fmt.Errorf("error: expected operator got %s", stack.Stack[1].Literal)
	}
	if stack.Stack[0].Type != STRING && stack.Stack[0].Type != INT && stack.Stack[0].Type != ASTERISK && stack.Stack[0].Type != DURATION && stack.Stack[0].Type != TIME {
		return fmt.Errorf("error: expected literal got %s", stack.Stack[0].Literal)
	}
	return nil
}

func EvaluateTransformStack(stack *TokenStack, d *Diff) bool {
	var identifier, operator, literal *Token
	// TODO (cbergoon): what if its not equal to 3
	if stack.Size() == 3 {
		identifier = stack.Pop()
		operator = stack.Pop()
		literal = stack.Pop()
	}

	//Rewrite
	exprChangeIdentifier := identifier

	identifierParts := strings.Split(exprChangeIdentifier.Literal, ".")
	for i := 1; i <= len(identifierParts); i++ {
		cumulativeParts := strings.Join(identifierParts[:i], ".")
		field, _ := d.GetStructFieldByName(cumulativeParts, d.B)
		fieldKind := reflect.ValueOf(field).Kind()
		if fieldKind == reflect.Array || fieldKind == reflect.Slice {
			length, _ := d.GetStructSliceFieldLenByName(cumulativeParts, d.B)
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

	// fmt.Println("rewrite summary", expandedPath, matchedChanges)

	foundValidChange := false
	if len(matchedChanges) == 0 && operator.Type == NOTGOESTO {
		foundValidChange = true
	} else {
		for _, mc := range matchedChanges {

			if operator.Type == GOESTO {
				// check that literal value of literal (converted to value type) == to value
				if literal.Type == INT {
					i, err := strconv.ParseInt(literal.Literal, 10, 64)
					if err != nil {

					}
					if i == cast.ToInt64(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == FLOAT {
					i, err := strconv.ParseFloat(literal.Literal, 64)
					if err != nil {

					}
					if i == cast.ToFloat64(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == STRING {
					s := literal.Literal
					if s == cast.ToString(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == DURATION {
					d, err := time.ParseDuration(literal.Literal)
					if err != nil {

					}
					if d == cast.ToDuration(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == TIME {
					t, err := time.Parse(time.RFC3339, literal.Literal)
					if err != nil {
						fmt.Println(err)
					}
					if t.Equal(cast.ToTime(mc.To)) {
						foundValidChange = true
					}
				} else if literal.Type == TRUE {
					bv := true
					if bv == mc.To {
						foundValidChange = true
					}
				} else if literal.Type == FALSE {
					bv := false
					if bv == mc.To {
						foundValidChange = true
					}
				} else if literal.Type == NIL {
					if mc.To == nil {
						foundValidChange = true
					}
				} else if literal.Type == ASTERISK {
					// Change matches; so change went to some value therefore true
					foundValidChange = true
				} else if literal.Type == CREATED {
					if mc.Type == "create" {
						foundValidChange = true
					}
				} else if literal.Type == DELETED {
					if mc.Type == "delete" {
						foundValidChange = true
					}
				}
			} else if operator.Type == NOTGOESTO {

				// check that literal value of literal (converted to value type) != to value
				if literal.Type == INT {
					i, err := strconv.ParseInt(literal.Literal, 10, 64)
					if err != nil {

					}
					if i != cast.ToInt64(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == FLOAT {
					i, err := strconv.ParseFloat(literal.Literal, 64)
					if err != nil {

					}
					if i != cast.ToFloat64(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == STRING {
					s := literal.Literal
					if s != cast.ToString(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == DURATION {
					d, err := time.ParseDuration(literal.Literal)
					if err != nil {

					}
					if d != cast.ToDuration(mc.To) {
						foundValidChange = true
					}
				} else if literal.Type == TIME {
					t, err := time.Parse(time.RFC3339, literal.Literal)
					if err != nil {

					}
					if !t.Equal(cast.ToTime(mc.To)) {
						foundValidChange = true
					}
				} else if literal.Type == TRUE {
					bv := true
					if bv != mc.To {
						foundValidChange = true
					}
				} else if literal.Type == FALSE {
					bv := false
					if bv != mc.To {
						foundValidChange = true
					}
				} else if literal.Type == NIL {
					if mc.To != nil {
						foundValidChange = true
					}
				} else if literal.Type == ASTERISK {
					notFound := true
					for _, ch := range matchedChanges {
						if wildcardPathMatch(expandedPath, ch.Path) {
							notFound = false
						}
					}
					if notFound {
						foundValidChange = notFound
					}
				} else if literal.Type == CREATED {
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
				} else if literal.Type == DELETED {
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

			if foundValidChange {
				return foundValidChange
			}

		}
	}

	return foundValidChange
}

func Validate(statement string) error {
	// TODO (cbergoon): ensure no illegal tokens
	// TODO (cbergoon): ensure len tokens > 0

	ts := &TokenStack{}

	lexer := NewLexer(statement)

	isBalanced := true
	isOpComplete := true

	token := lexer.NextToken()
	for isBalanced && token.Type != EOF {
		switch token.Type {
		case AND:
			fallthrough
		case OR:
			fallthrough
		case EVAL:
			fallthrough
		case LPAREN:
			ts.Push(token)
		case RPAREN:
			if ts.IsEmpty() {
				isBalanced = false
			} else {
				expectedOpenDelimeter := ts.Pop()
				if expectedOpenDelimeter.Type != LPAREN {
					isBalanced = false
				}
				expectedOpDelimeter := ts.Pop()
				if expectedOpDelimeter.Type != AND && expectedOpDelimeter.Type != OR && expectedOpDelimeter.Type != EVAL {
					isOpComplete = false
				}
			}
		}
		token = lexer.NextToken()
	}

	if !ts.IsEmpty() {
		isBalanced = false
	}

	if !isOpComplete {
		return errors.New("validation error: missing operation")
	}
	if !isBalanced {
		return errors.New("validation error: mismatched parentheses")
	}
	return nil
}

func Evaluate(statement string, d *Diff) bool {
	ts := &TokenStack{}

	lexer := NewLexer(statement)

	token := lexer.NextToken()
	for token.Type != EOF {
		if token.Type != RPAREN {
			ts.Push(token)
		} else {
			curexpts := &TokenStack{}
			for !ts.IsEmpty() {
				ct := ts.Pop()
				if ct.Type != COMMA && ct.Type != LPAREN && ct.Type != RPAREN {
					curexpts.Push(ct)
				}
				if ct.Type == LPAREN {
					break
				}
			}

			op := ts.Pop()

			if op.Type == EVAL {
				// TODO (cbergoon): return error?
				err := ValidateTransformStack(curexpts)
				if err != nil {
					// return err
				}
				expres := EvaluateTransformStack(curexpts, d)
				if expres == true {
					ts.Push(&Token{Type: TRUE, Literal: TRUE})
				} else {
					ts.Push(&Token{Type: FALSE, Literal: FALSE})
				}
			} else {
				seenTrue := false
				seenFalse := false
				for !curexpts.IsEmpty() {
					t := curexpts.Pop()
					if t.Type == FALSE {
						seenFalse = true
					} else {
						seenTrue = true
					}
				}
				if op.Type == OR {
					if seenTrue {
						ts.Push(&Token{Type: TRUE, Literal: TRUE})
					} else {
						ts.Push(&Token{Type: FALSE, Literal: FALSE})
					}
				} else if op.Type == AND {
					if !seenFalse {
						ts.Push(&Token{Type: TRUE, Literal: TRUE})
					} else {
						ts.Push(&Token{Type: FALSE, Literal: FALSE})
					}
				}
			}

		}
		token = lexer.NextToken()
	}

	result := ts.Pop().Literal
	if result == TRUE {
		return true
	} else if result == FALSE {
		return false
	} else {
		return false
	}
}
