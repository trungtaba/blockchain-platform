package cli

import (
	"fmt"

	"github.com/sotatek-dev/heta/core"
)

// CreateAccount ...
func (cli *CLI) CreateAccount(clientID string) {
	account := core.NewAccount()

	fmt.Println("===================== CREATE NEW ACCOUNT =====================")
	fmt.Printf("Your new address:  %s\n", account.Key.Address.ToString())
	fmt.Printf("Private key:       %x\n", account.Key.PrivateKey.D.Bytes())
	fmt.Println("==============================================================")
}
