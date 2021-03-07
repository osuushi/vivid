package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/osuushi/vivid/vivian"
)

func main() {
	str := os.Args[1]
	ast, err := vivian.ParseString(str)
	if err != nil {
		fmt.Printf("Failed with: %v", err)
	} else {
		pretty.Println(ast)
	}
}
