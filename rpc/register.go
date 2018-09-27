package rpc

import (
	"encoding/json"
	"fmt"
	"net/rpc"

	"github.com/sotatek-dev/heta/core"
)

// Server ...
type Server int

// GetInfo ...
func (sv *Server) GetInfo(args *string, rep *string) error {
	res := core.GetChainInfo()
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*rep = string(b)
	return nil
}

// response from core
func (sv *Server) CreateAccount(args *string, rep *string) error {
	reqJSON := core.AccountRequestFormat{}
	err := json.Unmarshal([]byte(*args), &reqJSON)
	if err != nil {
		panic(err)
	}

	res := core.NewAccount()
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*rep = string(b)
	return nil
}

// response from core
func (sv *Server) GetBlock(args *string, rep *string) error {

	reqJSON := core.BlockRequestFormat{}
	err := json.Unmarshal([]byte(*args), &reqJSON)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	res := core.GetBlock(reqJSON.BlockNumber)
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*rep = string(b)
	return nil
}

// response from core
func (sv *Server) GetTransaction(args *string, rep *string) error {
	reqJSON := core.TransactionRequestFormat{}
	err := json.Unmarshal([]byte(*args), &reqJSON)

	if err != nil {
		panic(err)
	}

	res := core.GetTransaction(reqJSON.Txid)
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*rep = string(b)
	return nil
}

func SetupRPCServer() {
	cal := new(Server)
	rpc.Register(cal)
	rpc.HandleHTTP()
}
