package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"rate_limiter/config"
	"rate_limiter/pb"
)

var conf, cmd, login, pass, ip string

func init() {
	flag.StringVar(&conf, "conf", "", "config file")
	flag.StringVar(&cmd, "cmd", "", "command to run")
	flag.StringVar(&login, "login", "", "login")
	flag.StringVar(&pass, "pass", "", "pass")
	flag.StringVar(&ip, "ip", "", "ip")
}

func main() {
	flag.Parse()
	if conf == "" {
		log.Fatal("conf is required")
	}
	if cmd == "" {
		log.Fatal("command is required")
	}

	cfg := config.GetEnvConfig(conf)
	grpcHost := fmt.Sprintf(":%d", cfg.Grpc.Port)

	conn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewRateLimiterServiceClient(conn)

	ctx := context.Background()

	switch cmd {
	case "check":
		err := check(ctx, client)
		if err != nil {
			log.Fatal(err)
		}
	case "reset":
		err := resetBucket(ctx, client)
		if err != nil {
			log.Fatal(err)
		}
	case "ab":
		err := addIPToBlackList(ctx, client)
		if err != nil {
			log.Fatal(err)
		}
	case "rb":
		err := removeIPFromBlackList(ctx, client)
		if err != nil {
			log.Fatal(err)
		}
	case "aw":
		err := addIPToWhiteList(ctx, client)
		if err != nil {
			log.Fatal(err)
		}
	case "rw":
		err := removeIPFromWhiteList(ctx, client)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func check(ctx context.Context, client pb.RateLimiterServiceClient) error {
	if login == "" {
		log.Fatal("login is required")
	}
	if pass == "" {
		log.Fatal("pass is required")
	}
	if ip == "" {
		log.Fatal("ip is required")
	}

	req := pb.LoginRequest{
		Login:    login,
		Password: pass,
		IP:       ip,
	}
	res, err := client.Allow(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%t", res.Ok)
	return err
}

func resetBucket(ctx context.Context, client pb.RateLimiterServiceClient) error {
	if login == "" {
		log.Fatal("login is required")
	}
	if ip == "" {
		log.Fatal("ip is required")
	}
	req := pb.ResetRequest{
		Login: login,
		IP:    ip,
	}
	_, err := client.ResetBucket(ctx, &req)

	return err
}

func addIPToBlackList(ctx context.Context, client pb.RateLimiterServiceClient) error {
	if ip == "" {
		log.Fatal("ip is required")
	}
	req := pb.IPRequest{
		IP: ip,
	}
	_, err := client.AddIPToBlackList(ctx, &req)

	return err
}

func removeIPFromBlackList(ctx context.Context, client pb.RateLimiterServiceClient) error {
	if ip == "" {
		log.Fatal("ip is required")
	}
	req := pb.IPRequest{
		IP: ip,
	}
	_, err := client.RemoveIPFromBlackList(ctx, &req)

	return err
}

func addIPToWhiteList(ctx context.Context, client pb.RateLimiterServiceClient) error {
	if ip == "" {
		log.Fatal("ip is required")
	}
	req := pb.IPRequest{
		IP: ip,
	}
	_, err := client.AddIPToWhiteList(ctx, &req)

	return err
}

func removeIPFromWhiteList(ctx context.Context, client pb.RateLimiterServiceClient) error {
	if ip == "" {
		log.Fatal("ip is required")
	}
	req := pb.IPRequest{
		IP: ip,
	}
	_, err := client.RemoveIPFromWhiteList(ctx, &req)

	return err
}
