package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vrnvgasu/metrics/internal/agent"
	pb "github.com/vrnvgasu/metrics/internal/proto"
	"github.com/vrnvgasu/metrics/pkg/crypto"
	"github.com/vrnvgasu/metrics/pkg/retry"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func info() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func run() error {
	info()
	cnf := parseFlags()

	bufSize := cnf.ReportInterval / cnf.PollInterval * 100

	agentClient := agent.NewAgent(retry.NewClient(nil), bufSize)

	if cnf.CryptoKey != "" {
		pub, err := crypto.LoadPublicKey(cnf.CryptoKey)
		if err != nil {
			return fmt.Errorf("could not load public key: %w", err)
		}
		agentClient.SetPublicKey(pub)
	}

	ip, err := agent.OutboundIP(cnf.Address)
	if err != nil {
		return fmt.Errorf("could not determine outbound IP: %w", err)
	}
	agentClient.SetRealIP(ip.String())

	if cnf.GRPCAddress != "" {
		conn, err := grpc.NewClient(cnf.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("could not connect to gRPC server: %w", err)
		}
		defer conn.Close()
		agentClient.SetGRPCClient(pb.NewMetricsClient(conn))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	errGroup, runtimeCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return agentClient.Collect(runtimeCtx, cnf)
	})
	errGroup.Go(func() error {
		return agentClient.SendMetrics(runtimeCtx, cnf)
	})

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}
