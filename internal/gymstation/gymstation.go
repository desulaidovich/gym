package gymstation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"github.com/desulaidovich/gym/config"
	"github.com/desulaidovich/gym/internal/dto"
)

type Client struct {
	host       string
	userAgent  string
	httpClient *http.Client
	reBitrix   *regexp.Regexp
}

func NewClient(cfg *config.Gymstation) (*Client, error) {
	re := regexp.MustCompile(`bitrix_sessid:function\(\)\{return '([^']+)'\}`)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %w", err)
	}

	return &Client{
		host:      cfg.Host,
		userAgent: cfg.UserAgent,
		reBitrix:  re,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Jar:     jar,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   cfg.DialTimeout,
					KeepAlive: cfg.KeepAlive,
				}).DialContext,
				TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
				ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
				MaxIdleConns:          cfg.MaxIdleConns,
				IdleConnTimeout:       cfg.IdleConnTimeout,
				ExpectContinueTimeout: cfg.ExpectContinueTimeout,
				ForceAttemptHTTP2:     cfg.ForceAttemptHTTP2,
			},
		},
	}, nil
}

func (c *Client) GetGroupsSchedule(ctx context.Context, sessionID string) ([]*dto.ScheduleItem, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   "/local/php_interface/api",
	}

	q := u.Query()
	q.Set("sessid", sessionID)
	q.Set("method", "getGroupsSchedule")
	u.RawQuery = q.Encode()

	resp, err := c.do(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("getting groups schedule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting groups schedule: bad status: %s", resp.Status)
	}

	var data dto.ScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("getting groups schedule: %w", err)
	}

	if data.Error {
		return nil, errors.New("groups schedule error: " + data.Message)
	}

	return data.Data, nil
}

func (c *Client) GetSessionID(ctx context.Context) (string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   "/",
	}

	resp, err := c.do(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("getting session ID: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getting session ID: bad status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parsing session ID response: %w", err)
	}

	if val, exists := doc.Find(`input[name="sessid"]`).First().Attr("value"); exists {
		return val, nil
	}

	html, err := doc.Html()
	if err != nil {
		return "", fmt.Errorf("getting session ID HTML: %w", err)
	}

	if m := c.reBitrix.FindStringSubmatch(html); len(m) > 1 {
		return m[1], nil
	}

	return "", errors.New("session ID not found")
}

func (c *Client) do(ctx context.Context, method, rawURL string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}
