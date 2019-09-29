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
func newline() *Line{
	endpoints:=&node{
		target:make(map[string][]string,16),
	}

	endpoints.target["A"]=[]string{"192.168.5.56:26657","192.168.5.56:36657","192.168.5.56:46657","192.168.5.56:56657"}
	endpoints.target["B"]=[]string{"192.168.5.57:26657","192.168.5.57:36657","192.168.5.57:46657","192.168.5.57:56657"}
	endpoints.target["C"]=[]string{"192.168.5.58:26657","192.168.5.58:36657","192.168.5.58:46657","192.168.5.58:56657"}
	endpoints.target["D"]=[]string{"192.168.5.60:26657","192.168.5.60:36657","192.168.5.60:46657","192.168.5.60:56657"}

	l1:=NewLine(endpoints.target)

	return l1
}
type Line struct {
	target map[string][]string
	conns  map[string][]*websocket.Conn

}
//type Cn struct {
//	conn *websocket.Conn
//	mu sync.Mutex
//}
type node struct{
	target map[string] []string
}
var l *Line
//var cn1 *Cn
func begin(){
	if err:=l.Start();err!=nil{
		fmt.Println("连接错误！")
		fmt.Println(err)
		return
	}

}
func init(){

	l=newline()
	go begin()
}
func receiveloop(conn *websocket.Conn,shard string,i int){
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("连接中断")

			fmt.Println(err)
			ip := l.target[shard][i]
			fmt.Println("无法连接",ip)
			c,_,err := l.connect(ip)
			if err!=nil{
				fmt.Println("连接失败")
				fmt.Println(err)
				go receiveloop(c,shard,i)
				return
			}
			go receiveloop(c,shard,i)
			return
		}
		//fmt.Println(string(p))
	}
}

func UseConnect(key string,ip string)*websocket.Conn{
	//cn1.mu.Lock()
	//var j int
	rand.Seed(time.Now().Unix())
	rnd := rand.Intn(4)
	//for i,siteip := range l.target[key]{
	//	if siteip==ip{
	//		j = i
	//	}else{
	//		j = 0
	//	}
	//}
	//fmt.Println(j)
	//fmt.Println(l.target[key][j],"connection","j:",j)
	c:=l.conns[key][rnd]

	return c
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
			l.conns[shard][i] = c
			fmt.Println("连接成功！！")
			go receiveloop(c,shard,i)
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

