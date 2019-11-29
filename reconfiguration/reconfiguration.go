package reconfiguration

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tendermint/tendermint/SGX"
	tmrpc"github.com/tendermint/tendermint/rpc/client"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	cs "github.com/tendermint/tendermint/consensus"
)

type Reconfiguration struct {
	timestamp time.Time
	ShardPlan [][]int
	Nodesinfo [][]Nodeinfo
	ShardCount int
	NodeCount int
	MoveCount int
	Txs    [][]Tx
	IsLeader bool
	cs      *cs.ConsensusState
}
const(
	sendTimeout =time.Second*10
	)
type Tx []byte
type BlockMetas struct {
	Height int

}
type Nodeinfo struct {
	PeerId string
	ChainId string
}
var Ticker = time.NewTicker(time.Second*10)
func init(){
	go PeriodReconfiguration()
}

func NewReconfiguration()*Reconfiguration{
	return &Reconfiguration{
		ShardCount:4,
		NodeCount:16,
		MoveCount:10,
		IsLeader:false,
	}
}
func PeriodReconfiguration(){
	//定时器实现
	Re:=NewReconfiguration()
	for{
		select {
		case <- Ticker.C:
			if Re.IsLeader{
				client        := tmrpc.NewHTTP("Leader:26657", "/websocket")
				BlockTimeStamp:=Re.LatestBlockTime(client)
				Re.GenerateReconfiguration(BlockTimeStamp)
				Re.SendReconfiguration()
			}
		}
	}
}
func (Reconfiguration *Reconfiguration)LatestBlockTime(client tmrpc.Client) time.Time {
	//获取区块的可信时间戳
	status, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return status.SyncInfo.LatestBlockTime
}
func  (Re *Reconfiguration)GenerateReconfiguration(BlockStamp time.Time){
	//生成可信方案
	CredibleTimeStamp:=SGX.GetCredibleTimeStamp()
	if CredibleTimeStamp.Sub(BlockStamp)>1{
		//若满足时间间隔，则进行生成调整方案
		Re.GenerateCrediblePlan()

	}
}

func (Re *Reconfiguration)GenerateCrediblePlan(){
	for i:=0;i<Re.ShardCount;i++{
		Re.ShardPlan = append(Re.ShardPlan,Re.GenerateCredibleShardPlan())
	}
}
func Contain(ShardPlan []int,data int)(bool){
	for i:=0;i<len(ShardPlan);i++{
		if ShardPlan[i]==data{
			return true
		}
	}
	return false
}
func (Re *Reconfiguration)GenerateCredibleShardPlan()[]int{
	//生成单片调整计划
	var ShardPlan []int
	var Total int
	Total=0
	for{
		if(Total<Re.MoveCount){
			data:=SGX.GetCredibleRand(Re.NodeCount)
			if !Contain(ShardPlan,data){
				ShardPlan = append(ShardPlan, data)
				Total++
			}
		}else{
			break
		}
	}
	return ShardPlan
}
func (Re *Reconfiguration)SendReconfiguration(){

		Re.Nodesinfo = Re.GetNodeInfo()
		//发送交易到区块
		Re.Txs=make([][]Tx,Re.ShardCount)
		for i:=0;i<Re.ShardCount;i++{
			Re.Txs[i]=make([]Tx,Re.NodeCount)
			for j:=0;j<Re.NodeCount;j++{
				tx:= []byte("PeerId:"+Re.Nodesinfo[i][j].PeerId + " ChainId:"+Re.Nodesinfo[i][j].ChainId)
				Re.Txs[i][j]=tx
			}
		}
		for i:=0;i<Re.ShardCount;i++{
			go Re.SendMessage(Re.Txs[i],len(Re.Txs[i]))
		}

}
func Connect(host string) (*websocket.Conn, *http.Response, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: "/websocket"}
	return websocket.DefaultDialer.Dial(u.String(), nil)
}
func  (Re *Reconfiguration)SendMessage(Txs []Tx,num int){
	res, _ := json.Marshal(Txs)
	rawParamsJSON := json.RawMessage(res)
	//第一层打包结束
	c, _, err :=Connect("Leader:26657")
	if err!=nil{
		fmt.Fprintln(os.Stderr, err)
	}
	c.SetPingHandler(func(message string) error {
		err := c.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(sendTimeout))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	err1 := c.WriteJSON(rpctypes.RPCRequest{
		JSONRPC: "2.0",
		Sender:  "flag",
		ID:      rpctypes.JSONRPCStringID("relay"),
		Method:  "broadcast_tx_async",
		Params:  rawParamsJSON,
	})
	if err1 != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	time.Sleep(time.Second*5)
	c.Close()
}
func (Re *Reconfiguration)GetNodeInfo()[][]Nodeinfo{
	var Nodesinfo [][]Nodeinfo
	Nodesinfo = make([][]Nodeinfo,Re.ShardCount)
	for i:=0;i<Re.ShardCount;i++{

		for j:=0;j<Re.NodeCount;j++{
			Node:=Nodeinfo{
				PeerId:strconv.Itoa(j+i),
				ChainId:strconv.Itoa(j*i),
			}
			Nodesinfo[i] = append(Nodesinfo[i],Node)
		}
	}
	return Nodesinfo
}

