package lisp

import "fmt"

var Env map[string]interface{}

func init() {
	Env = make(map[string]interface{})
}

func Execute(line string) (string, error) {
	tokenized := Tokenize(line)
	parsed, err := Parse(tokenized)
	if err != nil {
		return "", err
	}
	evaled, err := Eval(parsed)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", evaled), nil
}

func Eval(expr interface{}) (val interface{}, err error) {
	switch expr.(type) {
	case int: // Int
		val = expr
	case string: // Symbol
		sym := expr.(string)
		if v, ok := Env[sym]; ok {
			val = v
		} else if sym == "true" || sym == "false" {
			val = sym
		} else {
			err = fmt.Errorf("Unknown symbol: %v", expr)
		}
	case Sexp:
		tokens := expr.(Sexp)
		t := tokens[0]
		if _, ok := t.(Sexp); ok {
			val, err = Eval(t)
		} else if t == "quote" { // Quote
			val = tokens[1:]
		} else if t == "define" { // Define
			if len(tokens) > 2 {
				val, err = Eval(tokens[2])
			}
			if err == nil && len(tokens) > 1 {
				Env[tokens[1].(string)] = val
			}
		} else if t == "set!" { // Set!
			key := tokens[1].(string)
			if _, ok := Env[key]; ok {
				val, err = Eval(tokens[2])
				if err == nil {
					Env[key] = val
				}
			} else {
				err = fmt.Errorf("Can only set! variable that is previously defined")
			}
		} else if t == "if" { // If
			r, err := Eval(tokens[1])
			if err == nil {
				if r != "false" && len(tokens) > 2 {
					val, err = Eval(tokens[2])
				} else if len(tokens) > 3 {
					val, err = Eval(tokens[3])
				}
			}
		} else if t == "begin" { // Begin
			for _, val = range tokens[1:] {
				val, err = Eval(val)
				if err != nil {
					break
				}
			}
		} else if t == "+" { // Addition
			var sum int
			for _, i := range tokens[1:] {
				j, err := Eval(i)
				if err == nil {
					v, ok := j.(int)
					if ok {
						sum += int(v)
					} else {
						err = fmt.Errorf("Cannot only add numbers: %v", i)
						break
					}
				}
			}
			val = sum
		} else {
			return Eval(t)
		}
	default:
		err = fmt.Errorf("Unknown data type: %v", expr)
	}
	return
}
