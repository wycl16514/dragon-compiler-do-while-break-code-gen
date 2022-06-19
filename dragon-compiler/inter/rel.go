package inter

import (
	"lexer"
)

type Rel struct {
	logic *Logic
	expr1 ExprInterface
	expr2 ExprInterface
	token *lexer.Token
}


func relCheckType(p1 *Type, p2 *Type) *Type {
	if p1.Lexeme == p2.Lexeme {
		return NewType("bool", lexer.BASIC, 1)
	}

	return nil
}

func NewRel(line uint32, token *lexer.Token,
	expr1 ExprInterface, expr2 ExprInterface) *Rel {
	return &Rel{
		logic: NewLogic(line, token, expr1, expr2, relCheckType),
		expr1: expr1,
		expr2: expr2,
		token: token,
	}
}

func (r *Rel) Errors(s string) error {
	return r.logic.Errors(s)
}

func (r *Rel) NewLabel() uint32 {
	return r.logic.NewLabel()
}

func (r *Rel) EmitLabel(l uint32) {
	r.logic.EmitLabel(l)
}

func (r *Rel) Emit(code string) {
	r.logic.Emit(code)
}

func (r *Rel) Gen() ExprInterface {
	return r.logic.Gen()
}

func (r *Rel) Reduce() ExprInterface {
	return r
}

func (r *Rel) Type() *Type {
	return r.logic.Type()
}

func (r *Rel) ToString() string {
	return r.logic.ToString()
}



func (r *Rel) Jumping(t uint32, f uint32) {
	expr1 := r.expr1.Reduce()
	expr2 := r.expr2.Reduce()
	test := expr1.ToString() + " " + r.token.ToString() + " " + expr2.ToString()
	r.EmitJumps(test, t, f)
}

func (r *Rel) EmitJumps(test string, t uint32, l uint32) {
	r.logic.EmitJumps(test, t, l)
}
