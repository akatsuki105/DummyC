package main

import (
	"flag"
	"llvm.org/llvm/bindings/go/llvm"
	"path/filepath"
	"../pkg/lexer"
	"../pkg/parser"
	"../pkg/ast"
	"../pkg/generator"
)

var (
	inputFile, outputFile string
)

func main() {
	// コマンドライン引数のパース
	var (
		o = flag.String("o", "", "output target")
	)
	flag.Parse()
	input := flag.Arg(0)
	output := *o
	input, output = parseArgs(input, output)

	// LLVMの処理
	llvm.InitializeNativeTarget()

	// ソースコードの読み込み
	source, err := readFile(input)
	if err != nil {
		panic(err)
	}

	// 解析
	l := lexer.New(source)
	p := parser.New(l)
	tu := p.Parse()
	g := generator.New()
	g.Generate(tu, input)
	mod := g.GetModule()
}

func parseArgs(input, output string) (string, string) {
	if input == "" {
		panic("please enter input file path")
	} else if filepath.Ext(input) != ".dc" {
		panic("please select .dc file")
	}

	if output == "" {
		output = input[0:len(input)-3] + ".ll"
	} else if filepath.Ext(output) != ".ll" {
		panic("please select .ll file")
	}

	return input, output
}


func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := string(b)
	return lines, nil
}