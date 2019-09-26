package Line

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"

)

type Line struct {
	target map[string][]string
	conns  map[string][]*websocket.Conn
}

//连接函数
func (l *Line) connect(host string) (*websocket.Conn, *http.Response, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: "/websocket"}
	return websocket.DefaultDialer.Dial(u.String(), nil)
}

//产生新的连接类型
func NewLine(target map[string][]string) *Line {
	//sum是算整体网络的节点个数，为了开辟相当的空间
	var sum int
	sum = 0
	for shard := range target {
		sum += len(target[shard])
	}
	return &Line{
		target: target, //目标节点地址
		conns:  make(map[string][]*websocket.Conn, sum), //连接地址
	}
}
//func (l *Line)Connect_target(key string)*websocket.Conn{
//	e:=useetcd.NewEtcd()
//	ip := string(e.Query(key))
//	l.target[key][0]=ip
//	c,_,err:=l.connect(ip)
//	if err != nil {
//		return nil
//	}
//	return c
//}
//建立连接数组
func(l *Line) Test(){
	fmt.Println(l.target["A"])
}
func (l *Line) Start() error {
	time.Sleep(time.Second*20)
	for shard := range l.target {
		l.conns[shard]=make([]*websocket.Conn,len(l.target[shard]))

		for i, ip := range l.target[shard] {

			c, _, err := l.connect(ip)
			if err != nil {
				fmt.Println("连接出错:",ip)
				return err
			}
			fmt.Println("首次连接",l.conns[shard][i])
			l.conns[shard][i] = c
		}
	}
	return nil
}

//发送消息，随机取一个连接给目标节点发送信息
func (l *Line) SendMessageTrans(message json.RawMessage, Receiver string,Sender string) error {
	rc:=&rpctypes.RPCRequest{
		JSONRPC:  "2.0",
		Sender:   Sender,
		Receiver: Receiver,
		Flag:     0,
		ID:       rpctypes.JSONRPCStringID("trans"),
		Method:   "broadcast_tx_commit",
		Params:   message,
	}
	rand.Seed(time.Now().Unix())
	//rnd := rand.Intn(4)
	c := l.conns[Receiver][0]
	err := c.WriteJSON(rc)
	if err != nil {
		return err
	}
	return nil
}
func (l *Line) SendMessageCommit(message json.RawMessage, Receiver string,Sender string) error {
	rc:=&rpctypes.RPCRequest{
		JSONRPC:  "2.0",
		Sender:   Sender,
		Receiver: Receiver,
		Flag:     0,
		ID:       rpctypes.JSONRPCStringID("commit"),
		Method:   "broadcast_tx_commit",
		Params:   message,
	}
	rand.Seed(time.Now().Unix())
	rnd := rand.Intn(4)
	c := l.conns[Receiver][rnd]
	err := c.WriteJSON(rc)
	if err != nil {
		return err
	}
	return nil
}
func (l *Line) ReceiveMessage(key string,connindex int)error{
	c:=l.conns[key][connindex]
	for {
		//第二个下划线指的是返回的信息，在下一步进行使用
		_, _, err := c.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return err
			}
			return nil
		}

		//if t.stopped || t.connsBroken[connIndex] {
		//	return
		//}
	}
}

