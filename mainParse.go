package main

import (
	"go/ast"
)

type mainSource struct {
	assignments []*ast.AssignStmt // :=
	values []*ast.ValueSpec // consts, =
	literals []*ast.BasicLit // "asdf", 1234, 0xFFFFD00D
	imports []*ast.ImportSpec // all the imports
}

func parseMainSourceFromAST(node ast.Node) *mainSource{
	ret_val := mainSource{}

	// depth first iterate over each node in the AST Tree
	ast.Inspect(node, func (n ast.Node) bool{

		// if the node is an assignment
		assignments, ok := n.(*ast.AssignStmt)
		if ok {
			// add it to the list of assignments
			ret_val.assignments = append(ret_val.assignments, assignments)
			return true
		}


		imports, ok := n.(*ast.ImportSpec)
		if ok {
			ret_val.imports = append(ret_val.imports, imports)
			return true
		}
		var import_names []string
		for i := range ret_val.imports {
			import_names = append(import_names, ret_val.imports[i].Path.Value)
		}
		// ditto
		values, ok := n.(*ast.ValueSpec)
		if ok  {
			// ditto
			ret_val.values = append(ret_val.values, values)
			return true
		}

		// if the node is a literal
		vars, ok := n.(*ast.BasicLit)
		if ok && !strContains(import_names, vars.Value) {
			// add it to our list of literals
			ret_val.literals = append(ret_val.literals, vars)
			return true // our evaluation is done, don't recheck the same node
		}
		return true
	})
	return  &ret_val
}

