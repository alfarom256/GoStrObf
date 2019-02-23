package main

import (
	"go/token"
	"go/parser"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"go/printer"
	"os"
)

func main(){
	welcome_str := "Hello world!!!"
	fmt.Printf("%s\n", welcome_str)
	fPath := "C:/Users/Michael/go/src/helloWorld/main.go"

	mySource := mainSource{}
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, fPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal("fuckin a, then!")
		panic(err)
	}
	mySource = *parseMainSourceFromAST(node)
	parseFile(&mySource)
	myStrObf := xorStub()
	new_src := appendStub(node, fset, myStrObf)
	var obfStrings []string
	for i := range mySource.literals{
		obfStrings = append(obfStrings, mySource.literals[i].Value)
	}
	newObfStrings := generateStrings(obfStrings, myStrObf)
	//new_src = appendStringLiterals(new_src, fset, *newObfStrings)
	//printer.Fprint(os.Stdout, fset, node)
	fmt.Printf("\n\n\n\nNEW SOURCE CODE\n\n\n\n\n")
	printer.Fprint(os.Stdout, fset, new_src)
	// now we've added the temporary variables at the bottom, let's reparse the AST and try to swap nodes
	new_src = replaceTempVarStrings(new_src, *newObfStrings)
	printer.Fprint(os.Stdout, fset, new_src)
}

func parseFile(mySource *mainSource){
	literals := mySource.literals
	import_strings := make([]string, len(mySource.imports))
	for i := range mySource.imports{
		import_name :=  mySource.imports[i].Path.Value
		import_strings = append(import_strings, import_name)
		fmt.Printf("Ignoring string for Import: %s\n", import_name)
	}
	for i := range literals{
		// if it's a string but NOT an import, print it
		if literals[i].Kind == token.STRING && !strContains(import_strings, literals[i].Value){
			fmt.Printf("Found string: %s\n", literals[i].Value)
		}
	}
	// let's print all of the Objects in the assignment list and examine the values
	for i := range mySource.assignments{
		spew.Dump(mySource.assignments[i])
	}

}