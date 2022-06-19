package inter

import (
	"errors"
)
//if 继承 stmt
type If struct {
	stmt    StmtInterface
	expr    ExprInterface
	if_stmt StmtInterface
}

func NewIf(line uint32, expr ExprInterface, if_stmt StmtInterface) *If {
	if expr.Type().Lexeme != "bool" {
		err := errors.New("bool type required in if")
		panic(err)
	}
	return &If{
		stmt:    NewStmt(line),
		expr:    expr,
		if_stmt: if_stmt,
	}
}

func (i *If) Errors(str string) error {
	return i.stmt.Errors(str)
}

func (i *If) NewLabel() uint32 {
	return i.stmt.NewLabel()
}

func (i *If) EmitLabel(label uint32) {
	i.stmt.EmitLabel(label)
}

func (i *If) Emit(code string) {
	i.stmt.Emit(code)
}

func (i *If) Gen(_ uint32, end uint32) {
	label := i.NewLabel()
	i.expr.Jumping(0, end)
	i.EmitLabel(label)
	i.if_stmt.Gen(label, end)
}
