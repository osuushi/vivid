package render

import (
	"fmt"

	"github.com/osuushi/vivid/vivian"
)

type Row struct {
	Cells []Cell
}

func MakeRow(input string) (*Row, error) {
	ast, err := vivian.ParseString(input)
	if err != nil {
		return nil, err
	}
	return makeRowFromAst(ast)
}

func makeRowFromAst(ast *vivian.Ast) (*Row, error) {
	return nil, fmt.Errorf("TODO")
}
