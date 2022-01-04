package openapi2_test

import (
	"fmt"

	"github.com/go-chai/chai/openapi2"
)

//Comment
func Simple() string {
	return "hello"
}

type ZZ struct {
	A string `json:"a"`
}

func ExampleGetFuncInfo() {
	fmt.Println(openapi2.GetFuncInfo(Simple).Comment)
	// Output:
	// Comment
}
