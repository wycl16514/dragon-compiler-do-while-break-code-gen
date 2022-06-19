package inter

import (
	"errors"
)

type Do struct {
	stmt       *Stmt         //继承自Stmt节点
	expr       ExprInterface //对应while 后面的条件判断表达式
	while_stmt StmtInterface //对应while的循环体部分
	after      uint32
}

//为了确保break语句能与给定的while节点对应，我们需要进行修改
func NewDo(line uint32) *Do {
	return &Do{
		stmt:       NewStmt(line),
		expr:       nil,
		while_stmt: nil,
		after:      0,
	}
}

func (d *Do) InitDo(expr ExprInterface, while_stmt StmtInterface) {
	if expr.Type().Lexeme != "bool" {
		//用于while后面的表达式必须为bool类型
		err := errors.New("bool type required for Do...While")
		panic(err)
	}
	d.while_stmt = while_stmt
	d.expr = expr
}

//下面仅仅是调用其父类接口
func (d *Do) Errors(str string) error {
	return d.stmt.Errors(str)
}

func (d *Do) NewLabel() uint32 {
	return d.stmt.NewLabel()
}

func (d *Do) EmitLabel(label uint32) {
	d.stmt.EmitLabel(label)
}

func (d *Do) Emit(code string) {
	d.stmt.Emit(code)
}

func (d *Do) setAfter(after uint32) {
	d.after = after
}

func (d *Do) GetAfter() uint32 {
	return d.after
}

func (d *Do) Gen(start uint32, end uint32) {
	d.setAfter(end) //记录下循环体外第一句代码的标号
	label := d.NewLabel()
	d.while_stmt.Gen(start, label) //生成while循环体语句的起始标志
	d.EmitLabel(label)
	d.expr.Jumping(start, 0)
}
