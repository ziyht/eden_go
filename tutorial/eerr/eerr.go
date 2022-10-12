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

	eerr.PrintSourceColor(err, 0)

	f, _ := eerr.StackCause(err)
	fmt.Println(f.String())

	f = eerr.Call(0)
	fmt.Println(f.String())
}

