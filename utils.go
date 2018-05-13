package uzan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func PrintResult(res []byte, err error) {
	if err != nil {
		panic(err)
	}

	response, err := ParseRawResponse(res)

	jsonBytes, _ := json.Marshal(response)

	var out bytes.Buffer
	err := json.Indent(&out, jsonBytes, "", " ")
	if err != nil {
		fmt.Println(string(res))
	} else {
		fmt.Println(out.String())
	}
}

func PrintObject(ret interface{}) {
	spew.Config.Indent = " "
	spew.Config.SortKeys = true
	spew.Config.DisableCapacities = true
	spew.Dump(ret)
}
