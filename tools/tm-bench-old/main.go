package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log/term"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmrpc "github.com/tendermint/tendermint/rpc/client"
)

var logger = log.NewNopLogger()

func main() {
	//durationInt表示持续时间
	//txsRate 发送交易个数
	//connections表示链接的个数
	//txSize 表示tx的大小
	//verbose 冗余
	//outputFormat 输出格式
	//shard 分片
	//broadcastTxmethod 传播方法
	//allshard 全分片
	var durationInt, txsRate, connections, txSize int
	var verbose bool
	var outputFormat, broadcastTxMethod string
	//初始化数值
	flagSet := flag.NewFlagSet("tm-bench", flag.ExitOnError)
	flagSet.IntVar(&connections, "c", 1, "Connections to keep open per endpoint")
	flagSet.IntVar(&durationInt, "T", 10, "Exit after the specified amount of time in seconds")
	flagSet.IntVar(&txsRate, "r", 1000, "Txs per second to send in a connection")
	flagSet.IntVar(&txSize, "s", 40, "The size of a transaction in bytes, must be greater than or equal to 40.")
	flagSet.StringVar(&outputFormat, "output-format", "plain", "Output format: plain or json")
	flagSet.StringVar(&broadcastTxMethod, "broadcast-tx-method", "sync", "Broadcast method: async (no guarantees; fastest), sync (ensures tx is checked) or commit (ensures tx is checked and committed; slowest)")
	flagSet.BoolVar(&verbose, "v", false, "Verbose output")

	flagSet.Usage = func() {
		fmt.Println(`Tendermint blockchain benchmarking tool.

Usage:
	tm-bench [-c 1] [-T 10] [-r 1000] [-s 250] [endpoints] [-output-format <plain|json> [-broadcast-tx-method <async|sync|commit>]]

Examples:
	tm-bench localhost:26657`)
		flagSet.PrintDefaults()
	}

	flagSet.Parse(os.Args[1:]) //获取输入值

	if flagSet.NArg() == 0 {
		flagSet.Usage() //若没有设定值，使用默认值
		os.Exit(1)
	}

	if verbose {
		if outputFormat == "json" { //json格式发
			printErrorAndExit("Verbose mode not supported with json output.")
		}
		// Color errors red
		colorFn := func(keyvals ...interface{}) term.FgBgColor {
			for i := 1; i < len(keyvals); i += 2 {
				if _, ok := keyvals[i].(error); ok {
					return term.FgBgColor{Fg: term.White, Bg: term.Red}
				}
			}
			return term.FgBgColor{}
		}
		logger = log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)

		fmt.Printf("Running %ds test @ %s\n", durationInt, flagSet.Arg(0))
	}

	if txSize < 40 {
		printErrorAndExit("The size of a transaction must be greater than or equal to 40.")
	}

	if broadcastTxMethod != "async" &&
		broadcastTxMethod != "sync" &&
		broadcastTxMethod != "commit" {
		printErrorAndExit("broadcast-tx-method should be either 'sync', 'async' or 'commit'.")
	}

	var (
		endpoints     = strings.Split(flagSet.Arg(0), ",")        //可以对多个连接进行访问
		client        = tmrpc.NewHTTP(endpoints[0], "/websocket") //tmrpc.NewHTTP(连接地址，连接方式)
		initialHeight = latestBlockHeight(client)
	)

	logger.Info("Latest block height", "h", initialHeight)
	//设置交易器,初始化交易器
	transacters := startTransacters(
		endpoints,
		connections, //连接数量
		txsRate,
		txSize,
		"broadcast_tx_"+broadcastTxMethod,
	)

	// Stop upon receiving SIGTERM or CTRL-C.中断命令
	cmn.TrapSignal(logger, func() {
		for _, t := range transacters {
			t.Stop()
		}
	})

	// Wait until transacters have begun until we get the start time.
	timeStart := time.Now() //获取当前时间
	logger.Info("Time last transacter started", "t", timeStart)
	duration := time.Duration(durationInt) * time.Second //获取持续时间
	timeEnd := timeStart.Add(duration)                   //获取目标时间
	logger.Info("End time for calculation", "t", timeEnd)
	fmt.Println(timeEnd)
	<-time.After(duration) //等待时间
	//发送消息
	for i, t := range transacters {
		t.Stop()
		numCrashes := countCrashes(t.connsBroken)
		if numCrashes != 0 {
			fmt.Printf("%d connections crashed on transacter #%d\n", numCrashes, i)
		}
	}
	//结束消息
	logger.Debug("Time all transacters stopped", "t", time.Now())
	stats, err := calculateStatistics( //stats是统计数据
		client,
		initialHeight,
		timeStart,
		durationInt,
	)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	printStatistics(stats, outputFormat)
}

func latestBlockHeight(client tmrpc.Client) int64 {
	status, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return status.SyncInfo.LatestBlockHeight
}

func countCrashes(crashes []bool) int {
	count := 0
	for i := 0; i < len(crashes); i++ {
		if crashes[i] {
			count++
		}
	}
	return count
}

func startTransacters( //创建多线并行状态
	endpoints []string,
	connections,
	txsRate int,
	txSize int,
	broadcastTxMethod string, //这是参数列表
) []*transacter { //这是返回交易器
	transacters := make([]*transacter, len(endpoints)) //创建等价连接器

	wg := sync.WaitGroup{} //waitGroup的作用能够一直等到所有的goroutine执行完成
	wg.Add(len(endpoints)) //相当于线程个数
	for i, e := range endpoints {
		//e是目标地址 i是计数单位
		t := newTransacter(e, connections, txsRate, txSize, broadcastTxMethod)

		t.SetLogger(logger) //写日志
		go func(i int) {    //以并发的方式调用匿名函数func
			defer wg.Done() //开启一个t并行器
			if err := t.Start(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			transacters[i] = t
		}(i)
	}
	wg.Wait()

	return transacters
}

func printErrorAndExit(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
