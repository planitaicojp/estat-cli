package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	cerrors "github.com/planitaicojp/estat-cli/internal/errors"
)

var UserAgent = "planitaicojp/estat-cli/dev"

const defaultTimeout = 30 * time.Second

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	AppID      string
	Lang       string
	Verbose    bool
}

func NewClient(baseURL, appID, lang string) *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		BaseURL:    baseURL,
		AppID:      appID,
		Lang:       lang,
	}
}

func (c *Client) Get(path string, params map[string]string, result any) error {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return fmt.Errorf("URL解析エラー: %w", err)
	}

	q := u.Query()
	q.Set("appId", c.AppID)
	q.Set("lang", c.Lang)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "GET %s\n", u.String())
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return &cerrors.NetworkError{Err: err}
	}
	defer func() { _ = resp.Body.Close() }()

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "HTTP %d\n", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return &cerrors.APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTPステータス %d", resp.StatusCode),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("レスポンスのデコードに失敗: %w", err)
		}
	}
	return nil
}

func (c *Client) Post(path string, params map[string]string, result any) error {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return fmt.Errorf("URL解析エラー: %w", err)
	}

	q := u.Query()
	q.Set("appId", c.AppID)
	q.Set("lang", c.Lang)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "POST %s\n", u.String())
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return &cerrors.NetworkError{Err: err}
	}
	defer func() { _ = resp.Body.Close() }()

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "HTTP %d\n", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return &cerrors.APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTPステータス %d", resp.StatusCode),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("レスポンスのデコードに失敗: %w", err)
		}
	}
	return nil
}
