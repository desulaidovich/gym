package job

import (
	"fmt"

	"go.uber.org/zap"
)

type Option interface {
	apply(*GymStationSiteParser) error
}

type clientOption struct {
	client HTTPClient
}

func WithClient(client HTTPClient) Option {
	return clientOption{
		client: client,
	}
}

func (o clientOption) apply(p *GymStationSiteParser) error {
	if o.client == nil {
		return fmt.Errorf("nil client")
	}
	p.client = o.client
	return nil
}

type loggerOption struct {
	log *zap.Logger
}

func WithLogger(log *zap.Logger) Option {
	return loggerOption{
		log: log,
	}
}

func (o loggerOption) apply(p *GymStationSiteParser) error {
	if o.log == nil {
		return fmt.Errorf("nil logger")
	}
	p.log = o.log
	return nil
}
