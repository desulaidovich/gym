package cron

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	"github.com/desulaidovich/gym/config"
)

const DefaultShutdownTimeout = 15 * time.Second

type Cron struct {
	scheduler   gocron.Scheduler
	cfg         *config.Cron
	handlers    map[string]Job
	stopTimeout time.Duration
	log         *zap.Logger
}

func New(opts ...Option) (*Cron, error) {
	c := &Cron{
		handlers:    make(map[string]Job),
		stopTimeout: DefaultShutdownTimeout,
	}
	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("cron: nil option")
		}
		if err := opt.apply(c); err != nil {
			return nil, fmt.Errorf("cron: %w", err)
		}
	}
	if c.log == nil {
		return nil, fmt.Errorf("cron: nil logger")
	}
	if c.cfg == nil {
		return nil, fmt.Errorf("cron: nil config")
	}
	if err := c.build(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cron) Run(ctx context.Context) error {
	names := make([]string, 0, len(c.scheduler.Jobs()))
	for _, j := range c.scheduler.Jobs() {
		names = append(names, j.Name())
	}

	c.log.Info("started", zap.Strings("jobs", names))
	c.scheduler.Start()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), c.stopTimeout)
	defer cancel()

	if err := c.scheduler.ShutdownWithContext(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown scheduler: %w", err)
	}

	c.log.Info("stopped")
	return nil
}

func (c *Cron) build() error {
	for _, job := range c.cfg.Jobs {
		if job == nil {
			return fmt.Errorf("cron job: nil")
		}
		if job.Name == "" {
			return fmt.Errorf("cron job: empty name")
		}
		if _, ok := c.handlers[job.Name]; !ok {
			return fmt.Errorf("cron job %q: handler not found", job.Name)
		}
	}

	schedulerOpts := []gocron.SchedulerOption{
		gocron.WithStopTimeout(c.stopTimeout),
	}
	if c.cfg.Location != "" {
		loc, err := time.LoadLocation(c.cfg.Location)
		if err != nil {
			return fmt.Errorf("cron location %q: %w", c.cfg.Location, err)
		}
		schedulerOpts = append(schedulerOpts, gocron.WithLocation(loc))
	}

	s, err := gocron.NewScheduler(schedulerOpts...)
	if err != nil {
		return fmt.Errorf("scheduler: %w", err)
	}

	for _, job := range c.cfg.Jobs {
		if err := addJob(c.log, s, job, c.handlers[job.Name]); err != nil {
			return err
		}
	}

	c.scheduler = s
	return nil
}

func addJob(log *zap.Logger, s gocron.Scheduler, cfg *config.CronJob, handler Job) error {
	definition, err := jobDefinition(cfg)
	if err != nil {
		return fmt.Errorf("cron job %q: %w", cfg.Name, err)
	}

	_, err = s.NewJob(
		definition,
		gocron.NewTask(runner{log: log, name: cfg.Name, job: handler}.run),
		gocron.WithName(cfg.Name),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		return fmt.Errorf("cron job %q: %w", cfg.Name, err)
	}
	return nil
}

type runner struct {
	log  *zap.Logger
	name string
	job  Job
}

func (r runner) run(ctx context.Context) error {
	if err := r.job.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		r.log.Error("job",
			zap.String("name", r.name),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func jobDefinition(cfg *config.CronJob) (gocron.JobDefinition, error) {
	if cfg.Mask == "" {
		return nil, fmt.Errorf("empty mask")
	}
	return gocron.CronJob(cfg.Mask, len(strings.Fields(cfg.Mask)) > 5), nil
}
