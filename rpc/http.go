package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/sotatek-dev/heta/core"
)

var func_facts = map[string]interface{}{
	"get_info":             core.HandleChainInfoReq,
	"get_block":            core.HandleBlockInfoReq,
	"get_balance":          core.HandleBalanceReq,
	"create_address":       core.HandleAddressRegReq,
	"get_transaction":      core.HandleTxInfoRequest,
	"send_raw_transaction": core.HandleTxRequest,
	"sign_transaction":     core.HandleSignTxRequest,
	"get_peers":            core.HandlePeersReq,
	"get_total_staked":     core.HandleTotalStaked,
	"list_producers":       core.HandleListProducers,
	"get_blockByhash":      core.HandleBlockInfoByHashReq,
}

func getHandler(route string) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		fmt.Println(string(b[:]))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		defer r.Body.Close()
		value := reflect.ValueOf(b)
		inputs := make([]reflect.Value, 1)
		inputs[0] = value

		msg := reflect.ValueOf(func_facts[route]).Call(inputs)[0].Interface()
		output, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(output)
	})
	return handler
}

func SetupHTTPServer() {
	for route, _ := range func_facts {
		handler := getHandler(route)
		http.Handle("/"+route, handler)
	}
}
