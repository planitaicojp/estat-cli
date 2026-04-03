# estat-cli

## 概要

e-Stat（政府統計の総合窓口）APIを操作するためのCLIツール。

## 対象API

- **提供元**: 総務省統計局
- **API**: e-Stat API (https://www.e-stat.go.jp/api/)
- **認証**: アプリケーションID（appId）が必要（無料登録）
- **形式**: REST/JSON

## 主な機能

- 統計表の検索・一覧取得
- 統計データの取得・フィルタリング（地域、時期、分類項目など）
- メタ情報（統計分野、調査名、統計表ID）の取得
- データのCSV/JSON出力
- 統計分野コード・地域コードの参照

## 既存ツールの状況

- R言語パッケージ（estatapi）のみ存在
- Python/Go/Node.js向けのSDK・CLIは事実上存在しない

## 開発方針

- Go言語で実装
- サブコマンド構成（`estat search`, `estat get`, `estat meta` など）
- 出力形式: JSON（デフォルト）、CSV、テーブル表示
- Claude Code等のAIエージェントからの利用を想定した設計
- 日本語ドキュメント・ヘルプメッセージ
