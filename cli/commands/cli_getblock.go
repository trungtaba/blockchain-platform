package cli

import (
	"encoding/json"
	"fmt"

	"github.com/sotatek-dev/heta/core"
)

// GetBlockByNumber ...
func (cli *CLI) GetBlockByNumber(number int) string {
	info := core.GetBlockByNumber(uint64(int64(number)))
	b, err := json.Marshal(info)
	fmt.Println(string(b))
	if err != nil {
		fmt.Println(err)
		return "Error: Something went wrong"
	}
	return string(b)
}
