package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rooobot/zrpc/client"
	"github.com/rooobot/zrpc/log"

	hw "github.com/rooobot/zrpc/examples/helloworld2/helloworld"
)

var concurrency = flag.Int64("concurrency", 500, "concurrency")
var target = flag.String("target", "127.0.0.1:8000", "target")
var total = flag.Int64("total", 1000000, "total requests")

func main() {
	flag.Parse()
	request(*total, *concurrency, *target)
}

func request(totalReqs int64, concurrency int64, target string) {

	perClientReqs := totalReqs / concurrency

	counter := &Counter{
		Total:       perClientReqs * concurrency,
		Concurrency: concurrency,
	}

	opts := []client.Option{
		client.WithTarget("127.0.0.1:8000"),
		client.WithNetwork("tcp"),
		client.WithTimeout(2000 * time.Millisecond),
		client.WithProtocol("proto"),
	}
	proxy := hw.NewGreeterClientProxy(opts...)
	req := &hw.HelloRequest{
		Msg: "hello",
	}

	var wg sync.WaitGroup
	wg.Add(int(concurrency))

	startTime := time.Now().UnixNano()

	for i := int64(0); i < counter.Concurrency; i++ {

		go func(i int64) {
			for j := int64(0); j < perClientReqs; j++ {

				rsp, err := proxy.SayHello(context.Background(), req, opts...)

				if err == nil && rsp.Msg == "world" {
					atomic.AddInt64(&counter.Succ, 1)
				} else {
					log.Info("rsp fail : %v", err)
					atomic.AddInt64(&counter.Fail, 1)
				}
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	counter.Cost = (time.Now().UnixNano() - startTime) / 1000000

	fmt.Printf("took %d ms for %d requests \n", counter.Cost, counter.Total)
	fmt.Printf("sent     requests      : %d\n", counter.Total)
	fmt.Printf("received requests      : %d\n", atomic.LoadInt64(&counter.Succ)+atomic.LoadInt64(&counter.Fail))
	fmt.Printf("received requests succ : %d\n", atomic.LoadInt64(&counter.Succ))
	fmt.Printf("received requests fail : %d\n", atomic.LoadInt64(&counter.Fail))
	fmt.Printf("throughput  (TPS)      : %d\n", totalReqs*1000/counter.Cost)

}

type Counter struct {
	Succ        int64 // 成功量
	Fail        int64 // 失败量
	Total       int64 // 总量
	Concurrency int64 // 并发量
	Cost        int64 // 总耗时 ms
}
