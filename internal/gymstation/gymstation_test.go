package gymstation

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/desulaidovich/gym/internal/dto"
)

func TestGetGroupsSchedule(t *testing.T) {
	const sessionID = "16850c45cef0a1252e5832a7cc3166e3"

	tests := []struct {
		name    string
		status  int
		body    string
		want    []*dto.ScheduleItem
		wantErr bool
	}{
		{
			name:   "success",
			status: http.StatusOK,
			body: `{
				"error": false,
				"data": [
					{
						"id": "531c9251-b443-11f1-90a3-d85ed35e8660",
						"date": "2026-09-27T10:00:00+03:00",
						"title": "Функциональная тренировка АП",
						"club": "Клуб GYM STATION (Алые Паруса, Октябрьский пр-т, 117)",
						"direction": "Активные классы АП",
						"description": "<span>Описание</span>",
						"coach": "Астифеева Екатерина",
						"level": "-",
						"picture": "/upload/onecfitness/pictures/19e4c51d0319dc5f764a5a25a2f4f4c5.jpg"
					}
				]
			}`,
			want: []*dto.ScheduleItem{
				{
					ID:        "531c9251-b443-11f1-90a3-d85ed35e8660",
					Title:     "Функциональная тренировка АП",
					Direction: "Активные классы АП",
					Coach:     "Астифеева Екатерина",
					Level:     "-",
				},
			},
		},
		{
			name:   "empty data",
			status: http.StatusOK,
			body:   `{"error": false, "data": []}`,
			want:   []*dto.ScheduleItem{},
		},
		{
			name:    "api error",
			status:  http.StatusOK,
			body:    `{"error": true, "message": "invalid session"}`,
			wantErr: true,
		},
		{
			name:    "bad status",
			status:  http.StatusInternalServerError,
			body:    `{"error": false}`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			status:  http.StatusOK,
			body:    `not-a-json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				host: "gym-s.ru",
				httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, "https", r.URL.Scheme)
					assert.Equal(t, "gym-s.ru", r.URL.Host)
					assert.Equal(t, "/local/php_interface/api", r.URL.Path)
					assert.Equal(t, sessionID, r.URL.Query().Get("sessid"))
					assert.Equal(t, "getGroupsSchedule", r.URL.Query().Get("method"))
					assert.Equal(t, http.MethodGet, r.Method)

					return &http.Response{
						StatusCode: tt.status,
						Body:       io.NopCloser(strings.NewReader(tt.body)),
						Header:     make(http.Header),
						Request:    r,
					}, nil
				})},
			}

			got, err := client.GetGroupsSchedule(t.Context(), sessionID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, len(tt.want))

			for i := range tt.want {
				assert.Equal(t, tt.want[i].ID, got[i].ID)
				assert.Equal(t, tt.want[i].Title, got[i].Title)
				assert.Equal(t, tt.want[i].Direction, got[i].Direction)
				assert.Equal(t, tt.want[i].Coach, got[i].Coach)
				assert.Equal(t, tt.want[i].Level, got[i].Level)
			}
		})
	}
}

func TestGetSessionID(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		want    string
		wantErr bool
	}{
		{
			name:   "input",
			status: http.StatusOK,
			body:   `<input type="hidden" name="sessid" id="sessid" value="16850c45cef0a1252e5832a7cc3166e3" />`,
			want:   "16850c45cef0a1252e5832a7cc3166e3",
		},
		{
			name:   "script",
			status: http.StatusOK,
			body: "<script>window.BX=window.BX||{bitrix_sessid:function(){" +
				"return '16850c45cef0a1252e5832a7cc3166e3'}}</script>",
			want: "16850c45cef0a1252e5832a7cc3166e3",
		},
		{
			name:    "missing",
			status:  http.StatusOK,
			body:    `<html></html>`,
			wantErr: true,
		},
		{
			name:    "not ok",
			status:  http.StatusInternalServerError,
			body:    `<html><body><input name="sessid" value="abc123"></body></html>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				host:      "gym-s.ru",
				userAgent: "GymStation-Go-Client/1.0 (SorryForTheLoad)",
				reBitrix:  regexp.MustCompile(`bitrix_sessid:function\(\)\{return '([^']+)'\}`),
				httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, "https://gym-s.ru/", r.URL.String())
					assert.Equal(t, "GymStation-Go-Client/1.0 (SorryForTheLoad)", r.Header.Get("User-Agent"))
					return &http.Response{
						StatusCode: tt.status,
						Body:       io.NopCloser(strings.NewReader(tt.body)),
						Header:     make(http.Header),
						Request:    r,
					}, nil
				})},
			}

			got, err := client.GetSessionID(t.Context())
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
