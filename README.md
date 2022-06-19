我们的简易编译器完成了一大部分，但还有一些关键的语法结构没有处理，那就是for, while, do..while等循环语句对应的中间代码还没有生成，本节我们就针对这些语法结构进行相应的中间代码生成。

首先我们要了解循环语句对应的语法表达式：
```
stmt -> "while"  "( " bool ")" stmts
stmt -> "do" stmts "while" "(" bool ")" ";"
stmt-> "break"
```
为了简单起见，我们暂时不处理for循环，有兴趣的同学可以自己添加试试。下面我们先创建while, do..while语法结构对应的语法树节，在inter文件夹中创建while.go，然后添加代码如下：
```
package inter

import (
	"errors"
	"fmt"
)

type While struct {
	stmt       *Stmt         //继承自Stmt节点
	expr       ExprInterface //对应while 后面的条件判断表达式
	while_stmt StmtInterface //对应while的循环体部分
}

func NewWhile(line uint32, expr ExprInterface, while_stmt StmtInterface) *While {
	if expr.Type().Lexeme != "bool" {
		//用于while后面的表达式必须为bool类型
		err := errors.New("bool type required for while")
		panic(err)
	}

	return &While{
		stmt:       NewStmt(line),
		expr:       expr,
		while_stmt: while_stmt,
	}
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

func (w *While) Gen(start uint32, end uint32) {
	w.expr.Jumping(0, end)
	label := w.NewLabel()
	w.EmitLabel(label)
	w.while_stmt.Gen(label, start) //生成while循环体语句的起始标志
	emit_code := fmt.Sprintf("goto L%d", start)
	w.Emit(emit_code)
}

```
上面代码中需要注意的就是Gen函数，首先它创建跳转标签，注意这些标签对循环的正确执行有着非常重要的作用，然后它先对while后面的判断表达式生成代码，然后对while循环体内的语句集合生成代码，具体的逻辑讲解请参看b站搜索Coding迪斯尼参看我的调试演示。

接下来我们要做的是修改语法解析代码，在list_parser.go中修改stmt解析函数如下：
```
func (s *SimpleParser) stmt() inter.StmtInterface {
	/*
		if "(" bool ")"
		if -> "(" bool ")" ELSE  "{" stmt "}"

		bool -> bool "||"" join | join
		join -> join "&&" equality | equality
		equality -> equality "==" rel | equality != rel | rel
		rel -> expr < expr | expr <= expr | expr >= expr | expr > expr | expr
		rel : a > b , a < b, a <= b
		a < b && c > d || e < f
	*/
	switch s.cur_tok.Tag {
	case lexer.IF:
		s.move_forward()
		err := s.matchLexeme("(")
		if err != nil {
			panic(err)
		}
		s.move_forward()
		x := s.bool()
		err = s.matchLexeme(")")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过 ）
		s.move_forward() //越过{
		s1 := s.stmt()
		err = s.matchLexeme("}")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过}

		//判断if 后面是否跟着else
		if s.cur_tok.Tag != lexer.ELSE {
			return inter.NewIf(s.lexer.Line, x, s1)
		} else {
			s.move_forward() //越过else关键字
			err = s.matchLexeme("{")
			if err != nil {
				panic(err)
			}
			s.move_forward() //越过{
			s2 := s.stmt()   //else 里面包含的代码块
			err = s.matchLexeme("}")
			if err != nil {
				panic(err)
			}
			return inter.NewElse(s.lexer.Line, x, s1, s2)
		}

	case lexer.WHILE:
		s.move_forward()
		//while 后面跟着左括号， 然后是判断表达式，以右括号结尾
		err := s.matchLexeme("(")
		if err != nil {
			panic(err)
		}
		s.move_forward()
		while_bool := s.bool()
		err = s.matchLexeme(")")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过 ）
		s.move_forward() //越过{
		//解析while循环成立时要执行的语句块
		while_stmt := s.stmts()
		err = s.matchLexeme("}")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过}
		return inter.NewWhile(s.lexer.Line, while_bool, while_stmt)

	default:
		return s.expression()
	}
}
```
这里我们增加了对while关键字的判断，然后执行其对应的语法解析逻辑，完成上面代码后，我们在main.go中实现包含while语句的代码，这样就能运行上面代码并查看结果:
```
func main() {

	source := `{int a; int b; int c; 
		        a = 3;
				b = 0;
				while (a >= 0 && b <= 4) {
					a = a - 1;
					b = b + 1;
				}

				c = 2;
				
	}`
	my_lexer := lexer.NewLexer(source)
	parser := simple_parser.NewSimpleParser(my_lexer)
	parser.Parse()
}
```
代码运行后输出结果如下：
![请添加图片描述](https://img-blog.csdnimg.cn/d8e576eeeca342e2b5bd7b6c2c62675a.png)
我们简单分析一下输出结果，从L4开始就是while循环体输出的代码，L4对应的语句就是while后面条件判断对应的中间代码，它表明如果a >= 0 ， b <= 4 这两个条件只要有一个不成立 ，那么就跳转到L5,注意到L5正好对应while循环体出去后的第一条语句，因此生成的中间代码其逻辑符合我们在main.go中给定代码的意图。如果进入L6,也就是 a>=0和b <= 4都成立，那么就进入while循环体内部，从L6, L7可以看出他们确实是while循环体内两条语句对应的中间代码，注意到L7还有一条goto L4的语句，它表明循环体执行结束后再次调到循环体开头去对条件进行判断，如果条件依然成立，那么代码继续进入L6开始的语句进行执行，要不然就直接跳转到L5，因此从输出结果看，它是满足我们给定代码逻辑的。

接着我们看看break语句的实现，break必须要出现在循环中才能成立，因此我们在遇到该语句时，需要判断其是否位于while 或者do..while循环中，一旦执行break语句时，编译器会使用goto语句跳转到循环体外面接下来的语句，例如从上面例子中，接着循环体的第一条语句是L5,因此break执行时对应的输出就是goto L5，所以要生成break语句对应的中间代码就需要记录它所在循环体外边接下来第一条语句的标号。

在实现break时还有一点要注意，那就是循环嵌套，代码可能有多个while嵌套，于是在执行break时一定要对应到给定的while上，例如：
```
while() {
    while() {
        while() {
            break; //对应最里面的while
        }
        //对应中间while
    }
    break; //对应最外层while
}
```
因此为了应对这种情况，我们在语法解析时需要使用一个栈来记录while循环的嵌套，所以我们首先在list_parser.go中做一些修改：
```
type SimpleParser struct {
	lexer          lexer.Lexer
	top            *Env
	saved          *Env
	cur_tok        lexer.Token           //当前读取到的token
	used_storage   uint32                //当前用于存储变量的内存字节数
	loop_enclosing []inter.StmtInterface //用于循环体记录
}
```
在解析到while的时候，我们要把当前生成的while节点压入loop_enclosing栈，在解析到break语句时需要从堆栈上弹出与它对应的while节点，因此在parser函数的while部分我们要做一些修改:
```
case lexer.WHILE:
		s.move_forward()
		//while 后面跟着左括号， 然后是判断表达式，以右括号结尾
		err := s.matchLexeme("(")
		if err != nil {
			panic(err)
		}
		s.move_forward()
		while_bool := s.bool()
		err = s.matchLexeme(")")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过 ）
		s.move_forward() //越过{
		//解析while循环成立时要执行的语句块
		//这里需要注意可能解析到break语句，所以要提前生成while节点
		while_node := inter.NewWhile(s.lexer.Line, while_bool)
		//将当前while节点加入栈，解析break语句时从栈顶拿到包围它的循环语句
		s.loop_enclosing = append(s.loop_enclosing, while_node)

		while_stmt := s.stmts()
		err = s.matchLexeme("}")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过}
		while_node.InitWhileBody(while_stmt)
		return while_node
```
上面代码中我们对while的初始化也做了修改，原因是在解析它的循环体语句时可能会遇到break语句，这时候我们需要确保while节点已经生成了，所以代码改成了先构造while节点，然后再调用stmts()去解析while内部语句，这样解析到break语句时它才能找到对应的while节点，下面我们看看break节点的实现，在inter目录下创建break.go，实现代码如下：
```
package inter

import (
	"fmt"
)

type Break struct {
	stmt  StmtInterface //父节点
	enclosing   StmtInterface  //包裹break语句的循环体对象
}

func NewBreak(line uint32, enclosing StmtInterface) *Break {
	if _, ok := enclosing.(*While); !ok {
		//后面增加Do循环时还需修改这里的判断
		panic("unenclosed break") //break语句没有处于循环体中
	}

	return &Break {
		stmt: NewStmt(line),
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

func (b *Break)Gen(_ uint32, _ uint32) {
	enclosing_loop, _ := b.enclosing.(*While)
    code := fmt.Sprintf("goto L%d", enclosing_loop.GetAfter())
    b.Emit(code)
}
```
它的实现没有什么特别，唯一值得关注的就是Gen函数，它从对应的while节点取得循环体出去后的第一条语句地址，然后创建goto 指令直接跳转到给定语句处。最后我们在parse函数中增加对break语句的解析：
```
case lexer.BREAK:
		s.move_forward()
		s.matchLexeme(";")
		enclosing_while := s.loop_enclosing[len(s.loop_enclosing)-1]
		s.loop_enclosing = s.loop_enclosing[0 : len(s.loop_enclosing)-1]
		s.move_forward()
		return inter.NewBreak(s.lexer.Line, enclosing_while)
```
完成上面代码后，我们把main.go里面要解析的代码修改如下：
```
source := `{int a; int b; int c; 
		        a = 3;
				b = 0;
				while (a >= 0 && b <= 4) {
					a = a - 1;
					b = b + 1;
                    if (b < 2) {
                        break;
                    } else {
                       c = c + 1;
                    }
				}

				c = 2;
				
	}`
```
我们在while 循环中加了if判断，如果条件成立则执行break语句，我们看看代码运行结果：


![请添加图片描述](https://img-blog.csdnimg.cn/483221ed2b654f6f928a13f7b1308e58.png)
我们分析一下生成的指令，现在我们的代码已经比较复杂了，我们需要关注L7开始部分，L7开始对应的是while循环体里面的if语句，如果if判断不成立就跳转到L9,而L9正好对应else部分，也就是要执行c = c+1;如果if成立那么直接进入L8, 而在if内部直接运行break语句，由于break语句要跳出循环体直接指向循环体外面接下来的第一条语句，而代码中循环体外面第一条语句所在处就是L2，因此L8接下来就是goto L2，这条指令是break语句生成。问题在于后面还接着goto L4，这是为什么？goto L4其实是else节点生成，它的作用是指向if成立部分代码后就要跳过else部分代码，goto L4是else出来后接下来的第一条指令，而这条指令恰巧又对应while循环体最后一条指令，所以这里又产生了L4, 当然这条语句其实是冗余，在后面生成代码优化时我们再处理。

最后我们看看do...while...循环的实现。在inter里面创建do.go，添加代码如下：
```
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

```
do 节点实现跟while没有太大差别，只是跳转的位置稍微有些差异。接着我们在解析时要添加对do语句的处理，代码如下：
```
case lexer.DO:
		s.move_forward()
		do_node := inter.NewDo(s.lexer.Line)
		//将当前do节点加入栈，解析break语句时从栈顶拿到包围它的循环语句
		s.loop_enclosing = append(s.loop_enclosing, do_node)

		//解析do的循环体部分
		err := s.matchLexeme("{")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过{
		while_stmt := s.stmts()
		err = s.matchLexeme("}")
		if err != nil {
			panic(err)
		}
		s.move_forward() //越过}

		s.matchLexeme("while")
		s.move_forward()
		s.matchLexeme("(")
		s.move_forward()
		expr := s.bool()
		s.matchLexeme(")")
		s.move_forward()
		s.matchLexeme(";")
		s.move_forward()
		do_node.InitDo(expr, while_stmt)
		return do_node
```
完成上面代码后，我们修改一下要编译的代码，在main.go中修改如下：
```
func main() {
	/*
	 if (b < 2) {
	                        break;
	                    } else {
	                       c = c + 1;
	                    }
	*/
	source := `{int a; int b; int c; 
		        a = 3;
				b = 0;
				do {
					if (b < 2) {
	                        break;
	                    } else {
	                       c = c + 1;
	                    }
				} while (a >= 0 && b <= 4);

				c = 2;
				
	}`
	my_lexer := lexer.NewLexer(source)
	parser :=![请添加图片描述](https://img-blog.csdnimg.cn/1f14e7be0a74446dab2f689d40873c21.png)
 simple_parser.NewSimpleParser(my_lexer)
	parser.Parse()
}
```
最后我们看看运行结果：
![请添加图片描述](https://img-blog.csdnimg.cn/3e89366099254e64ae72f6f7c7d1d635.png)
我们分析一下结果，L4对应循环体内部的if语句，如果b<2不成立，那么跳转到L8,可以看到L8对应的正好是else部分语句，如果成立，那么直接进入L7,其中它有两条goto语句，第一条跳转到L5,那里对应正好是do..while循环出去后的第一条语句，goto L6是else语句块生成的跳转，它的目的是当if成立后，执行了if成立时的语句块，那么就要越过else部分，而L8就是else部分代码入口，显然这里两个goto语句是一种冗余，我们需要在代码优化部分再进行处理。

L6对应的正好就是while的判断语句，如果循环条件a>=0不成立，那么跳到L9，但是L9没有指令，因此直接进入L5,也就是跳出了循环，如果a >=0 成立，那么再判断b <= 4是否成立，不成立同样进入L9然后进入L5于是跳出循环，如果成立那么进入L4,而L4恰好就是循环体的入口，如此看来我们生成代码的逻辑基本正确。

更详细的讲解和演示请在B站搜索Coding迪斯尼，[更多干货](http://m.study.163.com/provider/7600199/index.htm?share=2&shareId=7600199)：http://m.study.163.com/provider/7600199/index.htm?share=2&shareId=7600199
