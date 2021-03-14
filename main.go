package main

import (
	"fmt"
	"os"

	"github.com/osuushi/vivid/render"
	"gopkg.in/yaml.v3"
)

func parseFile(path string) ([]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(file)

	var result interface{}
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	list, ok := result.([]interface{})
	if !ok {
		list = []interface{}{result}
	}

	return list, nil
}

func main() {
	path := os.Args[1]
	contexts, err := parseFile(path)
	if err != nil {
		panic(err)
	}

	template, err := render.MakeRow(os.Args[2])
	if err != nil {
		panic(err)
	}

	width, err := render.TerminalWidth()
	if err != nil {
		panic(err)
	}

	beam := render.DefaultBeam()

	for _, context := range contexts {
		rows, err := template.Render(width, beam, context)
		if err != nil {
			panic(err)
		}
		for _, row := range rows {
			fmt.Println(row)
		}
	}

	if err != nil {
		panic(err)
	}
}
