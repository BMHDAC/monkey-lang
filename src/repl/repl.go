package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/src/lexer"
	"monkey/src/parser"
)

/*

 */

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserError(out, p.Errors())
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserError(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
