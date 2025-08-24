package dbg

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any, name ...string) {
	b, err := json.MarshalIndent(&v, "", "  ")
	if err != nil {
		panic("DEBUG FAILED: " + err.Error())
	}
	if len(name) > 0 {
		fmt.Printf("%s\n%s\n", name[0], string(b))
	} else {
		fmt.Println(string(b))
	}
}
