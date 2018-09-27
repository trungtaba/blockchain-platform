package cli

import (
	"fmt"
	"github.com/sotatek-dev/heta/core"
	"log"
)

// CreateBlockchain ...
func (cli *CLI) CreateBlockchain(address string, clientID string) {
	if !core.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := core.CreateBlockchain(address, clientID)
	defer bc.GetDB().Close()

	fmt.Println("Done!")
}
