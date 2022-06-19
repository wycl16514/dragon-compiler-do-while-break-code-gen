package inter

import (
	"lexer"
)

type Or struct {
	logic *Logic        //它负责处理||， &&， !，等操作符一些共同的逻辑
	expr1 ExprInterface // "||"前面的表达式
	expr2 ExprInterface // "||"后面的表达式
}

func NewOr(line uint32, token *lexer.Token,
	expr1 ExprInterface, expr2 ExprInterface) *Or {
	return &Or{
		logic: NewLogic(line, token, expr1, expr2, logicCheckType),
		expr1: expr1,
		expr2: expr2,
	}
}

func (o *Or) Errors(s string) error {
	return o.logic.Errors(s)
}

func (o *Or) NewLabel() uint32 {
	return o.logic.NewLabel()
}

func (o *Or) EmitLabel(l uint32) {
	o.logic.EmitLabel(l)
}

func (o *Or) Emit(code string) {
	o.logic.Emit(code)
}

func (o *Or) Gen() ExprInterface {
	return o.logic.Gen()
}

func (o *Or) Reduce() ExprInterface {
	return o
}

func (o *Or) Type() *Type {
	return o.logic.Type()
}

func (o *Or) ToString() string {
	return o.logic.ToString()
}

func (o *Or) Jumping(t uint32, f uint32) {
	var label uint32
	if t != 0 {
		label = t
	} else {
		label = o.NewLabel()
	}

	o.expr1.Jumping(label, 0)
	o.expr2.Jumping(t, f)
	if t == 0 {
		o.EmitLabel(label)
	}
}

func (o *Or) EmitJumps(test string, t uint32, l uint32) {
	o.logic.EmitJumps(test, t, l)
}
