package cron

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/desulaidovich/gym/config"
)

type Option interface {
	apply(*Cron) error
}

type loggerOption struct {
	log *zap.Logger
}

func WithLogger(log *zap.Logger) Option {
	return loggerOption{
		log: log,
	}
}

func (o loggerOption) apply(c *Cron) error {
	if o.log == nil {
		return fmt.Errorf("nil logger")
	}
	c.log = o.log
	return nil
}

type configOption struct {
	cfg *config.Cron
}

func WithConfig(cfg *config.Cron) Option {
	return configOption{
		cfg: cfg,
	}
}

func (o configOption) apply(c *Cron) error {
	if o.cfg == nil {
		return fmt.Errorf("nil config")
	}
	c.cfg = o.cfg
	if o.cfg.StopTimeout > 0 {
		c.stopTimeout = o.cfg.StopTimeout
	}
	return nil
}

type jobOption struct {
	name string
	job  Job
}

func WithJob(name string, job Job) Option {
	return jobOption{
		name: name,
		job:  job,
	}
}

func (o jobOption) apply(c *Cron) error {
	if o.name == "" {
		return fmt.Errorf("empty job name")
	}
	if o.job == nil {
		return fmt.Errorf("nil job")
	}
	if _, ok := c.handlers[o.name]; ok {
		return fmt.Errorf("duplicate job %q", o.name)
	}
	c.handlers[o.name] = o.job
	return nil
}
