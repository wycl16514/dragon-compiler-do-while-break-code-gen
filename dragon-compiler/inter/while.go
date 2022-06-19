package inter

import (
	"errors"
	"fmt"
)



type While struct {
	stmt       *Stmt         //继承自Stmt节点
	expr       ExprInterface //对应while 后面的条件判断表达式
	while_stmt StmtInterface //对应while的循环体部分
	after      uint32
}

//为了确保break语句能与给定的while节点对应，我们需要进行修改
func NewWhile(line uint32, expr ExprInterface) *While {
	if expr.Type().Lexeme != "bool" {
		//用于while后面的表达式必须为bool类型
		err := errors.New("bool type required for while")
		panic(err)
	}

	return &While{
		stmt:       NewStmt(line),
		expr:       expr,
		while_stmt: nil,
		after:      0,
	}
}

func (w *While) InitWhileBody(while_stmt StmtInterface) {
	w.while_stmt = while_stmt
}


//下面仅仅是调用其父类接口
func (w *While) Errors(str string) error {
	return w.stmt.Errors(str)
}

func (w *While) NewLabel() uint32 {
	return w.stmt.NewLabel()
}

func (w *While) EmitLabel(label uint32) {
	w.stmt.EmitLabel(label)
}

func (w *While) Emit(code string) {
	w.stmt.Emit(code)
}

func (w *While) setAfter(after uint32) {
	w.after = after 
}

func (w *While) GetAfter() uint32 {
	return w.after 
}

func (w *While) Gen(start uint32, end uint32) {
	w.expr.Jumping(0, end)
    w.setAfter(end) //记录下循环体外第一句代码的标号
	label := w.NewLabel()
	w.EmitLabel(label)
	w.while_stmt.Gen(label, start) //生成while循环体语句的起始标志
	emit_code := fmt.Sprintf("goto L%d", start)
	w.Emit(emit_code)
}
