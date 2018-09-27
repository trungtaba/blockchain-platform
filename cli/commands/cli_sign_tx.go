package cli

import (
	"encoding/json"
	"fmt"
	"github.com/sotatek-dev/heta/core"
)

// HandleSignTxRequest ...
func (cli *CLI) SignTransaction(rawSign string) string {
	info := core.HandleSignTxRequest([]byte(rawSign))
	b, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
		return "Error: Something went wrong"
	}
	fmt.Println(string(b))
	return string(b)
}
