package inter

import (
	"errors"
	"lexer"
)

type Unary struct {
	/*
		它对应-a, -b, !a等表达式
	*/
	op *Op
	x  *Expr
}

func NewUnary(line uint32, token *lexer.Token, expr_type *Type, x *Expr) *Unary {
	unary := &Unary{
		op: NewOp(line, token, expr_type),
		x:  x,
	}

	if x.Type().Lexeme != "int" {
		err := errors.New("type error")
		panic(err)
	}

	return unary
}

func (u *Unary) Errors(s string) error {
	return u.op.expr.Errors(s)
}

func (u *Unary) NewLabel() uint32 {
	return u.op.expr.NewLabel()
}

func (u *Unary) EmitLabel(l uint32) {
	u.op.expr.EmitLabel(l)
}

func (u *Unary) Emit(code string) {
	u.op.expr.Emit(code)
}

func (u *Unary) Gen() ExprInterface {
	return u
}

func (u *Unary) Reduce() ExprInterface {
	return u.x.Reduce()
}

func (u *Unary) Type() *Type {
	return u.op.expr.Type()
}

func (u *Unary) ToString() string {
	return u.op.ToString() + " " + u.x.ToString()
}

func (u *Unary) Jumping(t uint32, f uint32) {
	u.op.expr.Jumping(t, f)
}

func (u *Unary) EmitJumps(test string, t uint32, f uint32) {
	u.op.expr.EmitJumps(test, t, f)
}
