package state

import (
	"fmt"
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"os"
	"net/http"
	"io"
	"bytes"
	tp "github.com/tendermint/tendermint/identypes"
	

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
func  Sendcptx(tx tp.TX,flag int){

	res, _ := json.Marshal(tx)
	fmt.Println("-----------------sendcheckpointtx-----------------------")
	client := &http.Client{}
	requestBody := new(bytes.Buffer)
	paramsJSON, err := json.Marshal(map[string]interface{}{"tx": res})	       
	if err != nil {
		fmt.Printf("failed to encode params: %v\n", err)
		os.Exit(1)
	}
	rawParamsJSON := json.RawMessage(paramsJSON)
	rc:=&RPCRequest{
		JSONRPC: "2.0",
		ID:      "tm-bench",
		Method:  "broadcast_tx_commit",
		Params:  rawParamsJSON,
	}
	json.NewEncoder(requestBody).Encode(rc)
	port:=[3]string{"26657","36657","46657"}
	url := "http://localhost:"+port[0]
	req, err := http.NewRequest("POST", url, requestBody)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, _ := client.Do(req)

	//将结果定位到标准输出 也可以直接打印出来 或者定位到其他地方进行相应的处理
	stdout := os.Stdout
	_, err = io.Copy(stdout, response.Body)

	//返回的状态码
	status := response.StatusCode

	fmt.Println(status)

}

