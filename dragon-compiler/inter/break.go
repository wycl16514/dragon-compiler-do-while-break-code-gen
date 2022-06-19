package inter

import (
	"fmt"
)

type Break struct {
	stmt      StmtInterface //父节点
	enclosing StmtInterface //包裹break语句的循环体对象
}

func NewBreak(line uint32, enclosing StmtInterface) *Break {

	if _, ok := enclosing.(*While); !ok {
		_, ok = enclosing.(*Do)
		if !ok {
			panic("unenclosed break") //break语句没有处于循环体中
		}
	}

	return &Break{
		stmt:      NewStmt(line),
		enclosing: enclosing,
	}
}

//下面仅仅是调用其父类接口
func (b *Break) Errors(str string) error {
	return b.stmt.Errors(str)
}

func (b *Break) NewLabel() uint32 {
	return b.stmt.NewLabel()
}

func (b *Break) EmitLabel(label uint32) {
	b.stmt.EmitLabel(label)
}

func (b *Break) Emit(code string) {
	b.stmt.Emit(code)
}

func (b *Break) Gen(_ uint32, _ uint32) {
	while_loop, while_ok := b.enclosing.(*While)
	do_loop, _ := b.enclosing.(*Do)
	var code string
	if while_ok {
		code = fmt.Sprintf("goto L%d", while_loop.GetAfter())
	} else {
		code = fmt.Sprintf("goto L%d", do_loop.GetAfter())
	}

	b.Emit(code)
}
