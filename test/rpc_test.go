package main

import (
	"fmt"
	rpcHeta "github.com/sotatek-dev/heta/rpc"
	"log"
	"net"
	"net/rpc"
	"testing"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

type ArithAddResp struct {
	ID     interface{} `json:"id"`
	Result Reply       `json:"result"`
	Error  interface{} `json:"error"`
}

func (t *Arith) Add(args *Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(args *Args, reply *Reply) error {
	if args.B == 0 {
		return fmt.Errorf("divide for 0")
	}
	reply.C = args.A / args.B
	return nil
}

func (t *Arith) Error(args *Args, reply *Reply) error {
	panic("ERROR")
}

func TestConnections(t *testing.T) {
	fmt.Println("Listen on port 9009 ...")
	rpc.Register(new(Arith))
	rwd, error := net.Listen("tcp", ":9009")
	if error != nil {
		log.Fatal("HTTP service error ", error)
	}
	defer rwd.Close()

	// Server start
	go func() {
		conn, err := rwd.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting codec ...")
		rpc.ServeCodec(rpcHeta.NewServerCodec(conn, nil))
	}()

	// Client listen
	connN, err := net.Dial("tcp", "127.0.0.1:9009")
	if err != nil {
		panic(err)
	}
	defer connN.Close()
	client := rpcHeta.NewClient(connN)
	// Synchronous calls
	args := &rpcHeta.Args{7, 8}
	reply := new(rpcHeta.Reply)
	err = client.Call("Arith.Add", args, reply)
	if err != nil {
		fmt.Printf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		fmt.Printf("Add: got %d expected %d\n", reply.C, args.A+args.B)
	}
}
