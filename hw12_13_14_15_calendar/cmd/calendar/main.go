package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	"github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/logger"
	//nolint:depguard
	internalhttp "github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/server/http"
	//nolint:depguard
	memorystorage "github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	//nolint:depguard
	sqlstorage "github.com/arelavvvv/vivakhnyuk_otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	logg := logger.New(config.Logger.LogFile, config.Logger.LogLevel)

	var storage app.Storage
	if config.Storage.StorageType == "in-memory" {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New(
			config.Database.DBUsername,
			config.Database.DBPassword,
			config.Database.DBHost,
			config.Database.DBPort,
			config.Database.DBDatabase,
		)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(*calendar, config.Server.HTTPHost, config.Server.HTTPPort)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
