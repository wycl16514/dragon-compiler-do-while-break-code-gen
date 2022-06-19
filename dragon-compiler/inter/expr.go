package inter

import (
	"lexer"
	"strconv"
)



type Expr struct {
	Node      *Node
	token     *lexer.Token
	expr_type *Type
}

func NewExpr(line uint32, token *lexer.Token, expr_type *Type) *Expr {
	expr := &Expr{
		Node:      NewNode(line),
		token:     token,
		expr_type: expr_type,
	}

	return expr
}

func (e *Expr) Errors(s string) error {
	return e.Node.Errors(s)
}

func (e *Expr) NewLabel() uint32 {
	return e.Node.NewLabel()
}

func (e *Expr) EmitLabel(i uint32) {
	e.Node.EmitLabel(i)
}

func (e *Expr) Emit(code string) {
	e.Node.Emit(code)
}

func (e *Expr) Gen() ExprInterface {
	return e
}

func (e *Expr) Reduce() ExprInterface {
	return e
}

func (e *Expr) ToString() string {
	return e.token.ToString()
}

func (e *Expr) Type() *Type {
	return e.expr_type
}

func (e *Expr) Jumping(t uint32, f uint32) {
	e.EmitJumps(e.ToString(), t, f)
}

func (e * Expr) EmitJumps(test string, t uint32, f uint32) {
	if t != 0 && f != 0 {
		e.Emit("if " + test + " got L" + strconv.Itoa(int(t)))
		e.Emit("goto L" + strconv.Itoa(int(f)))
	} else if t != 0 {
		e.Emit("if " + test + " goto L" + strconv.Itoa(int(t)))
	} else if f != 0 {
		e.Emit("iffalse " + test + " goto L" + strconv.Itoa(int(f)))
	}
}
