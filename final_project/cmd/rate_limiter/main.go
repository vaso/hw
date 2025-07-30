package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"rate_limiter/config"
	"rate_limiter/internal/app/service"
	"rate_limiter/pb"
)

var conf string

func init() {
	flag.StringVar(&conf, "conf", "", "config file")
}

func main() {
	flag.Parse()
	if conf == "" {
		log.Fatal("conf is required")
	}

	envConfig := config.GetEnvConfig(conf)
	fmt.Printf("envConfig: %+v\n", envConfig)
	appConfig := *config.GetAppConfig(*envConfig)
	fmt.Printf("appConfig: %+v\n", appConfig)
	ctx := context.Background()
	appService, err := service.NewService(ctx, appConfig)
	if err != nil {
		log.Fatal(err)
	}

	grpcHost := fmt.Sprintf(":%d", envConfig.Grpc.Port)
	var lc net.ListenConfig
	lsn, err := lc.Listen(ctx, "tcp", grpcHost)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterRateLimiterServiceServer(server, appService)

	log.Printf("starting server on %s", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
