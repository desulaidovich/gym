package gymstation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/desulaidovich/gym/config"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration")
	}

	ctx := t.Context()
	client, err := NewClient(&config.Gymstation{
		Host:                  "gym-s.ru",
		Timeout:               30 * time.Second,
		DialTimeout:           10 * time.Second,
		KeepAlive:             30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: time.Second,
		ForceAttemptHTTP2:     true,
		UserAgent:             "GymStation-Go-Client/1.0 (SorryForTheLoad)",
	})
	require.NoError(t, err)

	id, err := client.GetSessionID(ctx)
	assert.NoError(t, err)

	items, err := client.GetGroupsSchedule(ctx, id)
	assert.NoError(t, err)

	assert.NotEmpty(t, items)
}
