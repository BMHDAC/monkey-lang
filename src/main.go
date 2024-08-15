package main

import (
	"fmt"
	"monkey/src/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Monkey lang. Hello: \n", user.Username)
	fmt.Println("Type the source code")
	repl.Start(os.Stdin, os.Stdout)
}
