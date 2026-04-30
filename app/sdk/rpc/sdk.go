package main

import (
	"flag"
	"fmt"

	"amigo-api/app/sdk/rpc/internal/config"
	"amigo-api/app/sdk/rpc/internal/server"
	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/mqueue"
	"amigo-api/common/pb"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/sdk.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化队列客户端
	err := mqueue.InitGlobalMQueue(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       0,
	}, &mqueue.QueueConfig{
		Queues:      map[string]int{"default": 6, "critical": 3, "low": 1},
		Concurrency: 10,
		MaxRetry:    3,
	})
	if err != nil {
		fmt.Println("init mqueue failed:", err)
	}

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSdkServer(grpcServer, server.NewSdkServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
