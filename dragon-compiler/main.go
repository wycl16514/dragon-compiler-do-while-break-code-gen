package main

import (
	"lexer"
	"simple_parser"
	//"inter"
	//"fmt"
)

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
	parser := simple_parser.NewSimpleParser(my_lexer)
	parser.Parse()
}
