package diffq

import (
	"errors"
	"fmt"
	"strconv"
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

	// check goesto
	if operator.Type == GOESTO {
		x, ok := d.ChangeLogMap[identifier.Literal]
		if !ok {
			return false
		}
		// check that literal value of literal (converted to value type) == to value
		if literal.Type == INT {
			i, err := strconv.ParseInt(literal.Literal, 10, 64)
			if err != nil {

			}
			if i == cast.ToInt64(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == FLOAT {
			i, err := strconv.ParseFloat(literal.Literal, 64)
			if err != nil {

			}
			if i == cast.ToFloat64(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == STRING {
			s := literal.Literal
			if s == cast.ToString(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == DURATION {
			d, err := time.ParseDuration(literal.Literal)
			if err != nil {

			}
			if d == cast.ToDuration(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == TIME {
			t, err := time.Parse(time.RFC3339, literal.Literal)
			if err != nil {
				fmt.Println(err)
			}
			if t.Equal(cast.ToTime(x.To)) {
				return true
			} else {
				return false
			}
		} else if literal.Type == TRUE {
			bv := true
			if bv == x.To {
				return true
			} else {
				return false
			}
		} else if literal.Type == FALSE {
			bv := false
			if bv == x.To {
				return true
			} else {
				return false
			}
		}
	} else if operator.Type == NOTGOESTO {
		x, ok := d.ChangeLogMap[identifier.Literal]
		if !ok {
			return true
		}
		// check that literal value of literal (converted to value type) != to value
		if literal.Type == INT {
			i, err := strconv.ParseInt(literal.Literal, 10, 64)
			if err != nil {

			}
			if i != cast.ToInt64(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == FLOAT {
			i, err := strconv.ParseFloat(literal.Literal, 64)
			if err != nil {

			}
			if i != cast.ToFloat64(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == STRING {
			s := literal.Literal
			if s != cast.ToString(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == DURATION {
			d, err := time.ParseDuration(literal.Literal)
			if err != nil {

			}
			if d != cast.ToDuration(x.To) {
				return true
			} else {
				return false
			}
		} else if literal.Type == TIME {
			t, err := time.Parse(time.RFC3339, literal.Literal)
			if err != nil {

			}
			if !t.Equal(cast.ToTime(x.To)) {
				return true
			} else {
				return false
			}
		} else if literal.Type == TRUE {
			bv := true
			if bv != x.To {
				return true
			} else {
				return false
			}
		} else if literal.Type == FALSE {
			bv := false
			if bv != x.To {
				return true
			} else {
				return false
			}
		}
	}

	return false
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
		fmt.Println(token)
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
		// fmt.Println(token)
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
			//debug
			// for !curexpts.IsEmpty() {
			// 	cc := curexpts.Pop()
			// 	fmt.Println(cc.Literal, cc.Type)
			// }
			op := ts.Pop()

			if op.Type == EVAL {
				// TODO (cbergoon): return error?
				err := ValidateTransformStack(curexpts)
				if err != nil {
					// return err
				}
				expres := EvaluateTransformStack(curexpts, d)
				// fmt.Println("eval = ", expres)
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
			// fmt.Printf("Eval Set: [%v] %v\n", op, curexp)

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
