package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"fmt"
)

func traceLocks(file *ast.File) {
	statements := findAllLockStatements(file)
	if len(statements) == 0 {
		return
	}
	addAstraceInjectionImport(file)
	wrapAllStatements(file, statements)
}

func findAllLockStatements(file *ast.File) []ast.Stmt {
	statements := []ast.Stmt{}
	ast.Inspect(file, func(node ast.Node) bool {
		if statement, ok := node.(ast.Stmt); ok {
			if isStatementALock(statement) {
				statements = append(statements, statement)
			}
		}
		return true
	})
	return statements
}

func isStatementALock(statement ast.Stmt) bool {
	switch stmt := statement.(type) {
	case *ast.ExprStmt:
		if call, ok := stmt.X.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if isObjectAMutex(ident.Obj) {
						if sel.Sel.Name == "Lock" || sel.Sel.Name == "RLock" {
							return true
						}
					}
				}
			}
		}
	case *ast.SendStmt:
		return true
	case *ast.AssignStmt:
		for _, expr := range stmt.Rhs {
			if unary, ok := expr.(*ast.UnaryExpr); ok {
				if unary.Op == token.ARROW {
					return true
				}
			}
		}
	}
	return false
}

func isObjectAMutex(object *ast.Object) bool {
	if object == nil {
		return false
	}
	if field, ok := object.Decl.(*ast.Field); ok {
		if star, ok := field.Type.(*ast.StarExpr); ok {
			if sel, ok := star.X.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok {
					if pkg.Name != "sync" {
						return false
					}
				}
				if sel.Sel.Name == "Mutex" {
					return true
				}
				if sel.Sel.Name == "RWMutex" {
					return true
				}
			}
		}
	}
	return false
}

func wrapAllStatements(file *ast.File, statements []ast.Stmt) {
	ast.Inspect(file, func(node ast.Node) bool {
		if block, ok := node.(*ast.BlockStmt); ok {
			wrapStatementsInBlock(block, statements)
			return false
		}
		return true
	})
}

func wrapStatementsInBlock(block *ast.BlockStmt, statements []ast.Stmt) {
	for i := 0; i < len(block.List); i++ {
		for _, statement := range statements {
			if block.List[i] == statement {
				randomId := rand.Uint64()
				expr, _ := parser.ParseExpr(fmt.Sprintf("astraceInjection.BeforeLock(%d)", randomId))
				block.List[i] = &ast.ExprStmt{X: expr}
				block.List = append(block.List, nil, nil)
				copy(block.List[i+3:], block.List[i+1:])
				block.List[i+1] = statement
				expr, _ = parser.ParseExpr(fmt.Sprintf("astraceInjection.AfterLock(%d)", randomId))
				block.List[i+2] = &ast.ExprStmt{X: expr}
				i = i + 2
			}
		}
	}
}

func addAstraceInjectionImport(file *ast.File) {
	const importPath = `"github.com/orbs-network/go-ast-trace/injection"`
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			foundImport := false
			for _, spec := range genDecl.Specs {
				if importSpec, ok := spec.(*ast.ImportSpec); ok {
					if importSpec.Path.Value == importPath {
						foundImport = true
					}
				}
			}
			if !foundImport {
				genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
					Name: &ast.Ident{
						Name: "astraceInjection",
					},
					Path: &ast.BasicLit{
						Kind: token.STRING,
						Value: importPath,
					},
				})
				genDecl.Lparen = 1 // must be nonzero
			}
		}
	}
}