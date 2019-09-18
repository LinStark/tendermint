package useetcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Use_Etcd struct {

	Endpoints []string
}

func (e Use_Etcd)Update(key string,value string){
	cli, err := clientv3.New(clientv3.Config{
	Endpoints:e.Endpoints,
	DialKeepAliveTime:5*time.Second,
	})

	A:= "192.168.5.56"
	B:= "192.168.5.57"
	C:= "192.168.5.58"
	D:= "192.168.5.60"
	E:= "192.168.5.61"
	F:= "192.168.5.62"
	G:= "192.168.5.63"
	H:= "192.168.5.66"
	Value :=""
	if(key=="A"){
		Value=A+":"+value
		fmt.Println(Value)
	}
	if(key=="B"){
		Value=B+":"+value
		fmt.Println(Value)
	}
	if(key=="C"){
		Value=C+":"+value
		fmt.Println(Value)
	}
	if(key=="D"){
		Value=D+":"+value
		fmt.Println(Value)
	}
	if(key=="E"){
		Value=E+":"+value
		fmt.Println(Value)
	}
	if(key=="F"){
		Value=F+":"+value
		fmt.Println(Value)
	}
	if(key=="G"){
		Value=G+":"+value
		fmt.Println(Value)
	}
	if(key=="H"){
		Value=H+":"+value
		fmt.Println(Value)
	}
	key = "/"+key
	if err!=nil{
	fmt.Println("conn failure!")
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = cli.Put(ctx, key, Value)
	cancel()
	if err!=nil{
	fmt.Println("put failed!")
	return
	}
	fmt.Println("put success!")
	return
}


func (e Use_Etcd)Query(key string)(value []byte){
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:e.Endpoints,
		DialKeepAliveTime:5*time.Second,
	})
	key = "/"+key	
	fmt.Println(key)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Get(ctx,key)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		b:=ev.Value
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		return b
	}
	return
}
