# estat-cli プロジェクト構造 実装計画

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** e-Stat API 3.0 を操作する Go CLI ツールのプロジェクト骨格を構築し、`estat search` コマンドで統計表検索が動作するところまでを実装する。

**Architecture:** conoha-cli のレイヤー分離パターン（cmd → cmdutil → api → model → output）を踏襲し、e-Stat API の単純さ（appId 認証のみ、読み取り中心）に合わせて軽量化。Cobra によるサブコマンド構成、YAML 設定 + 環境変数オーバーライド。

**Tech Stack:** Go 1.22+, spf13/cobra, gopkg.in/yaml.v3, text/tabwriter, encoding/json, encoding/csv

---

## File Map

| File | Responsibility |
|------|---------------|
| `main.go` | エントリポイント — `cmd.Execute()` 呼び出し |
| `cmd/root.go` | ルートコマンド、グローバルフラグ、サブコマンド登録 |
| `cmd/version.go` | バージョン表示コマンド |
| `cmd/completion.go` | シェル補完コマンド |
| `cmd/cmdutil/client.go` | Cobra コンテキストから API クライアント生成 |
| `cmd/cmdutil/format.go` | 出力フォーマット解決（フラグ > 環境変数 > config） |
| `cmd/search/search.go` | `estat search` コマンド実装 |
| `cmd/get/get.go` | `estat get` コマンド（スタブ） |
| `cmd/meta/meta.go` | `estat meta` コマンド（スタブ） |
| `cmd/dataset/dataset.go` | `estat dataset` コマンド（スタブ） |
| `cmd/catalog/catalog.go` | `estat catalog` コマンド（スタブ） |
| `internal/config/config.go` | YAML 設定ロード・保存 |
| `internal/config/env.go` | 環境変数定数 + `EnvOr` ヘルパー |
| `internal/errors/exitcodes.go` | 終了コード定数 |
| `internal/errors/errors.go` | カスタムエラー型（ConfigError, APIError, NetworkError, ValidationError） |
| `internal/api/client.go` | ベース HTTP クライアント + 共通リクエスト処理 |
| `internal/api/statslist.go` | getStatsList API 呼び出し |
| `internal/model/statslist.go` | 統計表レスポンスモデル |
| `internal/output/formatter.go` | Formatter インターフェース + ファクトリ |
| `internal/output/table.go` | テーブル出力 |
| `internal/output/json.go` | JSON 出力 |
| `internal/output/csv.go` | CSV 出力 |

---

### Task 1: Go モジュール初期化 + Git 初期化

**Files:**
- Create: `go.mod`
- Create: `.gitignore`
- Create: `main.go`

- [ ] **Step 1: Git リポジトリ初期化**

```bash
cd /root/dev/planitai/estat-cli
git init
```

- [ ] **Step 2: .gitignore 作成**

```gitignore
# Binaries
estat
*.exe

# IDE
.idea/
.vscode/
*.swp

# OS
.DS_Store

# Build
dist/
```

- [ ] **Step 3: Go モジュール初期化**

```bash
go mod init github.com/planitaicojp/estat-cli
```

- [ ] **Step 4: main.go 作成**

```go
package main

import "github.com/planitaicojp/estat-cli/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 5: cobra 依存追加**

```bash
go get github.com/spf13/cobra@latest
go get gopkg.in/yaml.v3@latest
```

- [ ] **Step 6: コミット**

```bash
git add .gitignore go.mod go.sum main.go
git commit -m "feat: Go モジュール初期化 + エントリポイント作成"
```

---

### Task 2: エラー型 + 終了コード

**Files:**
- Create: `internal/errors/exitcodes.go`
- Create: `internal/errors/errors.go`
- Test: `internal/errors/errors_test.go`

- [ ] **Step 1: テスト作成**

```go
// internal/errors/errors_test.go
package errors

import "testing"

func TestGetExitCode_nil(t *testing.T) {
	if got := GetExitCode(nil); got != ExitOK {
		t.Errorf("GetExitCode(nil) = %d, want %d", got, ExitOK)
	}
}

func TestGetExitCode_APIError(t *testing.T) {
	err := &APIError{StatusCode: 400, Message: "bad request"}
	if got := GetExitCode(err); got != ExitAPI {
		t.Errorf("GetExitCode(APIError) = %d, want %d", got, ExitAPI)
	}
}

func TestGetExitCode_ConfigError(t *testing.T) {
	err := &ConfigError{Message: "appId未設定"}
	if got := GetExitCode(err); got != ExitConfig {
		t.Errorf("GetExitCode(ConfigError) = %d, want %d", got, ExitConfig)
	}
}

func TestGetExitCode_NetworkError(t *testing.T) {
	err := &NetworkError{Err: nil}
	if got := GetExitCode(err); got != ExitNetwork {
		t.Errorf("GetExitCode(NetworkError) = %d, want %d", got, ExitNetwork)
	}
}

func TestGetExitCode_ValidationError(t *testing.T) {
	err := &ValidationError{Message: "invalid"}
	if got := GetExitCode(err); got != ExitValidation {
		t.Errorf("GetExitCode(ValidationError) = %d, want %d", got, ExitValidation)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 100, Code: "100", Message: "パラメータ不正"}
	got := err.Error()
	want := "APIエラー (コード 100): パラメータ不正"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
```

- [ ] **Step 2: テスト実行して失敗確認**

```bash
go test ./internal/errors/
```

Expected: FAIL (files don't exist yet)

- [ ] **Step 3: exitcodes.go 作成**

```go
// internal/errors/exitcodes.go
package errors

const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitConfig     = 2
	ExitAPI        = 3
	ExitNetwork    = 4
	ExitValidation = 5
)
```

- [ ] **Step 4: errors.go 作成**

```go
// internal/errors/errors.go
package errors

import "fmt"

// ExitCoder is implemented by errors that carry a process exit code.
type ExitCoder interface {
	ExitCode() int
}

// ConfigError represents a configuration problem.
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("設定エラー: %s", e.Message)
}

func (e *ConfigError) ExitCode() int {
	return ExitConfig
}

// APIError represents an error returned by the e-Stat API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIエラー (コード %s): %s", e.Code, e.Message)
}

func (e *APIError) ExitCode() int {
	return ExitAPI
}

// NetworkError wraps an underlying network-level error.
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("ネットワークエラー: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

func (e *NetworkError) ExitCode() int {
	return ExitNetwork
}

// ValidationError represents invalid user input.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("バリデーションエラー (%s): %s", e.Field, e.Message)
	}
	return fmt.Sprintf("バリデーションエラー: %s", e.Message)
}

func (e *ValidationError) ExitCode() int {
	return ExitValidation
}

// GetExitCode returns the exit code for the given error.
func GetExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}
```

- [ ] **Step 5: テスト実行して成功確認**

```bash
go test ./internal/errors/ -v
```

Expected: PASS (6 tests)

- [ ] **Step 6: コミット**

```bash
git add internal/errors/
git commit -m "feat: カスタムエラー型と終了コードを追加"
```

---

### Task 3: 設定システム

**Files:**
- Create: `internal/config/env.go`
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`

- [ ] **Step 1: テスト作成**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_noFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESTAT_CONFIG_DIR", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Format != DefaultFormat {
		t.Errorf("Format = %q, want %q", cfg.Format, DefaultFormat)
	}
	if cfg.Lang != DefaultLang {
		t.Errorf("Lang = %q, want %q", cfg.Lang, DefaultLang)
	}
}

func TestLoad_withFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESTAT_CONFIG_DIR", dir)

	content := []byte("app_id: test-id\nformat: json\nlang: E\n")
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), content, 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.AppID != "test-id" {
		t.Errorf("AppID = %q, want %q", cfg.AppID, "test-id")
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %q, want %q", cfg.Format, "json")
	}
	if cfg.Lang != "E" {
		t.Errorf("Lang = %q, want %q", cfg.Lang, "E")
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESTAT_CONFIG_DIR", dir)

	cfg := &Config{AppID: "save-test", Format: "csv", Lang: "J"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.AppID != "save-test" {
		t.Errorf("AppID = %q, want %q", loaded.AppID, "save-test")
	}
}

func TestEnvOr(t *testing.T) {
	t.Setenv("TEST_VAR", "hello")
	if got := EnvOr("TEST_VAR", "fallback"); got != "hello" {
		t.Errorf("EnvOr = %q, want %q", got, "hello")
	}
	if got := EnvOr("NONEXISTENT_VAR", "fallback"); got != "fallback" {
		t.Errorf("EnvOr = %q, want %q", got, "fallback")
	}
}
```

- [ ] **Step 2: テスト実行して失敗確認**

```bash
go test ./internal/config/
```

Expected: FAIL

- [ ] **Step 3: env.go 作成**

```go
// internal/config/env.go
package config

import "os"

const (
	EnvAppID     = "ESTAT_APP_ID"
	EnvFormat    = "ESTAT_FORMAT"
	EnvLang      = "ESTAT_LANG"
	EnvConfigDir = "ESTAT_CONFIG_DIR"
)

// EnvOr returns the environment variable value if set, otherwise the fallback.
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
```

- [ ] **Step 4: config.go 作成**

```go
// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFormat = "table"
	DefaultLang   = "J"
	configFile    = "config.yaml"
)

// Config represents the estat CLI configuration.
type Config struct {
	AppID  string `yaml:"app_id"`
	Format string `yaml:"format"`
	Lang   string `yaml:"lang"`
}

// ConfigDir returns the configuration directory path.
func ConfigDir() string {
	if d := os.Getenv(EnvConfigDir); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "estat")
}

// Load reads the configuration from disk.
// If the file does not exist, returns default values.
func Load() (*Config, error) {
	path := filepath.Join(ConfigDir(), configFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗: %w", err)
	}
	if cfg.Format == "" {
		cfg.Format = DefaultFormat
	}
	if cfg.Lang == "" {
		cfg.Lang = DefaultLang
	}
	return &cfg, nil
}

// Save writes the configuration to disk.
func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("設定のマーシャルに失敗: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, configFile), data, 0600)
}

func defaultConfig() *Config {
	return &Config{
		Format: DefaultFormat,
		Lang:   DefaultLang,
	}
}
```

- [ ] **Step 5: テスト実行して成功確認**

```bash
go test ./internal/config/ -v
```

Expected: PASS (4 tests)

- [ ] **Step 6: コミット**

```bash
git add internal/config/
git commit -m "feat: YAML設定システムと環境変数オーバーライドを追加"
```

---

### Task 4: 出力フォーマッター

**Files:**
- Create: `internal/output/formatter.go`
- Create: `internal/output/table.go`
- Create: `internal/output/json.go`
- Create: `internal/output/csv.go`
- Test: `internal/output/formatter_test.go`

- [ ] **Step 1: テスト作成**

```go
// internal/output/formatter_test.go
package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type testRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestTableFormatter(t *testing.T) {
	var buf bytes.Buffer
	rows := []testRow{
		{ID: "001", Name: "テスト1"},
		{ID: "002", Name: "テスト2"},
	}

	f := New("table")
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ID") {
		t.Errorf("table output missing header 'ID'")
	}
	if !strings.Contains(out, "001") {
		t.Errorf("table output missing value '001'")
	}
	if !strings.Contains(out, "テスト1") {
		t.Errorf("table output missing value 'テスト1'")
	}
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	rows := []testRow{{ID: "001", Name: "テスト"}}

	f := New("json")
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	var result []testRow
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if result[0].ID != "001" {
		t.Errorf("ID = %q, want %q", result[0].ID, "001")
	}
}

func TestCSVFormatter(t *testing.T) {
	var buf bytes.Buffer
	rows := []testRow{
		{ID: "001", Name: "テスト"},
	}

	f := New("csv")
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("CSV lines = %d, want 2 (header + 1 row)", len(lines))
	}
	if !strings.Contains(lines[0], "id") {
		t.Errorf("CSV header missing 'id': %q", lines[0])
	}
	if !strings.Contains(lines[1], "001") {
		t.Errorf("CSV row missing '001': %q", lines[1])
	}
}

func TestNew_default(t *testing.T) {
	f := New("unknown")
	if _, ok := f.(*TableFormatter); !ok {
		t.Errorf("New('unknown') should return TableFormatter")
	}
}

func TestTableFormatter_emptySlice(t *testing.T) {
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, []testRow{}); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("empty slice should produce no output, got %q", buf.String())
	}
}
```

- [ ] **Step 2: テスト実行して失敗確認**

```bash
go test ./internal/output/
```

Expected: FAIL

- [ ] **Step 3: formatter.go 作成**

```go
// internal/output/formatter.go
package output

import "io"

// Formatter formats and writes data to a writer.
type Formatter interface {
	Format(w io.Writer, data any) error
}

// New creates a formatter for the given format name.
func New(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "csv":
		return &CSVFormatter{}
	default:
		return &TableFormatter{}
	}
}
```

- [ ] **Step 4: table.go 作成**

```go
// internal/output/table.go
package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

// TableFormatter formats data as a human-readable table.
type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		_, err := fmt.Fprintf(w, "%v\n", data)
		return err
	}

	if val.Len() == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()

	headers := make([]string, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Name
		}
		headers[i] = strings.ToUpper(name)
	}
	if _, err := fmt.Fprintln(tw, strings.Join(headers, "\t")); err != nil {
		return err
	}

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}
		fields := make([]string, row.NumField())
		for j := 0; j < row.NumField(); j++ {
			fields[j] = fmt.Sprintf("%v", row.Field(j).Interface())
		}
		if _, err := fmt.Fprintln(tw, strings.Join(fields, "\t")); err != nil {
			return err
		}
	}

	return tw.Flush()
}
```

- [ ] **Step 5: json.go 作成**

```go
// internal/output/json.go
package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter formats data as pretty-printed JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}
```

- [ ] **Step 6: csv.go 作成**

```go
// internal/output/csv.go
package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

// CSVFormatter formats data as CSV with headers.
type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, data any) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		return fmt.Errorf("CSVフォーマッターにはスライスが必要です（%T が渡されました）", data)
	}
	if val.Len() == 0 {
		return nil
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()

	headers := make([]string, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Name
		}
		headers[i] = name
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}
		record := make([]string, row.NumField())
		for j := 0; j < row.NumField(); j++ {
			record[j] = fmt.Sprintf("%v", row.Field(j).Interface())
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
```

- [ ] **Step 7: テスト実行して成功確認**

```bash
go test ./internal/output/ -v
```

Expected: PASS (5 tests)

- [ ] **Step 8: コミット**

```bash
git add internal/output/
git commit -m "feat: 出力フォーマッター（table/json/csv）を追加"
```

---

### Task 5: API クライアント基盤

**Files:**
- Create: `internal/api/client.go`
- Test: `internal/api/client_test.go`

- [ ] **Step 1: テスト作成**

```go
// internal/api/client_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Get_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("appId") != "test-app-id" {
			t.Errorf("appId = %q, want %q", r.URL.Query().Get("appId"), "test-app-id")
		}
		if r.URL.Query().Get("lang") != "J" {
			t.Errorf("lang = %q, want %q", r.URL.Query().Get("lang"), "J")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-app-id", "J")
	var result map[string]string
	if err := client.Get("/json/getStatsList", nil, &result); err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("status = %q, want %q", result["status"], "ok")
	}
}

func TestClient_Get_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"GET_STATS_LIST": map[string]any{
				"RESULT": map[string]any{
					"STATUS":    100,
					"ERROR_MSG": "パラメータ不正",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad-id", "J")
	var result map[string]any
	err := client.Get("/json/getStatsList", nil, &result)
	// The raw response will be parsed by callers; the client just decodes JSON
	if err != nil {
		t.Fatalf("Get() should not return transport error: %v", err)
	}
}

func TestClient_Get_networkError(t *testing.T) {
	client := NewClient("http://localhost:1", "test", "J")
	var result map[string]string
	err := client.Get("/json/getStatsList", nil, &result)
	if err == nil {
		t.Fatal("Get() should return error for bad host")
	}
}
```

- [ ] **Step 2: テスト実行して失敗確認**

```bash
go test ./internal/api/
```

Expected: FAIL

- [ ] **Step 3: client.go 作成**

```go
// internal/api/client.go
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

// UserAgent is the User-Agent header value sent with all requests.
var UserAgent = "planitaicojp/estat-cli/dev"

const defaultTimeout = 30 * time.Second

// Client is the base HTTP client for e-Stat API.
type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	AppID      string
	Lang       string
	Verbose    bool
}

// NewClient creates a new e-Stat API client.
func NewClient(baseURL, appID, lang string) *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		BaseURL:    baseURL,
		AppID:      appID,
		Lang:       lang,
	}
}

// Get performs a GET request to the given path with query parameters.
// params can be nil. The response is decoded into result.
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
	defer resp.Body.Close()

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

// Post performs a POST request to the given path.
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
	defer resp.Body.Close()

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
```

- [ ] **Step 4: テスト実行して成功確認**

```bash
go test ./internal/api/ -v
```

Expected: PASS (3 tests)

- [ ] **Step 5: コミット**

```bash
git add internal/api/
git commit -m "feat: e-Stat APIベースHTTPクライアントを追加"
```

---

### Task 6: 統計表モデル + getStatsList API

**Files:**
- Create: `internal/model/statslist.go`
- Create: `internal/api/statslist.go`
- Test: `internal/api/statslist_test.go`

- [ ] **Step 1: テスト作成**

```go
// internal/api/statslist_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatsList(t *testing.T) {
	responseJSON := `{
		"GET_STATS_LIST": {
			"RESULT": {
				"STATUS": 0,
				"ERROR_MSG": "正常に終了しました。",
				"DATE": "2026-04-03T10:00:00.000+09:00"
			},
			"PARAMETER": {
				"LANG": "J",
				"SEARCH_WORD": "人口"
			},
			"DATALIST_INF": {
				"NUMBER": 1,
				"RESULT_INF": {
					"FROM_NUMBER": 1,
					"TO_NUMBER": 1
				},
				"TABLE_INF": [
					{
						"@id": "0003410379",
						"STAT_NAME": {"@code": "00200521", "$": "国勢調査"},
						"GOV_ORG": {"@code": "00200", "$": "総務省"},
						"STATISTICS_NAME": "国勢調査 人口等基本集計",
						"TITLE": {"@no": "001", "$": "男女別人口－全国，都道府県"},
						"SURVEY_DATE": "202001",
						"OPEN_DATE": "2021-11-30",
						"OVERALL_TOTAL_NUMBER": 12345,
						"UPDATED_DATE": "2022-01-15"
					}
				]
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-id", "J")
	resp, err := GetStatsList(client, map[string]string{"searchWord": "人口"})
	if err != nil {
		t.Fatalf("GetStatsList() error: %v", err)
	}

	if resp.Result.Status != 0 {
		t.Errorf("Status = %d, want 0", resp.Result.Status)
	}

	tables := resp.DatalistInf.TableInf
	if len(tables) != 1 {
		t.Fatalf("len(TableInf) = %d, want 1", len(tables))
	}

	table := tables[0]
	if table.ID != "0003410379" {
		t.Errorf("ID = %q, want %q", table.ID, "0003410379")
	}
	if table.StatName.Name != "国勢調査" {
		t.Errorf("StatName = %q, want %q", table.StatName.Name, "国勢調査")
	}
	if table.Title.Name != "男女別人口－全国，都道府県" {
		t.Errorf("Title = %q, want %q", table.Title.Name, "男女別人口－全国，都道府県")
	}
}

func TestGetStatsList_apiError(t *testing.T) {
	responseJSON := `{
		"GET_STATS_LIST": {
			"RESULT": {
				"STATUS": 100,
				"ERROR_MSG": "アプリケーションIDが不正です。",
				"DATE": "2026-04-03T10:00:00.000+09:00"
			},
			"PARAMETER": {}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(json.RawMessage(responseJSON))
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad-id", "J")
	_, err := GetStatsList(client, nil)
	if err == nil {
		t.Fatal("GetStatsList() should return error for STATUS != 0")
	}
}
```

- [ ] **Step 2: テスト実行して失敗確認**

```bash
go test ./internal/api/ -run TestGetStatsList
```

Expected: FAIL

- [ ] **Step 3: model/statslist.go 作成**

```go
// internal/model/statslist.go
package model

// StatsListResponse is the top-level response from getStatsList.
type StatsListResponse struct {
	GetStatsList struct {
		Result      ResultInf   `json:"RESULT"`
		Parameter   interface{} `json:"PARAMETER"`
		DatalistInf DatalistInf `json:"DATALIST_INF"`
	} `json:"GET_STATS_LIST"`
}

// ResultInf is the common result section in API responses.
type ResultInf struct {
	Status   int    `json:"STATUS"`
	ErrorMsg string `json:"ERROR_MSG"`
	Date     string `json:"DATE"`
}

// DatalistInf contains the list of statistical tables.
type DatalistInf struct {
	Number    int         `json:"NUMBER"`
	ResultInf PageInfo    `json:"RESULT_INF"`
	TableInf  []TableInfo `json:"TABLE_INF"`
}

// PageInfo contains pagination information.
type PageInfo struct {
	FromNumber int `json:"FROM_NUMBER"`
	ToNumber   int `json:"TO_NUMBER"`
	NextKey    int `json:"NEXT_KEY,omitempty"`
}

// TableInfo represents a single statistical table entry.
type TableInfo struct {
	ID                 string       `json:"@id"`
	StatName           CodeNamePair `json:"STAT_NAME"`
	GovOrg             CodeNamePair `json:"GOV_ORG"`
	StatisticsName     string       `json:"STATISTICS_NAME"`
	Title              NoNamePair   `json:"TITLE"`
	SurveyDate         string       `json:"SURVEY_DATE"`
	OpenDate           string       `json:"OPEN_DATE"`
	OverallTotalNumber int          `json:"OVERALL_TOTAL_NUMBER"`
	UpdatedDate        string       `json:"UPDATED_DATE"`
}

// CodeNamePair represents a JSON object with @code and $ fields.
type CodeNamePair struct {
	Code string `json:"@code"`
	Name string `json:"$"`
}

// NoNamePair represents a JSON object with @no and $ fields.
type NoNamePair struct {
	No   string `json:"@no"`
	Name string `json:"$"`
}

// TableRow is a flattened row for output formatting.
type TableRow struct {
	ID             string `json:"id"`
	StatName       string `json:"stat_name"`
	Title          string `json:"title"`
	SurveyDate     string `json:"survey_date"`
	OpenDate       string `json:"open_date"`
}

// ToRows converts a list of TableInfo to output-ready rows.
func ToTableRows(tables []TableInfo) []TableRow {
	rows := make([]TableRow, len(tables))
	for i, t := range tables {
		rows[i] = TableRow{
			ID:         t.ID,
			StatName:   t.StatName.Name,
			Title:      t.Title.Name,
			SurveyDate: t.SurveyDate,
			OpenDate:   t.OpenDate,
		}
	}
	return rows
}
```

- [ ] **Step 4: api/statslist.go 作成**

```go
// internal/api/statslist.go
package api

import (
	"fmt"

	cerrors "github.com/planitaicojp/estat-cli/internal/errors"
	"github.com/planitaicojp/estat-cli/internal/model"
)

// StatsListResult is the unwrapped response for command usage.
type StatsListResult struct {
	Result      model.ResultInf
	DatalistInf model.DatalistInf
}

// GetStatsList calls the getStatsList API endpoint.
func GetStatsList(c *Client, params map[string]string) (*StatsListResult, error) {
	var resp model.StatsListResponse
	if err := c.Get("/json/getStatsList", params, &resp); err != nil {
		return nil, err
	}

	inner := resp.GetStatsList
	if inner.Result.Status != 0 {
		return nil, &cerrors.APIError{
			StatusCode: inner.Result.Status,
			Code:       fmt.Sprintf("%d", inner.Result.Status),
			Message:    inner.Result.ErrorMsg,
		}
	}

	return &StatsListResult{
		Result:      inner.Result,
		DatalistInf: inner.DatalistInf,
	}, nil
}
```

- [ ] **Step 5: テスト実行して成功確認**

```bash
go test ./internal/api/ -v
go test ./internal/model/ -v
```

Expected: PASS

- [ ] **Step 6: コミット**

```bash
git add internal/model/statslist.go internal/api/statslist.go internal/api/statslist_test.go
git commit -m "feat: getStatsList APIクライアントとレスポンスモデルを追加"
```

---

### Task 7: ルートコマンド + cmdutil

**Files:**
- Create: `cmd/root.go`
- Create: `cmd/version.go`
- Create: `cmd/completion.go`
- Create: `cmd/cmdutil/client.go`
- Create: `cmd/cmdutil/format.go`

- [ ] **Step 1: cmd/cmdutil/format.go 作成**

```go
// cmd/cmdutil/format.go
package cmdutil

import (
	"github.com/spf13/cobra"

	"github.com/planitaicojp/estat-cli/internal/config"
)

// GetFormat resolves the output format from flag > env > config > default.
func GetFormat(cmd *cobra.Command) string {
	format, _ := cmd.Flags().GetString("format")
	if format != "" {
		return format
	}
	if f := config.EnvOr(config.EnvFormat, ""); f != "" {
		return f
	}
	cfg, err := config.Load()
	if err != nil {
		return config.DefaultFormat
	}
	return cfg.Format
}
```

- [ ] **Step 2: cmd/cmdutil/client.go 作成**

```go
// cmd/cmdutil/client.go
package cmdutil

import (
	"github.com/spf13/cobra"

	"github.com/planitaicojp/estat-cli/internal/api"
	"github.com/planitaicojp/estat-cli/internal/config"
	cerrors "github.com/planitaicojp/estat-cli/internal/errors"
)

// NewClient creates an API client from the cobra command context.
func NewClient(cmd *cobra.Command) (*api.Client, error) {
	appID, _ := cmd.Flags().GetString("app-id")
	if appID == "" {
		appID = config.EnvOr(config.EnvAppID, "")
	}
	if appID == "" {
		cfg, err := config.Load()
		if err != nil {
			return nil, err
		}
		appID = cfg.AppID
	}
	if appID == "" {
		return nil, &cerrors.ConfigError{
			Message: "アプリケーションIDが設定されていません。--app-id フラグ、ESTAT_APP_ID 環境変数、または ~/.config/estat/config.yaml で設定してください",
		}
	}

	lang, _ := cmd.Flags().GetString("lang")
	if lang == "" {
		lang = config.EnvOr(config.EnvLang, "")
	}
	if lang == "" {
		cfg, err := config.Load()
		if err != nil {
			lang = config.DefaultLang
		} else {
			lang = cfg.Lang
		}
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	client := api.NewClient("https://api.e-stat.go.jp/rest/3.0/app", appID, lang)
	client.Verbose = verbose
	return client, nil
}
```

- [ ] **Step 3: cmd/root.go 作成**

```go
// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/estat-cli/cmd/catalog"
	"github.com/planitaicojp/estat-cli/cmd/dataset"
	"github.com/planitaicojp/estat-cli/cmd/get"
	"github.com/planitaicojp/estat-cli/cmd/meta"
	"github.com/planitaicojp/estat-cli/cmd/search"
	cerrors "github.com/planitaicojp/estat-cli/internal/errors"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:           "estat",
	Short:         "e-Stat API CLI",
	Long:          "e-Stat（政府統計の総合窓口）APIを操作するコマンドラインツール",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().String("app-id", "", "アプリケーションID")
	rootCmd.PersistentFlags().String("format", "", "出力形式: table, json, csv")
	rootCmd.PersistentFlags().String("lang", "", "言語: J(日本語), E(英語)")
	rootCmd.PersistentFlags().Bool("no-color", false, "色出力を無効にする")
	rootCmd.PersistentFlags().Bool("verbose", false, "詳細出力")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(search.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(meta.Cmd)
	rootCmd.AddCommand(dataset.Cmd)
	rootCmd.AddCommand(catalog.Cmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}
```

- [ ] **Step 4: cmd/version.go 作成**

```go
// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("estat version %s\n", version)
	},
}
```

- [ ] **Step 5: cmd/completion.go 作成**

```go
// cmd/completion.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "シェル補完スクリプトを生成",
	Long: `指定したシェルの補完スクリプトを標準出力に出力します。

使用例:
  # bash
  estat completion bash > /etc/bash_completion.d/estat

  # zsh
  estat completion zsh > "${fpath[1]}/_estat"

  # fish
  estat completion fish > ~/.config/fish/completions/estat.fish
`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
```

- [ ] **Step 6: コミット**

```bash
git add cmd/root.go cmd/version.go cmd/completion.go cmd/cmdutil/
git commit -m "feat: ルートコマンド、バージョン、補完、cmdutilを追加"
```

---

### Task 8: サブコマンドスタブ（get, meta, dataset, catalog）

**Files:**
- Create: `cmd/get/get.go`
- Create: `cmd/meta/meta.go`
- Create: `cmd/dataset/dataset.go`
- Create: `cmd/catalog/catalog.go`

- [ ] **Step 1: cmd/get/get.go 作成**

```go
// cmd/get/get.go
package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd is the top-level get command.
var Cmd = &cobra.Command{
	Use:   "get <statsDataId>",
	Short: "統計データを取得",
	Long:  "指定した統計表IDの統計データを取得します。",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat get は今後のバージョンで実装予定です")
	},
}

func init() {
	Cmd.Flags().String("area", "", "地域コード")
	Cmd.Flags().String("time", "", "時間軸コード")
	Cmd.Flags().String("category", "", "分類項目フィルタ（例: cat01=A01）")
	Cmd.Flags().Bool("section-header", false, "セクションヘッダを含める")
	Cmd.Flags().Bool("bulk", false, "一括取得モード")
}
```

- [ ] **Step 2: cmd/meta/meta.go 作成**

```go
// cmd/meta/meta.go
package meta

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd is the top-level meta command.
var Cmd = &cobra.Command{
	Use:   "meta <statsDataId>",
	Short: "メタ情報を取得",
	Long:  "指定した統計表IDのメタ情報（分類項目、地域コードなど）を取得します。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat meta は今後のバージョンで実装予定です")
	},
}

func init() {
	Cmd.Flags().String("class", "", "特定の分類のみ取得（例: area, time）")
}
```

- [ ] **Step 3: cmd/dataset/dataset.go 作成**

```go
// cmd/dataset/dataset.go
package dataset

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd is the top-level dataset command.
var Cmd = &cobra.Command{
	Use:   "dataset",
	Short: "データセットを管理",
	Long:  "データセットの登録・参照を行います。",
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "データセットを登録",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat dataset register は今後のバージョンで実装予定です")
	},
}

var showCmd = &cobra.Command{
	Use:   "show <datasetId>",
	Short: "データセットを参照",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat dataset show は今後のバージョンで実装予定です")
	},
}

func init() {
	registerCmd.Flags().String("stats-data-id", "", "統計表ID")
	registerCmd.Flags().String("filter", "", "フィルタ条件")

	Cmd.AddCommand(registerCmd)
	Cmd.AddCommand(showCmd)
}
```

- [ ] **Step 4: cmd/catalog/catalog.go 作成**

```go
// cmd/catalog/catalog.go
package catalog

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd is the top-level catalog command.
var Cmd = &cobra.Command{
	Use:   "catalog",
	Short: "データカタログを取得",
	Long:  "統計データのカタログ情報を取得します。",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("未実装: estat catalog は今後のバージョンで実装予定です")
	},
}

func init() {
	Cmd.Flags().String("survey", "", "調査名で絞り込み")
	Cmd.Flags().String("field", "", "統計分野コードで絞り込み")
	Cmd.Flags().String("dataset-type", "", "データセット種別（db または file）")
}
```

- [ ] **Step 5: コミット**

```bash
git add cmd/get/ cmd/meta/ cmd/dataset/ cmd/catalog/
git commit -m "feat: サブコマンドスタブ（get, meta, dataset, catalog）を追加"
```

---

### Task 9: `estat search` コマンド実装

**Files:**
- Create: `cmd/search/search.go`
- Test: `cmd/search/search_test.go`

- [ ] **Step 1: テスト作成**

```go
// cmd/search/search_test.go
package search

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"
)

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"GET_STATS_LIST": map[string]any{
				"RESULT": map[string]any{
					"STATUS":    0,
					"ERROR_MSG": "正常に終了しました。",
					"DATE":      "2026-04-03T10:00:00.000+09:00",
				},
				"PARAMETER": map[string]any{},
				"DATALIST_INF": map[string]any{
					"NUMBER": 1,
					"RESULT_INF": map[string]any{
						"FROM_NUMBER": 1,
						"TO_NUMBER":   1,
					},
					"TABLE_INF": []map[string]any{
						{
							"@id":                   "0003410379",
							"STAT_NAME":             map[string]any{"@code": "00200521", "$": "国勢調査"},
							"GOV_ORG":               map[string]any{"@code": "00200", "$": "総務省"},
							"STATISTICS_NAME":        "国勢調査 人口等基本集計",
							"TITLE":                  map[string]any{"@no": "001", "$": "男女別人口"},
							"SURVEY_DATE":            "202001",
							"OPEN_DATE":              "2021-11-30",
							"OVERALL_TOTAL_NUMBER":   12345,
							"UPDATED_DATE":           "2022-01-15",
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestSearchCmd_json(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	var buf bytes.Buffer
	root := &cobra.Command{Use: "estat"}
	root.PersistentFlags().String("app-id", "", "")
	root.PersistentFlags().String("format", "", "")
	root.PersistentFlags().String("lang", "", "")
	root.PersistentFlags().Bool("verbose", false, "")
	root.AddCommand(Cmd)

	root.SetOut(&buf)
	root.SetArgs([]string{"search", "人口", "--app-id", "test", "--format", "json", "--base-url", server.URL})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	var result []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal error: %v\noutput: %s", err, buf.String())
	}
	if len(result) != 1 {
		t.Fatalf("len(result) = %d, want 1", len(result))
	}
	if result[0]["id"] != "0003410379" {
		t.Errorf("id = %v, want %q", result[0]["id"], "0003410379")
	}
}
```

- [ ] **Step 2: テスト実行して失敗確認**

```bash
go test ./cmd/search/ -v
```

Expected: FAIL

- [ ] **Step 3: cmd/search/search.go 作成**

```go
// cmd/search/search.go
package search

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/estat-cli/cmd/cmdutil"
	"github.com/planitaicojp/estat-cli/internal/api"
	"github.com/planitaicojp/estat-cli/internal/model"
	"github.com/planitaicojp/estat-cli/internal/output"
)

// Cmd is the top-level search command.
var Cmd = &cobra.Command{
	Use:   "search [キーワード]",
	Short: "統計表を検索",
	Long: `e-Stat APIで統計表を検索します。

使用例:
  estat search 人口
  estat search --survey "国勢調査"
  estat search --field "02" --open-year 2020
  estat search 人口 --limit 5 --format json`,
	RunE: runSearch,
}

func init() {
	Cmd.Flags().String("survey", "", "調査名で絞り込み")
	Cmd.Flags().String("field", "", "統計分野コードで絞り込み")
	Cmd.Flags().String("open-year", "", "公開年で絞り込み")
	Cmd.Flags().String("stats-code", "", "政府統計コードで絞り込み")
	Cmd.Flags().Int("limit", 0, "取得件数の上限")
	Cmd.Flags().Int("start", 0, "取得開始位置")
	// Hidden flag for testing
	Cmd.Flags().String("base-url", "", "APIベースURL（テスト用）")
	Cmd.Flags().MarkHidden("base-url")
}

func runSearch(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient(cmd)
	if err != nil {
		return err
	}

	// Override base URL for testing
	if baseURL, _ := cmd.Flags().GetString("base-url"); baseURL != "" {
		client.BaseURL = baseURL
	}

	params := buildParams(cmd, args)

	result, err := api.GetStatsList(client, params)
	if err != nil {
		return err
	}

	rows := model.ToTableRows(result.DatalistInf.TableInf)

	w := cmd.OutOrStdout()
	if w == os.Stdout {
		w = os.Stdout
	}
	format := cmdutil.GetFormat(cmd)
	return output.New(format).Format(w, rows)
}

func buildParams(cmd *cobra.Command, args []string) map[string]string {
	params := make(map[string]string)

	if len(args) > 0 {
		params["searchWord"] = args[0]
	}
	if v, _ := cmd.Flags().GetString("survey"); v != "" {
		params["surveyYears"] = v
	}
	if v, _ := cmd.Flags().GetString("field"); v != "" {
		params["statsField"] = v
	}
	if v, _ := cmd.Flags().GetString("open-year"); v != "" {
		params["openYears"] = v
	}
	if v, _ := cmd.Flags().GetString("stats-code"); v != "" {
		params["statsCode"] = v
	}
	if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
		params["limit"] = fmt.Sprintf("%d", v)
	}
	if v, _ := cmd.Flags().GetInt("start"); v > 0 {
		params["startPosition"] = fmt.Sprintf("%d", v)
	}

	return params
}
```

Note: Add missing `"fmt"` import to the imports section:

```go
import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/estat-cli/cmd/cmdutil"
	"github.com/planitaicojp/estat-cli/internal/api"
	"github.com/planitaicojp/estat-cli/internal/model"
	"github.com/planitaicojp/estat-cli/internal/output"
)
```

- [ ] **Step 4: テスト実行して成功確認**

```bash
go test ./cmd/search/ -v
```

Expected: PASS

- [ ] **Step 5: コミット**

```bash
git add cmd/search/
git commit -m "feat: estat search コマンドを実装（getStatsList API連携）"
```

---

### Task 10: ビルド確認 + エンドツーエンド動作確認

**Files:**
- Modify: `main.go` (already created in Task 1)

- [ ] **Step 1: ビルド**

```bash
cd /root/dev/planitai/estat-cli
go build -o estat .
```

Expected: バイナリ `estat` が生成される

- [ ] **Step 2: ヘルプ表示確認**

```bash
./estat --help
./estat search --help
./estat version
```

Expected: 日本語ヘルプメッセージが表示される

- [ ] **Step 3: 全テスト実行**

```bash
go test ./... -v
```

Expected: 全テスト PASS

- [ ] **Step 4: go vet + 整形チェック**

```bash
go vet ./...
gofmt -l .
```

Expected: エラーなし、未整形ファイルなし

- [ ] **Step 5: estat バイナリを .gitignore 確認 + コミット**

```bash
git status
git add -A
git commit -m "feat: estat-cli v0.1.0 プロジェクト骨格完成

- Go モジュール初期化
- Cobra ベースのサブコマンド構造
- YAML 設定 + 環境変数オーバーライド
- 出力フォーマッター（table/json/csv）
- e-Stat API クライアント基盤
- estat search コマンド（getStatsList API 連携）
- エラー型 + 終了コード体系
- サブコマンドスタブ（get, meta, dataset, catalog）"
```

---

### Task 11: .goreleaser.yaml + .golangci.yml

**Files:**
- Create: `.goreleaser.yaml`
- Create: `.golangci.yml`

- [ ] **Step 1: .goreleaser.yaml 作成**

```yaml
# .goreleaser.yaml
version: 2
builds:
  - main: .
    binary: estat
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X github.com/planitaicojp/estat-cli/cmd.version={{.Version}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
```

- [ ] **Step 2: .golangci.yml 作成**

```yaml
# .golangci.yml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
```

- [ ] **Step 3: コミット**

```bash
git add .goreleaser.yaml .golangci.yml
git commit -m "chore: GoReleaser とリンター設定を追加"
```
