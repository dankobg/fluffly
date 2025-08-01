package dbg

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any, s ...string) {
	b, err := json.MarshalIndent(&v, "", "  ")
	if err != nil {
		panic("DEBUG FAILED: " + err.Error())
	}
	if len(s) > 0 {
		fmt.Printf("%s\n%s\n", s[0], string(b))
	} else {
		fmt.Println(string(b))
	}
}
