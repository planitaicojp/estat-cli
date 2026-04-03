# estat - e-Stat API CLI

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

e-Stat（政府統計の総合窓口）API 用のコマンドラインインターフェースです。Go で書かれたシングルバイナリで、AI エージェントフレンドリーな設計を採用しています。

## 特徴

- シングルバイナリ、クロスプラットフォーム対応（Linux / macOS / Windows）
- 構造化出力（`--format json/csv/table`）
- AI エージェントフレンドリー設計（JSON 出力、決定的な終了コード、stderr/stdout 分離）
- 日本語ヘルプ・エラーメッセージ

## インストール

### Homebrew (macOS / Linux)

```bash
brew install planitaicojp/tap/estat
```

### Scoop (Windows)

```powershell
scoop bucket add planitaicojp https://github.com/planitaicojp/scoop-bucket
scoop install estat
```

### ソースからビルド

```bash
go install github.com/planitaicojp/estat-cli@latest
```

### リリースバイナリ

[Releases](https://github.com/planitaicojp/estat-cli/releases) ページからダウンロード

**Linux (amd64)**
```bash
curl -Lo estat https://github.com/planitaicojp/estat-cli/releases/latest/download/estat-cli_linux_amd64.tar.gz
tar xzf estat-cli_linux_amd64.tar.gz
chmod +x estat
sudo mv estat /usr/local/bin/
```

**macOS (Apple Silicon)**
```bash
curl -Lo estat.tar.gz https://github.com/planitaicojp/estat-cli/releases/latest/download/estat-cli_darwin_arm64.tar.gz
tar xzf estat.tar.gz
chmod +x estat
sudo mv estat /usr/local/bin/
```

## 事前準備

e-Stat API を利用するにはアプリケーション ID が必要です（無料）:

1. [e-Stat API 利用登録](https://www.e-stat.go.jp/api/) でアカウント作成
2. アプリケーション ID を取得
3. 以下のいずれかの方法で設定:

```bash
# 環境変数で設定
export ESTAT_APP_ID="your-app-id"

# または設定ファイルに保存
mkdir -p ~/.config/estat
echo 'app_id: "your-app-id"' > ~/.config/estat/config.yaml
```

## クイックスタート

```bash
# 統計表を検索
estat search 人口

# 国勢調査の統計表を検索
estat search --survey "国勢調査"

# JSON 形式で出力
estat search 人口 --format json

# CSV 形式で出力
estat search 人口 --format csv

# 取得件数を制限
estat search 人口 --limit 5
```

## コマンド一覧

| コマンド | 説明 | 状態 |
|---------|------|------|
| `estat search` | 統計表を検索 | ✅ 実装済み |
| `estat get` | 統計データを取得 | 🚧 開発中 |
| `estat meta` | メタ情報を取得 | 🚧 開発中 |
| `estat dataset` | データセットを管理 | 🚧 開発中 |
| `estat catalog` | データカタログを取得 | 🚧 開発中 |
| `estat version` | バージョン情報を表示 | ✅ 実装済み |
| `estat completion` | シェル補完を生成 | ✅ 実装済み |

## 設定

設定ファイル: `~/.config/estat/config.yaml`

```yaml
app_id: "your-app-id"
format: "table"    # デフォルト出力形式
lang: "J"          # J:日本語, E:英語
```

### 環境変数

| 変数 | 説明 |
|-----|------|
| `ESTAT_APP_ID` | アプリケーション ID |
| `ESTAT_FORMAT` | 出力形式 |
| `ESTAT_LANG` | 言語 |
| `ESTAT_CONFIG_DIR` | 設定ディレクトリ |

優先順位: コマンドフラグ > 環境変数 > 設定ファイル > デフォルト値

### グローバルフラグ

```
--app-id     アプリケーション ID
--format     出力形式（table / json / csv）
--lang       言語（J / E）
--no-color   カラー出力を無効化
--verbose    詳細出力
```

## 終了コード

| コード | 意味 |
|-------|------|
| 0 | 成功 |
| 1 | 一般エラー |
| 2 | 設定エラー |
| 3 | API エラー |
| 4 | ネットワークエラー |
| 5 | バリデーションエラー |

## AI エージェント連携

本 CLI は AI エージェントからの利用を想定して設計されています:

```bash
# JSON 形式で統計表を検索
estat search 人口 --format json

# 終了コードでエラーハンドリング
estat search 人口 || echo "Exit code: $?"
```

## 開発

```bash
make build     # バイナリをビルド
make test      # テストを実行
make lint      # リンターを実行（golangci-lint 必要）
make vet       # go vet を実行
make clean     # 成果物を削除
```

## API ドキュメント

- [e-Stat API 仕様 3.0](https://www.e-stat.go.jp/api/api-info/e-stat-manual3-0)
- [e-Stat API 利用ガイド](https://www.e-stat.go.jp/api/)

## ライセンス

MIT License
