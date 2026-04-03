# estat-cli プロジェクト構造設計

## 概要

e-Stat（政府統計の総合窓口）APIを操作するためのGo言語CLIツール。
conoha-cliのレイヤー分離パターンを参考に、e-Stat APIの単純さに合わせて軽量化した構造を採用する。

## 決定事項

| 項目 | 決定 |
|------|------|
| モジュール名 | `github.com/planitaicojp/estat-cli` |
| バイナリ名 | `estat` |
| CLIフレームワーク | spf13/cobra |
| 設定方式 | `~/.config/estat/config.yaml` + 環境変数オーバーライド |
| 出力形式 | table (デフォルト), JSON, CSV |
| 言語 | 日本語（ヘルプ・エラーメッセージ） |
| AI エージェント対応 | `--format json` で十分（追加のマシンモード不要） |

## 対象 API

- **ベースURL**: `https://api.e-stat.go.jp/rest/3.0/app`
- **認証**: `appId` クエリパラメータ（無料登録）
- **レスポンス**: JSON（`/json/` パスプレフィックス使用）

### エンドポイント一覧

| 機能 | エンドポイント | メソッド |
|------|---------------|---------|
| 統計表情報取得 | `/rest/3.0/app/json/getStatsList` | GET |
| メタ情報取得 | `/rest/3.0/app/json/getMetaInfo` | GET |
| 統計データ取得 | `/rest/3.0/app/json/getStatsData` | GET |
| データセット登録 | `/rest/3.0/app/postDataset` | POST |
| データセット参照 | `/rest/3.0/app/json/refDataset` | GET |
| データカタログ取得 | `/rest/3.0/app/json/getDataCatalog` | GET |
| 統計データ一括取得 | `/rest/3.0/app/getStatsDatas` | POST |

## プロジェクト構造

```
estat-cli/
├── main.go                      # エントリポイント（cmd.Execute()呼び出し）
├── go.mod
├── .goreleaser.yaml             # リリース設定
├── .golangci.yml                # Linter設定
├── cmd/
│   ├── root.go                  # ルートコマンド + グローバルフラグ
│   ├── version.go               # バージョン表示
│   ├── completion.go            # シェル補完（bash/zsh/fish/powershell）
│   ├── search/                  # 統計表検索（getStatsList）
│   │   └── search.go
│   ├── get/                     # 統計データ取得（getStatsData, getStatsDatas）
│   │   └── get.go
│   ├── meta/                    # メタ情報取得（getMetaInfo）
│   │   └── meta.go
│   ├── dataset/                 # データセット管理（postDataset, refDataset）
│   │   └── dataset.go
│   ├── catalog/                 # データカタログ（getDataCatalog）
│   │   └── catalog.go
│   └── cmdutil/
│       ├── client.go            # APIクライアント生成ヘルパー
│       └── args.go              # 引数バリデーションヘルパー
├── internal/
│   ├── api/
│   │   ├── client.go            # ベースHTTPクライアント
│   │   ├── statslist.go         # getStatsList
│   │   ├── metainfo.go          # getMetaInfo
│   │   ├── statsdata.go         # getStatsData / getStatsDatas
│   │   ├── dataset.go           # postDataset / refDataset
│   │   └── datacatalog.go       # getDataCatalog
│   ├── config/
│   │   ├── config.go            # YAML設定ロード・保存
│   │   └── env.go               # 環境変数オーバーライド
│   ├── model/
│   │   ├── statslist.go         # 統計表モデル
│   │   ├── metainfo.go          # メタ情報モデル
│   │   ├── statsdata.go         # 統計データモデル
│   │   ├── dataset.go           # データセットモデル
│   │   └── datacatalog.go       # カタログモデル
│   ├── output/
│   │   ├── formatter.go         # Formatterインターフェース + ファクトリ
│   │   ├── table.go             # テーブル出力（tabwriter）
│   │   ├── json.go              # JSON出力
│   │   └── csv.go               # CSV出力
│   └── errors/
│       ├── errors.go            # カスタムエラー型
│       └── exitcodes.go         # 終了コード定義
└── docs/
```

## 設定システム

### 設定ファイル

パス: `~/.config/estat/config.yaml`

```yaml
app_id: "your-app-id-here"
format: "table"
lang: "J"
```

### 環境変数

| 環境変数 | 用途 |
|---------|------|
| `ESTAT_APP_ID` | アプリケーションID |
| `ESTAT_FORMAT` | 出力形式 |
| `ESTAT_LANG` | 言語（J/E） |
| `ESTAT_CONFIG_DIR` | 設定ディレクトリ変更 |

### 優先順位

コマンドフラグ > 環境変数 > config.yaml > デフォルト値

### グローバルフラグ

```
--app-id string    アプリケーションID
--format string    出力形式: table, json, csv（デフォルト: table）
--lang string      言語: J(日本語), E(英語)（デフォルト: J）
--no-color         色出力を無効にする
--verbose          詳細出力
```

## サブコマンド設計

### `estat search` — 統計表検索

```
estat search [キーワード]
estat search --survey "国勢調査"
estat search --field "02"           # 統計分野コード
estat search --open-year 2020
estat search --limit 20
```

getStatsList を呼び出し。統計表ID、統計表名、調査名などを表示。

### `estat get` — 統計データ取得

```
estat get <statsDataId>
estat get <statsDataId> --area "13000"     # 東京都
estat get <statsDataId> --time "2020"
estat get <statsDataId> --category "cat01=A01"
estat get <statsDataId> --section-header
estat get --bulk <statsDataId1> <statsDataId2>
```

- 通常: getStatsData（GET）
- `--bulk`: getStatsDatas（POST、一括取得）

### `estat meta` — メタ情報取得

```
estat meta <statsDataId>
estat meta <statsDataId> --class "area"    # 特定分類のみ
```

getMetaInfo を呼び出し。分類項目、地域コード、時間軸などを表示。

### `estat dataset` — データセット管理

```
estat dataset register --stats-data-id <id> --filter "..."   # 登録
estat dataset show <datasetId>                                # 参照
```

- `register`: postDataset（POST）
- `show`: refDataset（GET）

### `estat catalog` — データカタログ

```
estat catalog
estat catalog --survey "国勢調査"
estat catalog --field "02"
estat catalog --dataset-type db    # db or file
```

getDataCatalog を呼び出し。

## APIクライアント

### Client構造体

```go
type Client struct {
    HTTPClient *http.Client
    AppID      string
    BaseURL    string  // https://api.e-stat.go.jp/rest/3.0/app
    Lang       string
    Verbose    bool
}
```

### 設計方針

- タイムアウト: 30秒
- 共通パラメータ（appId, lang）自動付与
- JSON レスポンス（`/json/` パスプレフィックス）
- Verbose モード時: リクエストURL + レスポンスステータスをstderrに出力
- User-Agent: `planitaicojp/estat-cli/{version}`

### データフロー

```
コマンド（cmd/search/）
  ↓ フラグパース
cmdutil.NewClient(cmd)
  ↓ config ロード → 環境変数オーバーライド → フラグオーバーライド
api.Client 生成
  ↓
api.GetStatsList(params)
  ↓ HTTP GET → JSONパース
model.StatsListResponse
  ↓
output.Format(os.Stdout, rows)
  ↓ table / json / csv
stdout 出力
```

## 出力フォーマッター

### インターフェース

```go
type Formatter interface {
    Format(w io.Writer, data any) error
}
```

`output.New(format string) Formatter` ファクトリで生成。

### Table（デフォルト）

- `text/tabwriter` 使用
- 構造体のJSONタグからヘッダー名を抽出
- 出力例:
```
統計表ID         統計表名                          調査名
0003410379       男女別人口－全国，都道府県          国勢調査
0003410380       年齢各歳別人口                     国勢調査
```

### JSON

- `json.MarshalIndent`（2スペースインデント）

### CSV

- `encoding/csv` 使用
- ヘッダー行を含む

## エラー処理

### エラー型

| エラー型 | 状況 | 終了コード |
|---------|------|-----------|
| `ConfigError` | appId未設定など | 2 |
| `APIError` | e-Stat APIエラーレスポンス | 3 |
| `NetworkError` | ネットワーク障害 | 4 |
| `ValidationError` | 不正な引数 | 5 |

### API エラーハンドリング

e-Stat APIのエラーレスポンスの `STATUS` フィールド（コード + メッセージ）をパースし、日本語エラーメッセージとして表示。

## ビルド・リリース

- GoReleaser を使用
- バイナリ名: `estat`
- ターゲット: Linux, macOS, Windows（amd64, arm64）
- バージョン注入: ldflags `-X github.com/planitaicojp/estat-cli/cmd.version={{.Version}}`
