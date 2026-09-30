package job

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/desulaidovich/gym/internal/dto"
)

type HTTPClient interface {
	GetGroupsSchedule(ctx context.Context, sessionID string) ([]*dto.ScheduleItem, error)
	GetSessionID(ctx context.Context) (string, error)
}

type GymStationSiteParser struct {
	client HTTPClient
	log    *zap.Logger
}

func NewGymStationSiteParser(opts ...Option) (*GymStationSiteParser, error) {
	p := &GymStationSiteParser{}
	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("job: nil option")
		}
		if err := opt.apply(p); err != nil {
			return nil, fmt.Errorf("job: %w", err)
		}
	}
	if p.client == nil {
		return nil, fmt.Errorf("job: nil client")
	}
	if p.log == nil {
		return nil, fmt.Errorf("job: nil logger")
	}
	return p, nil
}

func (p *GymStationSiteParser) Run(ctx context.Context) error {
	id, err := p.client.GetSessionID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session ID: %w", err)
	}

	schedules, err := p.client.GetGroupsSchedule(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get groups schedule: %w", err)
	}

	p.log.Info("get schedules",
		zap.Int("schedules", len(schedules)),
	)
	return nil
}
