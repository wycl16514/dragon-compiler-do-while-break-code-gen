package inter

import (
	"lexer"
)

type Not struct {
	logic *Logic
	expr1 *Expr
	expr2 *Expr
}

func NewNot(line uint32, token *lexer.Token,
	expr1 *Expr, expr2 *Expr) *Not {
	return &Not{
		logic: NewLogic(line, token, expr1, expr2, logicCheckType),
		expr1: expr1,
		expr2: expr2,
	}
}

func (n *Not) Errors(s string) error {
	return n.logic.Errors(s)
}

func (n *Not) NewLabel() uint32 {
	return n.logic.NewLabel()
}

func (n *Not) EmitLabel(l uint32) {
	n.logic.EmitLabel(l)
}

func (n *Not) Emit(code string) {
	n.logic.Emit(code)
}

func (n *Not) Gen() ExprInterface {
	return n.logic.Gen()
}

func (n *Not) Reduce() ExprInterface {
	return n
}

func (n *Not) Type() *Type {
	return n.logic.Type()
}

func (n *Not) ToString() string {
	return n.logic.ToString()
}

func (n *Not) Jumping(t uint32, f uint32) {
	n.expr2.Jumping(f, t) //not 正好要反过来
}

func (n *Not) EmitJumps(test string, t uint32, l uint32) {
	n.logic.EmitJumps(test, t, l)
}
