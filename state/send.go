package state

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	tp "github.com/tendermint/tendermint/identypes"
	"os"
	"strconv"
	myline"github.com/tendermint/tendermint/line"
)

func  conver2cptx(cpTxs []tp.TX,height int64) tp.TX{
	
	var content []string
	fmt.Println("cpTxs length is ",len(cpTxs))
	for i:=0;i<len(cpTxs);i++{
		marshalTx ,_:=json.Marshal(cpTxs[i])
		content=append(content,string(marshalTx))
	}
	cptx :=&tp.TX{
		Txtype:"checkpoint",
		Sender: strconv.FormatInt(height,10), //用sender记录高度
		Receiver: "",
		ID      : sha256.Sum256([]byte("checkpoint")),
		Content :content} 
    return *cptx
}
func Sendcptx(tx tp.TX, flag int) {

	res, _ := json.Marshal(tx)
	fmt.Println("-----------------sendcheckpointtx-----------------------")
	paramsJSON, err := json.Marshal(map[string]interface{}{"tx": res})
	if err != nil {
		fmt.Printf("failed to encode params: %v\n", err)
		os.Exit(1)
	}
	rawParamsJSON := json.RawMessage(paramsJSON)
	rc := &RPCRequest{
		JSONRPC: "2.0",
		ID:      "tm-bench",
		Method:  "broadcast_tx_async",
		Params:  rawParamsJSON,
	}
	c,_:=myline.UseConnect("Localhost","localhost")
	c.WriteJSON(rc)
	myline.Flag_conn["Localhost"][0]=false
}

