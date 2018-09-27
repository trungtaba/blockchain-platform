package cli

import (
	"encoding/json"
	"fmt"

	"github.com/sotatek-dev/heta/core"
)

// GetChainInfo ...
func (cli *CLI) GetChainInfo() string {
	info := core.GetChainInfo()
	b, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
		return "Error: Something went wrong"
	}
	return string(b)
}
