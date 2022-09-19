package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"

	"github.com/cyril-jump/shortener/internal/app/config"
	"github.com/cyril-jump/shortener/internal/app/dto"
	servergrpc "github.com/cyril-jump/shortener/internal/app/server-grpc"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/storage/postgres"
	"github.com/cyril-jump/shortener/internal/app/storage/ram"
	"github.com/cyril-jump/shortener/internal/app/storage/rom"
	"github.com/cyril-jump/shortener/internal/app/storage/users"
	"github.com/cyril-jump/shortener/internal/app/utils"
	"github.com/cyril-jump/shortener/internal/app/workerpool"
)

func init() {
	// it outputs a message to stdout
	printAssemblyData()
	// evn vars

	var cfg config.EnvVar
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// flags
	flag.StringVarP(&config.Flags.ServerAddress, "address", "a", cfg.ServerAddress, "server address")
	flag.StringVarP(&config.Flags.BaseURL, "base", "b", cfg.BaseURL, "base url")
	flag.StringVarP(&config.Flags.FileStoragePath, "file", "f", cfg.FileStoragePath, "file storage path")
	flag.StringVarP(&config.Flags.DatabaseDSN, "psqlConn", "d", cfg.DatabaseDSN, "database URL conn")
	flag.StringVarP(&config.Flags.ConfigJSON, "json", "c", cfg.ConfigJSON, "JSON configuration")
	flag.StringVarP(&config.Flags.TrustedSubnet, "trusted", "t", cfg.TrustedSubnet, "trusted subnet")
	flag.Parse()

}

func main() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//db
	var db storage.DB

	//config
	cfg := config.NewConfig(
		config.Flags.ServerAddress,
		config.Flags.BaseURL,
		config.Flags.FileStoragePath,
		config.Flags.DatabaseDSN,
		config.Flags.ConfigJSON,
		config.Flags.EnableHTTPS,
		config.Flags.TrustedSubnet,
	)

	psqlConn, err := cfg.Get("database_dsn_str")
	utils.CheckErr(err, "database_dsn_str")

	fileStoragePath, err := cfg.Get("file_storage_path_str")
	utils.CheckErr(err, "file_storage_path_str")

	if fileStoragePath != "" {
		db, err = rom.NewDB(ctx, fileStoragePath)
		utils.CheckErr(err, "")
	} else if psqlConn != "" {
		db = postgres.New(ctx, psqlConn)
	} else {
		db = ram.NewDB(ctx)
	}
	usr := users.New(ctx)

	// Init Workers
	g, _ := errgroup.WithContext(ctx)
	recordCh := make(chan dto.Task, 50)
	doneCh := make(chan struct{})
	mu := &sync.Mutex{}

	inWorker := workerpool.NewInputWorker(recordCh, doneCh, ctx, mu)
	for i := 1; i <= runtime.NumCPU(); i++ {
		outWorker := workerpool.NewOutputWorker(i, recordCh, doneCh, ctx, db, mu)
		g.Go(outWorker.Do)
	}

	g.Go(inWorker.Loop)

	// Init GRPCServer

	grpcServer := servergrpc.InitSrv(db, cfg, usr, inWorker)

	// Start Server

	serverAddress, err := cfg.Get("server_address_str")
	utils.CheckErr(err, "server_address_str")

	grpcListener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("unable to listen on %v addr for grpc: %v", serverAddress, err)
	}

	go func() {
		err = grpcServer.Serve(grpcListener)
		if err != nil {
			log.Fatalf("grpc serve failed: %v", err)
		}
	}()

	<-signalChan

	log.Println("Shutting down...")

	cancel()

	grpcServer.GracefulStop()

	if err = db.Close(); err != nil {
		log.Fatal(err)
	}

	close(recordCh)
	close(doneCh)
	err = g.Wait()
	if err != nil {
		log.Println(err)
	}
}
