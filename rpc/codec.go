package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/rpc"
	"sync"
)

type ServerCodec struct {
	encmutex sync.Mutex    // protects enc
	dec      *json.Decoder // for reading JSON values
	enc      *json.Encoder // for writing JSON values
	c        io.Closer
	srv      *rpc.Server
	ctx      context.Context

	// temporary work space
	req serverRequest

	// JSON-RPC clients can use arbitrary json values as request IDs.
	// Package rpc expects uint64 request IDs.
	// We assign uint64 sequence numbers to incoming requests
	// but save the original request ID in the pending map.
	// When rpc responds, we use the sequence number in
	// the response to find the original request ID.
	mutex   sync.Mutex // protects seq, pending
	seq     uint64
	pending map[uint64]*json.RawMessage
}

func NewServerCodec(conn io.ReadWriteCloser, srv *rpc.Server) rpc.ServerCodec {
	if srv == nil {
		srv = rpc.DefaultServer
	}
	// register functions for server
	srv.Register(APIRegister{})
	// new server codec that implement rpc.server codec interface
	newServerCodec := ServerCodec{
		dec:     json.NewDecoder(conn),
		enc:     json.NewEncoder(conn),
		c:       conn,
		srv:     srv,
		ctx:     context.Background(),
		pending: make(map[uint64]*json.RawMessage),
	}
	return &newServerCodec
}

// Support for single request
func (kodec *ServerCodec) ReadRequestHeader(r *rpc.Request) (err error) {
	var rawHeaderJson json.RawMessage

	if err := kodec.dec.Decode(&rawHeaderJson); err != nil {
		// ***
		kodec.enc.Encode(serverResponse{Version: "2.0", ID: &null, Error: errParse})
		// ***
		return err
	}
	// if (decode error, send response with error encoded
	// ***
	// continue to unmarshal rawdata (json.RawMessage ~~ binary data)
	// ***
	json.Unmarshal(rawHeaderJson, &kodec.req)
	// if (unmarshal error, send response with error encoded
	// ***
	r.ServiceMethod = kodec.req.Method

	// ***
	kodec.seq++
	kodec.pending[kodec.seq] = kodec.req.ID
	kodec.req.ID = nil
	r.Seq = kodec.seq
	// ***
	return nil
}

func (kodec *ServerCodec) ReadRequestBody(x interface{}) error {
	err := json.Unmarshal(*kodec.req.Params, x)

	if err != nil {
		panic(err)
	}
	return nil
}

func (kodec *ServerCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	b, ok := kodec.pending[r.Seq]
	if !ok {
		return errors.New("invalid sequence number in response")
	}
	resp := serverResponse{Version: "2.0", ID: b}

	resp.Result = x
	return kodec.enc.Encode(resp)
}

func (kodec *ServerCodec) Close() error {
	return kodec.c.Close()
}

func (APIRegister) Batch(arg BatchArg, replies *[]*json.RawMessage) (err error) {
	return nil
}
