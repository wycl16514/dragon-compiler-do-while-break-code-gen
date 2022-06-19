package inter

import (
	"errors"
	"strconv"
)

type Else struct {
	stmt  *Stmt
	expr  ExprInterface
	stmt1 StmtInterface
	stmt2 StmtInterface
}

func NewElse(line uint32, expr ExprInterface, stmt1 StmtInterface, stmt2 StmtInterface) *Else {
	if expr.Type().Lexeme != "bool" {
		err := errors.New("bool type required in if")
		panic(err)
	}
	return &Else{
		stmt:  NewStmt(line),
		expr:  expr,
		stmt1: stmt1,
		stmt2: stmt2,
	}
}

func (e *Else) Errors(str string) error {
	return e.stmt.Errors(str)
}

func (e *Else) NewLabel() uint32 {
	return e.stmt.NewLabel()
}

func (e *Else) EmitLabel(i uint32) {
	e.stmt.EmitLabel(i)
}

func (e *Else) Emit(code string) {
	e.stmt.Emit(code)
}

func (e *Else) Gen(_ uint32, end uint32) {
	label1 := e.NewLabel()
	label2 := e.NewLabel()

	e.expr.Jumping(0, label2)
	e.EmitLabel(label1) //生成if条件判断中代码
	e.stmt1.Gen(label1, end) //生成if成立后大括号里面代码的中间代码
	e.Emit("goto L" + strconv.Itoa(int(end))) //增加goto语句跳过else部分代码
	e.EmitLabel(label2) 
	e.stmt2.Gen(label2, end) //生成else里面代码对应中间代码
}
