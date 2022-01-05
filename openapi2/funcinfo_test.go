package openapi2_test

import (
	"fmt"
	"go/ast"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chai/chai/openapi2"
)

//Comment
func Simple() string {
	return "hello"
}

func TT() float32 {
	fifio := func(mm int) (string, int, string) {
		println(mm)

		if true {
			return "", http.StatusAccepted, "123"
		}
		return "45m", http.StatusForbidden, "4435"
	}

	fifio(123)

	return 1.23
}

type ZZ struct {
	A string `json:"a"`
}

func ExampleGetFuncInfo() {
	fi := openapi2.GetFuncInfo(Simple)

	fmt.Println(fi.Comment)
	fmt.Println(fi.File)
	fmt.Println(fi.Line)

	// ast.ReturnStmt

	// ast.Walk()

	ast.Inspect(fi.ASTFile, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			fmt.Println("x.Name.Name")
			spew.Dump(x.Name.Name)
			// spew.Dump(x.Body)
			spew.Dump(x.Doc)
		case *ast.ReturnStmt:
			fmt.Println("x.Results")
			spew.Dump(x.Results)
		}

		return true
	})
	// Output:
	// Comment
}
