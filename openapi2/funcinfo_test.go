package openapi2_test

import (
	"fmt"

	"github.com/go-chi/docgen"
)

//Comment
func Simple() string {
	return "hello"
}

type ZZ struct {
	A string `json:"a"`
}

func ExampleGetFuncInfo() {
	fmt.Println(docgen.GetFuncInfo(Simple).Comment)
	// Output:
	// Comment
}
