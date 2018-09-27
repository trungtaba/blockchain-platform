package cli

import (
	"encoding/json"
	"fmt"
	"github.com/sotatek-dev/heta/core"
)

// Get Block ...
func (cli *CLI) GetTransaction(txid string) string {
	info := core.GetTransaction(txid)
	b, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
		return "Error: Something went wrong"
	}
	fmt.Println(string(b))
	return string(b)
}
