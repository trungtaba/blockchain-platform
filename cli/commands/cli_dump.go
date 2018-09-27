package cli

import (
	"fmt"

	"github.com/sotatek-dev/heta/core"
	"github.com/sotatek-dev/heta/types"
)

// DumpChainData ...
func (cli *CLI) DumpChainData(clientID string) {
	bc := core.NewBlockChain(clientID)

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("========================================== BLOCK %d ==========================================\n", block.GetHeight())
		fmt.Printf("  Hash:         %x\n", block.GetHash())
		fmt.Printf("  Prev. block:  %x\n", block.GetPrevBlockHash())
		fmt.Printf("  Transactions: \n")
		for i, tx := range block.GetTransactions() {
			fmt.Printf("    %d: %x\n", i, tx)
		}
		fmt.Printf("==============================================================================================\n\n")

		if types.IsZeroHash(block.PrevBlockHash) {
			break
		}
	}
}
