package main

import (
	"context"
	"github.com/rooobot/zrpc"
	hw "github.com/rooobot/zrpc/examples/helloworld2/helloworld"
	"time"
)

type greeterService struct{}

func (g *greeterService) SayHello(ctx context.Context, req *hw.HelloRequest) (*hw.HelloReply, error) {
	rsp := &hw.HelloReply{
		Msg: "world",
	}
	return rsp, nil
}

func main() {
	opts := []zrpc.ServerOption{
		zrpc.WithAddress("127.0.0.1:8000"),
		zrpc.WithNetwork("tcp"),
		zrpc.WithProtocol("proto"),
		zrpc.WithTimeout(time.Millisecond * 2000),
	}
	s := zrpc.NewServer(opts ...)
	hw.RegisterService(s, &greeterService{})

	pprof()
	s.Serve()
}
