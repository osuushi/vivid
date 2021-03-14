package main

import (
	"fmt"
	"os"
	"time"

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

	width, err := render.TerminalWidth()
	if err != nil {
		panic(err)
	}

	beam := render.DefaultBeam()

	start := time.Now()
	templates := []*render.Row{}
	for _, arg := range os.Args[2:] {
		template, err := render.MakeRow(arg)
		templates = append(templates, template)
		if err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(start)
	start = time.Now()

	for _, context := range contexts {
		for _, template := range templates {
			rows, err := template.Render(width, beam, context)
			if err != nil {
				panic(err)
			}

			// Don't measure actual stdout time
			elapsed += time.Since(start)
			for _, row := range rows {
				fmt.Println(row)
			}
			start = time.Now()
		}
	}

	if err != nil {
		panic(err)
	}

	elapsed += time.Since(start)

	fmt.Println("Rendering time:", elapsed)
	fmt.Println("Per row:", time.Duration(int64(elapsed)/int64(len(contexts))))
}
