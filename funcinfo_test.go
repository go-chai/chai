package chai_test

import (
	"fmt"

	"github.com/go-chi/docgen"
)

//Comment
func Simple() string {
	return "hello"
}

func ExampleGetFuncInfo() {
	fmt.Println(docgen.GetFuncInfo(Simple).Comment)
	// Output:
	// Comment
}
