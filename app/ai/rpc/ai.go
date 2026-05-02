package main

import (
	"flag"
	"fmt"

	"amigo-api/app/ai/rpc/internal/config"
	"amigo-api/app/ai/rpc/internal/server"
	"amigo-api/app/ai/rpc/internal/svc"
	"amigo-api/common/mqueue"
	"amigo-api/common/pb"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/ai.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	err := mqueue.InitGlobalMQueue(&redis.Options{
		Addr:     c.AiRedis.Host,
		Password: c.AiRedis.Pass,
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
		pb.RegisterAiRpcServer(grpcServer, server.NewAiRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
