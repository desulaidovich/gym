package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/desulaidovich/gym/config"
	"github.com/desulaidovich/gym/internal/cron"
	"github.com/desulaidovich/gym/internal/gymstation"
	"github.com/desulaidovich/gym/internal/job"
	"github.com/desulaidovich/gym/internal/logger"
	"github.com/desulaidovich/gym/internal/server"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	runCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Parse("./config.yml")
	if err != nil {
		return err
	}

	log, err := logger.New(&cfg.Application)
	if err != nil {
		return err
	}
	defer func() {
		_ = log.Sync()
	}()

	client, err := gymstation.NewClient(&cfg.Gymstation)
	if err != nil {
		return err
	}

	cronLog := log.Named("cron")
	parser, err := job.NewGymStationSiteParser(
		job.WithClient(client),
		job.WithLogger(cronLog.Named("gym_station")),
	)
	if err != nil {
		return err
	}

	scheduler, err := cron.New(
		cron.WithLogger(cronLog),
		cron.WithConfig(&cfg.Cron),
		cron.WithJob("gymstation-site-parser", parser),
	)
	if err != nil {
		return err
	}

	debugServer, err := server.New(
		server.WithConfig(cfg.Server.Debug),
		server.WithHandler(http.NewServeMux()),
		server.WithLogger(log.Named("debug")),
	)
	if err != nil {
		return err
	}

	group, groupCtx := errgroup.WithContext(runCtx)

	group.Go(func() error {
		return scheduler.Run(groupCtx)
	})

	group.Go(func() error {
		return debugServer.Run(groupCtx)
	})

	return group.Wait()
}
