package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net/rpc"
)

var null = json.RawMessage([]byte("null"))

type serverRequest struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
	ID      *json.RawMessage `json:"id"`
}

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type serverResponse struct {
	Version string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  interface{}      `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}

// Ctx can be embedded into your struct with RPC method parameters (if
// that method parameters type is a struct) to make it implement
// WithContext interface and thus have access to request context.
type Ctx struct {
	ctx context.Context
}

// BatchArg is a param for internal RPC JSONRPC2.Batch.
type BatchArg struct {
	srv  *rpc.Server
	reqs []*json.RawMessage
	Ctx
}

type APIRegister struct {
}

func (r *serverRequest) reset() {
	r.Version = ""
	r.Method = ""
	r.Params = nil
	r.ID = nil
}

func (r *serverRequest) UnmarshalJSON(raw []byte) error {
	r.reset()
	type req *serverRequest
	if err := json.Unmarshal(raw, req(r)); err != nil {
		return errors.New("bad request")
	}

	var o = make(map[string]*json.RawMessage)
	if err := json.Unmarshal(raw, &o); err != nil {
		return errors.New("bad request")
	}
	if o["jsonrpc"] == nil || o["method"] == nil {
		return errors.New("bad request")
	}
	_, okID := o["id"]
	_, okParams := o["params"]
	if len(o) == 3 && !(okID || okParams) || len(o) == 4 && !(okID && okParams) || len(o) > 4 {
		return errors.New("bad request")
	}
	if r.Version != "2.0" {
		return errors.New("bad request")
	}
	if okParams {
		if r.Params == nil || len(*r.Params) == 0 {
			return errors.New("bad request")
		}
		switch []byte(*r.Params)[0] {
		case '[', '{':
		default:
			return errors.New("bad request")
		}
	}
	if okID && r.ID == nil {
		r.ID = &null
	}
	if okID {
		if len(*r.ID) == 0 {
			return errors.New("bad request")
		}
		switch []byte(*r.ID)[0] {
		case 't', 'f', '{', '[':
			return errors.New("bad request")
		}
	}
	return nil
}
