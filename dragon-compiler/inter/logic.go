package inter

import (
	"errors"
	"lexer"
	"strconv"
)



/*
实现or, and , !等操作
*/

type Logic struct {
	expr      ExprInterface
	token     *lexer.Token
	expr1     ExprInterface
	expr2     ExprInterface
	expr_type *Type
	line      uint32
}

type CheckType func(type1 *Type, type2 *Type) *Type

func logicCheckType(type1 *Type, type2 *Type) *Type {

	if type1.Lexeme == "bool" && type2.Lexeme == "bool" {
		return type1
	}


	return nil
}

func NewLogic(line uint32, token *lexer.Token,
	expr1 ExprInterface, expr2 ExprInterface, checkType CheckType) *Logic {
	expr_type := checkType(expr1.Type(), expr2.Type())
	if expr_type == nil {
		err := errors.New("type error")
		panic(err)
	}

	return &Logic{
		expr:      NewExpr(line, token, expr_type),
		token:     token,
		expr1:     expr1,
		expr2:     expr2,
		expr_type: expr_type,
		line:      line,
	}
}

func (l *Logic) Errors(s string) error {
	return l.expr.Errors(s)
}

func (l *Logic) NewLabel() uint32 {
	return l.expr.NewLabel()
}

func (l *Logic) EmitLabel(label uint32) {
	l.expr.EmitLabel(label)
}

func (l *Logic) Emit(code string) {
	l.expr.Emit(code)
}

func (l *Logic) Type() *Type {
	return l.expr_type
}

func (l *Logic) Gen() ExprInterface {
	f := l.NewLabel()
	a := l.NewLabel()
	temp := NewTemp(l.line, l.expr_type)
	l.Jumping(0, f)
	l.Emit(temp.ToString() + " = true")
	l.Emit("goto L" + strconv.Itoa(int(a)))
	l.EmitLabel(f)
	l.Emit(temp.ToString() + "=false")
	l.EmitLabel(a)
	return temp
}

func (l *Logic) Reduce() ExprInterface {
	return l
}

func (l *Logic) ToString() string {
	return l.expr1.ToString() + " " + l.token.ToString() + " " + l.expr2.ToString()
}

func (l *Logic) Jumping(t uint32, f uint32) {
	l.expr.Jumping(t, f)
}

func (l *Logic) EmitJumps(test string, t uint32, f uint32) {
	l.expr.EmitJumps(test, t, f)
}
