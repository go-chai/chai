// IT IS IMPORTANT TO KEEP THE FORMATTING OF THIS FILE EXACTLY AS IT IS.

package tests

import (
	"net/http"

	"github.com/go-chai/chai/chai"
)

//Simple correct comment
func Simple(string,
	string) chai.ResHandlerFunc[string, error] {

	//

	//

	//

	//

	// Simple wrong comment
	return func(w http.ResponseWriter, r *http.Request) (string, int, error) { return "", 3, nil }; return nil; return nil

	//

	//
}

//Simple2 correct comment
func Simple2(string,
	string) chai.ResHandlerFunc[string, error] {

	//

	//

	//

	//

	// Simple2 wrong comment
	return func(w http.ResponseWriter, r *http.Request) (string, int, error) { return "", 3, nil }; return nil; return nil

	//

	//
}

//NotSimple correct comment
func NotSimple(string,
	string) chai.ResHandlerFunc[string, error] {

	//

	//

	//

	//

	x := 3

	//NotSimple wrong comment
	return func(w http.ResponseWriter, r *http.Request) (string, int, error) { return "", x, nil }; return nil; return nil

	//

	//
}

//NotSimple2 correct comment
func NotSimple2(string,
	string) chai.ResHandlerFunc[string, error] {

	//

	//

	//

	//

	x := 3

	//NotSimple2 wrong comment
	return func(w http.ResponseWriter, r *http.Request) (string, int, error) { return "", x, nil }; return nil; return nil

	//

	//
}

//Simple3 correct comment
func Simple3(string,




	string) chai.ResHandlerFunc[string, error] {

	//

	//
	//

	//

	//

	// Simple3 wrong comment
	
	if false {return nil}
	
	{ return nil} 
	
	return nil; { return func(w http.ResponseWriter, r *http.Request) (string, int, error) { return "", 3, nil }; return nil; return nil}
	//

	//
}

// Simple4 outer comment
func Simple4() func() (int) {


	//


	//
	

	// Simple4 inner comment
	return func() (int) { return 1}
} 