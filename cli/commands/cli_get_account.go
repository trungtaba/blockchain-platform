package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sotatek-dev/heta/core"
	"github.com/sotatek-dev/heta/types"
)

// 35396bc9aa57bee30ebf81e813950dc4f631a308
// GetBlockByNumber ...
func (cli *CLI) GetAccount(accountString string) string {
	accountBytes, err := hex.DecodeString(accountString)
	if err != nil {
		panic(err)
	}
	account := new(types.Address)
	account.SetBytes(accountBytes)

	info := core.GetAccount(*account)
	b, err := json.Marshal(info)
	fmt.Println(string(b))
	if err != nil {
		fmt.Println(err)
		return "Error: Something went wrong"
	}
	return string(b)
}
