package myp2p

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
	log "github.com/inconshreveable/log15"
	"github.com/sotatek-dev/heta/core"
	"github.com/sotatek-dev/heta/crypto"
	"github.com/sotatek-dev/heta/utils"
)

var (
	chainInfoMessage        uint64
	requestBlockMessage     uint64 = 1
	ackRequestBlockMessage  uint64 = 2
	sendTxMessage           uint64 = 3
	sendCandidateMessage    uint64 = 4
	requestCandidateMessage uint64 = 5
	voteMessage             uint64 = 6
	mutex                          = &sync.Mutex{}
	// Transactions ...
	Transactions = []*core.Transaction{}
)

var (
	proto = p2p.Protocol{
		Name:    "sync",
		Version: 1,
		Length:  10, // number of message code
		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
			go SendCandidate(rw)
			go SendChainInfo(rw)
			go SendVote(rw)
			ReceiveMessage(rw)
			return nil
		},
	}
	// Protocols ...
	Protocols = []p2p.Protocol{proto}
)

// ReceiveMessage ...
func ReceiveMessage(rw p2p.MsgReadWriter) {
	for {
		received, err := rw.ReadMsg()
		if err != nil {
			_ = fmt.Errorf("Receive message fail: %v", err)
			return
		}
		mutex.Lock()
		myChainInfo := core.GetChainInfo()
		mutex.Unlock()

		switch received.Code {
		case chainInfoMessage:
			handleChainInfo(rw, received, myChainInfo)
		case requestBlockMessage:
			handleBlockRequest(rw, received, myChainInfo)
		case ackRequestBlockMessage:
			handleBlockReturn(rw, received)
		case sendTxMessage:
			handleSendTx(rw, received)
		case sendCandidateMessage:
			handleSendCandidate(rw, received)
		case requestCandidateMessage:
			handleCandidateRequest(rw)
		case voteMessage:
			handleVoteRequest(rw, received)
		}
	}
}

// SendChainInfo ...
func SendChainInfo(rw p2p.MsgReadWriter) {
	ticker := time.NewTicker(5 * time.Second) // 5s
	go func() {
		for range ticker.C {
			mutex.Lock()
			message := core.GetChainInfo()
			mutex.Unlock()
			err := p2p.Send(rw, chainInfoMessage, message)
			if err != nil {
				_ = fmt.Errorf("Send message fail: %v", err)
				return
			}

			// create dummy transaction
			account := core.NewAccount()
			tx := core.NewTX(&account.Key.Address, uint64(0), &account.Key.Address, new(big.Int), []byte{})
			err = p2p.Send(rw, sendTxMessage, tx)
			if err != nil {
				_ = fmt.Errorf("Send message fail: %v", err)
				return
			}
		}
	}()
}

// SendCandidate ...
func SendCandidate(rw p2p.MsgReadWriter) {
	if !core.IsProducerCandidate {
		p2p.Send(rw, requestCandidateMessage, core.Candidate{})
		return
	}
	clientID := os.Getenv("LOCAL_CLIENT_ID")
	now := uint64(time.Now().Unix())
	candidate := core.Candidate{Timestamp: now, Address: clientID, Vote: new(big.Int)}
	p2p.Send(rw, sendCandidateMessage, candidate)
}

// SendVote ...
func SendVote(rw p2p.MsgReadWriter) {
	ticker := time.NewTicker(23 * time.Second)
	unvoteTicker := time.NewTicker(43 * time.Second)
	nodeKeyHex := os.Getenv("NODE_PRIVATE_KEY")
	nodeKey, err := crypto.HexToECDSA(nodeKeyHex)
	if err != nil {
		log.Error(err.Error())
	}
	go func() {
		for range ticker.C {
			if core.IsProducerCandidate {
				return
			}
			voteTo(rw, "2", nodeKey, true) // vote
		}
	}()

	go func() {
		for range unvoteTicker.C {
			if core.IsProducerCandidate {
				return
			}
			voteTo(rw, "2", nodeKey, false) // unvote
		}
	}()
}

func handleChainInfo(rw p2p.MsgReadWriter, received p2p.Msg, myChainInfo core.InfoResponse) {
	var chainInfo core.InfoResponse
	err := received.Decode(&chainInfo)
	if err != nil {
		fmt.Println(err)
	}
	if myChainInfo.LastBlock > chainInfo.LastBlock {
		p2p.Send(rw, chainInfoMessage, myChainInfo)
	} else if myChainInfo.LastBlock < chainInfo.LastBlock {
		for i := myChainInfo.LastBlock + 1; i <= chainInfo.LastBlock; i++ {
			p2p.Send(rw, requestBlockMessage, i)
		}
	}
}

func handleBlockRequest(rw p2p.MsgReadWriter, received p2p.Msg, myChainInfo core.InfoResponse) {
	var blockNumber uint64
	err := received.Decode(&blockNumber)
	if err != nil {
		log.Error(err.Error())
	}

	if blockNumber > myChainInfo.LastBlock {
		return
	}

	mutex.Lock()
	block := core.GetBlockByNumber(blockNumber)
	mutex.Unlock()
	err = p2p.Send(rw, ackRequestBlockMessage, block)
	if err != nil {
		fmt.Println(err)
	}
}

func handleBlockReturn(rw p2p.MsgReadWriter, received p2p.Msg) {
	var block core.Block
	err := received.Decode(&block)
	if err != nil {
		log.Error(err.Error())
	}

	clientID := os.Getenv("LOCAL_CLIENT_ID")
	mutex.Lock()
	bc := core.NewBlockChain(clientID)
	isAdded := bc.AddBlock(&block)
	mutex.Unlock()
	if isAdded {
		msg := fmt.Sprintf("Accept Block Hash: %x\n", block.GetHash())
		log.Debug(msg)
	}
}

func handleSendTx(rw p2p.MsgReadWriter, received p2p.Msg) {
	var tx core.Transaction
	err := received.Decode(&tx)
	if err != nil {
		log.Error(err.Error())
	}
	Transactions = append(Transactions, &tx)
}

func handleSendCandidate(rw p2p.MsgReadWriter, received p2p.Msg) {
	var candidate core.Candidate
	err := received.Decode(&candidate)
	if err != nil {
		log.Error(err.Error())
	}
	candidates := core.Candidates
	if oldCandidate, exist := candidates[candidate.Address]; !exist {
		core.Candidates[candidate.Address] = &candidate
		for _, candidate := range core.Candidates {
			p2p.Send(rw, sendCandidateMessage, candidate)
		}
	} else {
		if oldCandidate.Timestamp < candidate.Timestamp {
			core.Candidates[candidate.Address] = &candidate
			p2p.Send(rw, sendCandidateMessage, candidate)
		} else if oldCandidate.Timestamp > candidate.Timestamp {
			p2p.Send(rw, sendCandidateMessage, oldCandidate)
		}
	}
}

func handleCandidateRequest(rw p2p.MsgReadWriter) {
	for _, candidate := range core.Candidates {
		p2p.Send(rw, sendCandidateMessage, candidate)
	}
}

func handleVoteRequest(rw p2p.MsgReadWriter, received p2p.Msg) {
	var voteReceiveMessage VoteMsg
	err := received.Decode(&voteReceiveMessage)
	if err != nil {
		log.Error(err.Error())
	}

	vote := voteReceiveMessage.VoteContent
	signature := voteReceiveMessage.Signature
	voteHash := utils.Hash(vote)
	messageOk := crypto.VerifySignature(vote.PublicKey, voteHash[:], signature)

	if !messageOk {
		return
	}

	// get address of voter
	publicKey, err := crypto.UnmarshalPubkey(vote.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)
	account := core.GetAccount(address)
	balance := new(big.Int)
	balance.SetString(account.Balance, 10)
	voter := core.Voter{Address: address, Timestamp: vote.Timestamp, VoteBalance: balance}

	// update candidate
	mutex.Lock()
	oldCandidate := core.Candidates[vote.Address]
	if oldCandidate.Timestamp < vote.Timestamp {
		oldCandidate.Timestamp = vote.Timestamp
		oldCandidate.Voters = updateVoters(oldCandidate.Voters, voter, vote.IsUnvote)
		oldCandidate.Vote = getVoteValue(oldCandidate.Voters)
		core.Candidates[vote.Address] = oldCandidate
		p2p.Send(rw, sendCandidateMessage, oldCandidate)
	}
	mutex.Unlock()
}

func updateVoters(voters []core.Voter, newVoter core.Voter, isUnVote bool) []core.Voter {
	isUpdated := false
	for index, voter := range voters {
		if voter.Address == newVoter.Address && voter.Timestamp < newVoter.Timestamp {
			if isUnVote {
				newVoter.VoteBalance = new(big.Int)
			}

			voters[index] = newVoter
			isUpdated = true
			break
		}
	}
	if !isUpdated {
		voters = append(voters, newVoter)
	}
	return voters
}

func getVoteValue(voters []core.Voter) *big.Int {
	result := new(big.Int)
	for _, voter := range voters {
		result.Add(result, voter.VoteBalance)
	}
	return result
}

func voteTo(rw p2p.MsgReadWriter, receiverID string, nodeKey *ecdsa.PrivateKey, isVote bool) {
	pubkeyByte := crypto.FromECDSAPub(&nodeKey.PublicKey)
	voteContent := VoteContent{Address: receiverID, Timestamp: uint64(time.Now().Unix()), Vote: big.NewInt(1), PublicKey: pubkeyByte}

	if !isVote {
		voteContent.IsUnvote = true
	}

	voteContentHash := utils.Hash(voteContent)
	signature, err := crypto.Sign(nodeKey, voteContentHash[:])
	if err != nil {
		log.Error(err.Error())
	}
	vote := VoteMsg{VoteContent: voteContent, Signature: signature}
	p2p.Send(rw, voteMessage, vote)
}
