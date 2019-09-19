package line
import (
	"fmt"
	"github.com/gorilla/websocket"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type line struct {
	target map[string] []string
	conns map[string][]*websocket.Conn
}
//连接函数
func (l *line)connect(host string) (*websocket.Conn, *http.Response, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: "/websocket"}
	return websocket.DefaultDialer.Dial(u.String(), nil)
}
//产生新的连接类型
func NewLine(target map[string] []string) *line {
	//sum是算整体网络的节点个数，为了开辟相当的空间
	var sum int
	sum=0
	for shard :=range target{
			sum+=len(target[shard])
	}
	return &line{
		target: target,//目标节点地址
		conns: make(map[string][]*websocket.Conn,sum),//连接地址
	}
}
//建立连接数组
func (l *line) start() error {
	for shard :=range l.target{
		fmt.Println(shard)
		for i,ip :=range l.target[shard]{
			c,_,err:=l.connect(ip)
			if err != nil {
				return err
			}
			l.conns[shard][i]=c
		}
	}
	return nil
}
//发送消息，随机取一个连接给目标节点发送信息
func (l *line) SendMessage(message rpctypes.RPCRequest,key string) error {
	rand.Seed(time.Now().Unix())
	rnd := rand.Intn(4)
	c := l.conns[key][rnd]
	err := c.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}
//func (l *line) ReceiveMessage(key string)error{
//
//}