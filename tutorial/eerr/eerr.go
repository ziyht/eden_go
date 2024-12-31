package main

import (
	"fmt"

	"github.com/ziyht/eden_go/eerr"
)

func raise() error {
	return eerr.New("new err")
}

func main() {

	err := raise()

	eerr.PrintSourceColor(err, 1)

	f, _ := eerr.StackCause(err)
	fmt.Println(f.String())

	f = eerr.Call(1)
	fmt.Println(f.String())
}

