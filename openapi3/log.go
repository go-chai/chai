package openapi3

import (
	"encoding/json"
	"fmt"
	"os"
)

func LogYAML(v any) {
	bytes, err := MarshalYAML(v)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, string(bytes))

	return
}

func LogJSON(v any) {
	bytes, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}
