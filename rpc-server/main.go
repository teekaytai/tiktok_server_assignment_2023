package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	rpc "github.com/teekaytai/tiktok_server_assignment_2023/rpc-server/kitex_gen/rpc/imservice"
)

func main() {
	ctx := context.Background()

	rdb := &RedisClient{}
	err := rdb.InitClient(ctx, "redis:6379", "")
	if err != nil {
		errMsg := fmt.Sprintf("failed to initialise Redis client, err: %v", err)
		log.Fatal(errMsg)
	}

	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	svr := rpc.NewServer(&IMServiceImpl{rdb}, server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
