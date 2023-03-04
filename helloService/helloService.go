package HelloService

import (
	"fmt"
	"log"
	"middleware/rpc_middleware/defs"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

var (
	helloServiceClient *HelloServiceClient
	helloServiceOnce   sync.Once
)

type HelloServiceClient struct {
	cli *rpc.Client
}

type HelloReq struct {
	Tar string
}

type HelloResp struct {
	Resp string
}

func mustGetJSONRPCClient(proto, host string) *rpc.Client {
	conn, err := net.Dial(proto, host)
	if err != nil {
		panic(fmt.Sprintf("Proto=%s,host=%s,err=%v", proto, host, err))
		return nil
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return client
}

func MustNewHelloServiceClient() *HelloServiceClient {
	helloServiceOnce.Do(func() {
		helloServiceClient = &HelloServiceClient{
			cli: mustGetJSONRPCClient(defs.connectProto, defs.HelloServiceHost),
		}
	})
	if helloServiceClient == nil {
		panic("get hello service client failed")
	}
	return helloServiceClient
}

func (c *HelloServiceClient) Hello(req *HelloReq) (resp *HelloResp, err error) {
	fmt.Println("get req=%s", req.Tar)
	resp = new(HelloResp)
	err = c.cli.Call("HelloService."+defs.HelloFunc, req, resp)
	return resp, err
}

type HelloServiceServer interface {
	Hello(req *HelloReq, resp *HelloResp) error
}

func HelloServiceStart(service HelloServiceServer) {
	rpc.Register(service)

	l, err := net.Listen(defs.connectProto, defs.LocalHost+":"+defs.HelloServicePost)
	if err != nil {
		panic(err)
	}

	go func() {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}()
}
