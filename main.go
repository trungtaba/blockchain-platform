package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/p2p/discover"
	env "github.com/joho/godotenv"
	cli "github.com/sotatek-dev/heta/cli/commands"
	"github.com/sotatek-dev/heta/core"
	"github.com/sotatek-dev/heta/heta"
	"github.com/sotatek-dev/heta/myp2p"
	"github.com/sotatek-dev/heta/rpc"
	"github.com/sotatek-dev/heta/trie"
	"github.com/sotatek-dev/heta/utils"
)

var (
	stateTrie trie.Trie
)

func start(heta heta.Heta) {
	listener, error := net.Listen("tcp", ":"+os.Getenv("HTTP_SERVICE_PORT"))
	if error != nil {
		log.Fatal("HTTP service error ", error)
	}
	defer listener.Close()
	defer heta.ChainDB.Close()

	http.Serve(listener, nil)
}

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		listenAddr   = flag.String("addr", ":30301", "Listen address")
		producerFlag = flag.Bool("producer", false, "Is set if node want to be a producer")
		bootFlag     = flag.String("bootnodes", "", "Comma separated bootnode enode URLs to seed with")
		port         = flag.Uint64("port", 9100, "P2P server port")
		netrestrict  = flag.String("netrestrict", "", "Restrict network communication to the given IP networks (CIDR masks)")
		clientID     = os.Getenv("LOCAL_CLIENT_ID")
		blockTime, _ = strconv.ParseUint(os.Getenv("BLOCK_TIME"), 10, 64)
		mutex        = &sync.Mutex{}
	)
	flag.Parse()

	if *bootFlag == "" {
		myp2p.NewBootNodeServer(*listenAddr, *netrestrict)
	} else {
		// create global db
		heta.New()
		stateTrie := core.GenStateTrie()
		heta.FullNode.StateTrie = stateTrie

		// Convert the bootnodes to internal enode representations
		var enodes []*discover.Node
		for _, boot := range strings.Split(*bootFlag, ",") {
			if url, err := discover.ParseNode(boot); err == nil {
				enodes = append(enodes, url)
			} else {
				log.Panic("Failed to parse bootnode URL", "url", boot, "err", err)
			}
		}

		serverName := "heta" + string(*port)
		server, err := myp2p.NewServer(utils.MakeName(serverName, "1.0"), *port, myp2p.Protocols, enodes, true)
		if err != nil {
			log.Panic("start server failed", err)
		}

		log.Println(server.NodeInfo().Enode)
		server.Start()

		if *producerFlag {
			core.IsProducerCandidate = true
			ticker := time.NewTicker(time.Duration(blockTime) * time.Second)
			go func() {
				for range ticker.C {
					candidates := core.Candidates
					if _, exist := candidates["1"]; !exist {
						continue
					}

					producerPrivateKey := os.Getenv("NODE_PRIVATE_KEY")
					producerKey := core.NewAccountFromPrivateKey(producerPrivateKey)
<<<<<<< HEAD

					accountBytes := heta.FullNode.StateTrie.Find(heta.FullNode.StateTrie.Root, producerKey.Key.Address[:])
					accountData := new(core.State)
					utils.Deserialize(accountBytes, accountData)
					log.Println("producer balance", accountData.Balance)
=======
>>>>>>> c838dea44f9ab4826534580d6b8aaecba2d35422

					accountBytes := heta.FullNode.StateTrie.Find(heta.FullNode.StateTrie.Root, producerKey.Key.Address[:])
					accountData := new(core.State)
					utils.Deserialize(accountBytes, accountData)
					// log.Println("producer balance", accountData.Balance)

					// mine with producers
					numberOfProducers, _ := strconv.ParseUint(os.Getenv("NUMBER_OF_PRODUCERS"), 10, 64)
					now := uint64(time.Now().Unix())
					blockTimestamp := utils.NewBlockTimestamp(now, blockTime)
					index := uint64(1) + blockTimestamp.Slot%numberOfProducers
					indexString := strconv.FormatUint(index, 10)

					// candidate
					var candidateArray []*core.Candidate
					for _, candidate := range core.Candidates {
						candidateArray = append(candidateArray, candidate)
					}
					sort.Slice(candidateArray, func(i, j int) bool {
						return candidateArray[i].Vote.Cmp(candidateArray[j].Vote) > 0
					})
					if len(candidateArray) >= 3 {
						candidateArray = candidateArray[:numberOfProducers]
					} else {
						log.Printf("WARNING: the number of candidates is %d but number of producers is %d", len(candidateArray), numberOfProducers)
					}
					candidateIndex := sort.Search(len(candidateArray), func(i int) bool {
						return candidateArray[i].Address == clientID
					})

					if clientID == indexString && candidateIndex >= 0 {
						mutex.Lock()
						cli := cli.CLI{}
						cli.CreateBlockWithTransactions(clientID, myp2p.Transactions)
						myp2p.Transactions = []*core.Transaction{} // clear old transactions
						mutex.Unlock()
					}
				}
			}()
		}
	}

	rpc.SetupHTTPServer()
	rpc.SetupRPCServer()
	start(heta.FullNode)
}
