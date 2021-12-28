package chai

import (
	"encoding/json"
	"fmt"
)

func LogYAML(v any) {
	bytes, err := MarshalYAML(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

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
