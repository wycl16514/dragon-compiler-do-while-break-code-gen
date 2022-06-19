package inter

import (
	"lexer"
)

type And struct {
	logic *Logic
	expr1 ExprInterface
	expr2 ExprInterface
}

func NewAnd(line uint32, token *lexer.Token,
	expr1 ExprInterface, expr2 ExprInterface) *And {
	return &And{
		logic: NewLogic(line, token, expr1, expr2, logicCheckType),
		expr1: expr1,
		expr2: expr2,
	}
}

func (a *And) Errors(s string) error {
	return a.logic.Errors(s)
}

func (a *And) NewLabel() uint32 {
	return a.logic.NewLabel()
}

func (a *And) EmitLabel(l uint32) {
	a.logic.EmitLabel(l)
}

func (a *And) Emit(code string) {
	a.logic.Emit(code)
}

func (a *And) Gen() ExprInterface {
	return a.logic.Gen()
}

func (a *And) Reduce() ExprInterface {
	return a
}

func (a *And) Type() *Type {
	return a.logic.Type()
}

func (a *And) ToString() string {
	return a.logic.ToString()
}

func (a *And) Jumping(t uint32, f uint32) {
	var label uint32
	if f != 0 {
		label = f
	} else {
		label = a.NewLabel()
	}
	a.expr1.Jumping(0, label)
	a.expr2.Jumping(t, f)
	if f == 0 {
		a.EmitLabel(label)
	}
}

func (a *And) EmitJumps(test string, t uint32, l uint32) {
	a.logic.EmitJumps(test, t, l)
}
