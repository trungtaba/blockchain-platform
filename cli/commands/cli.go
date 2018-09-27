package cli

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"os"
)

const (
	cmdInit              = "init"
	cmdGetInfo           = "get_info"
	cmdGetBlock          = "get_block"
	cmdGetBalance        = "get_balance"
	cmdGetTransaction    = "get_transaction"
	cmdSignTransaction   = "sign_transaction"
	cmdCreateAddress     = "create_address"
	cmdPrintChain        = "printchain"
	cmdCreateDummmyBlock = "createdummyblock"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("USAGE:")
	fmt.Println("  init -address <ADDRESS>                     Create a new blockchain and grant pre-minted coins to <ADDRESS>")
	fmt.Println("  get_info                                    Get current blockchain's general information")
	fmt.Println("  get_block -block <BLOCK_NUMBER>             Get block information")
	fmt.Println("  get_balance -address <ADDRESS>              Get balance information")
	fmt.Println("  get_transaction -id <TX_HASH>               Get transaction information")
	fmt.Println("  sign_transaction -data <JSON_DATA>          Sign transaction")
	fmt.Println("  create_address                              Generates a new address")

	fmt.Println("COMMANDS FOR DEBUGGING:")
	fmt.Println("  printchain                                  Print all the blocks of the blockchain")
	fmt.Println("  createdummyblock -n <NUMBER>              	  Generates <NUMBER> of dummy transactions")
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	// Validate cli arguments
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	// Validate environment variables
	clientID := os.Getenv("LOCAL_CLIENT_ID")
	if clientID == "" {
		fmt.Printf("LOCAL_CLIENT_ID is not set in environment variables!")
		os.Exit(1)
	}

	switch os.Args[1] {
	case cmdInit:
		cli.handleCreateBlockchain(clientID)
		return
	case cmdGetInfo:
		cli.handleGetInfo(clientID)
		return
	case cmdCreateAddress:
		cli.handleCreateAddress(clientID)
		return
	case cmdPrintChain:
		cli.handleDumpChain(clientID)
		return
	case cmdGetBlock:
		cli.handleGetBlock(clientID)
	case cmdGetBalance:
		cli.handleGetBalance(clientID)
	case cmdGetTransaction:
		cli.handleGetTransaction(clientID)
	case cmdSignTransaction:
		cli.handleSignTransaction(clientID)
	case cmdCreateDummmyBlock:
		cli.handeCreateDummyBlock(clientID)
		return

	default:
		fmt.Printf("Invalid command: %s\n", os.Args[1])
		cli.printUsage()
		os.Exit(1)
	}

}

func (cli *CLI) handleGetBlock(clientID string) {
	cmd := flag.NewFlagSet(cmdGetBlock, flag.ExitOnError)
	blockNum := cmd.String("block", "", "Require block number to find")

	err := cmd.Parse(os.Args[2:])

	if err != nil {
		log.Panic(err)
	}
	blockNumInt, err := strconv.Atoi(*blockNum)
	cli.GetBlockByNumber(blockNumInt)
}

func (cli *CLI) handleGetBalance(clientID string) {
	cmd := flag.NewFlagSet(cmdGetBalance, flag.ExitOnError)
	accountString := cmd.String("address", "", "Require address to find")

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
	}
	cli.GetAccount(*accountString)
}

func (cli *CLI) handleGetTransaction(clientID string) {
	cmd := flag.NewFlagSet(cmdGetTransaction, flag.ExitOnError)
	txid := cmd.String("id", "", "Require transaction id to find")

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
	}
	cli.GetTransaction(*txid)
}

func (cli *CLI) handleSignTransaction(clientID string) {
	cmd := flag.NewFlagSet(cmdSignTransaction, flag.ExitOnError)
	rawSign := cmd.String("data", "", "Require input data to sign")

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
	}
	cli.SignTransaction(*rawSign)
}

func (cli *CLI) handleDumpChain(clientID string) {
	cli.DumpChainData(clientID)
}

func (cli *CLI) handleGetInfo(clientID string) {
	result := cli.GetChainInfo()
	fmt.Println(result)
}

func (cli *CLI) handleCreateAddress(clientID string) {
	cli.CreateAccount(clientID)
}

func (cli *CLI) handleCreateBlockchain(clientID string) {
	cmd := flag.NewFlagSet(cmdInit, flag.ExitOnError)
	address := cmd.String("address", "", "The address to send genesis block reward to")

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
	}

	cli.CreateBlockchain(*address, clientID)
}

func (cli *CLI) handeCreateDummyBlock(clientID string) {
	cmd := flag.NewFlagSet(cmdCreateDummmyBlock, flag.ExitOnError)
	numOfTxs := cmd.Int("n", 1, "The number of dummy transactions will be created")

	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Will create %d dummy transactions.", *numOfTxs)
	cli.CreateDummyBlock(clientID, *numOfTxs)
}
